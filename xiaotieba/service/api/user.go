package api

import (
	"errors"
	"net/http"
	"time"
	db "xiaotieba/db/sqlc"
	"xiaotieba/util"

	"xiaotieba/service/worker"

	"github.com/google/uuid"
	"github.com/hibiken/asynq"
	_ "github.com/lib/pq" // 导入PostgreSQL驱动

	//匿名引用是因为不直接调用但需要驱动程序
	"github.com/gin-gonic/gin"
)

type ListDiligentUserReq struct {
	MinPostCounts int `uri:"minPostCounts"` //上榜用户最小发帖数
}

// 新功能引入：根据发帖数量对用户排名 complex.sql.go里名为 func ListUser()
// complex.sql.go里名为 func ListUser()，这个函数返回一个名为UserRow的切片，元素为User结构体，按post数量排名
func (server *Server) ListDiligentUser(ctx *gin.Context) {
	var req ListDiligentUserReq
	if err := ctx.ShouldBindUri(req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	UserRow, err := server.querier.ListUser(ctx, req.MinPostCounts)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	ctx.JSON(http.StatusOK, UserRow)
}

type createUserRequest struct {
	Username     string   `json:"username" binding:"require,alphanum"`
	HashPassword string   `json:"hashpassword" binding:"require,min=6"`
	Power        string   `json:"power" binding:"require"`
	Email        string   `json:"email" binding:"require"`
	Phone        string   `json:"phone" binding:"require"`
	Roles        []string `json:"roles" binding:"require"`
}

// 不用return，因为这个函数用于http协议通信，只需要发送响应即可
func (server *Server) SignUpUser(ctx *gin.Context) {
	var req createUserRequest
	if err := ctx.ShouldBindJSON(req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
	}
	createuser := db.CreateUserParams{
		HashPassword: req.HashPassword,
		Name:         req.Username,
		Power:        req.Power,
		Email:        req.Email,
		Phone:        req.Phone,
	}

	user, err := server.querier.CreateUser(ctx, createuser)
	if err != nil {
		errCode := db.ErrorCode(err)
		if errCode == db.ForeignKeyViolation || errCode == db.UniqueViolation {
			ctx.JSON(http.StatusForbidden, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	for _, role := range req.Roles {
		server.enforcer.AddRoleForUser(user.Name, role)
	}
	ctx.JSON(http.StatusOK, normalResponce("user", user))

}

type verifyEmailrequest struct {
	Username string `json:"username" binding:"require"`
}

func (server *Server) VerifyEmail(ctx *gin.Context) {
	var req verifyEmailrequest
	if err := ctx.ShouldBindJSON(req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	payload := &worker.PayloadSendVerifyEmail{
		Username: req.Username,
	}
	//配置参数
	opts := []asynq.Option{
		asynq.MaxRetry(10),                 //最大重试次数
		asynq.ProcessIn(10 * time.Second),  //异步处理发生于发出任务的多长时间之后
		asynq.Queue(worker.QueueXiaotieba), //队列名称
	}
	server.distributor.DistributorTaskSendEmail(ctx, payload, opts...)
	// user,err:=server.querier.GetUserByUsername(ctx,req.Username)
	// if err!=nil{
	// 	ctx.JSON(http.StatusInternalServerError, errorResponse(err))
	// 	return
	// }

}

type loginUserRequest struct {
	Username     string `json:"username" binding:"require,alphanum"`
	Hashpassword string `json:"hashpassword" binding:"require,min=6"`
}

// type userResponse struct {
// 	Username          string    `json:"username"`
// 	FullName          string    `json:"full_name"`
// 	Email             string    `json:"email"`
// 	PasswordChangedAt time.Time `json:"password_changed_at"`
// 	CreatedAt         time.Time `json:"created_at"`
// }

type loginUserResponse struct {
	SessionID             uuid.UUID `json:"session_id"`
	AccessToken           string    `json:"access_token"`
	AccessTokenExpiresAt  time.Time `json:"access_token_expires_at"`
	RefreshToken          string    `json:"refresh_token"`
	RefreshTokenExpiresAt time.Time `json:"refresh_token_expires_at"`
	Username              string    `json:"username"`
	Password              string    `json:"password"`
}

func (server *Server) loginUser(ctx *gin.Context) {
	var req loginUserRequest
	if err := ctx.ShouldBindJSON(req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	user, err := server.querier.GetUserByUsername(ctx, req.Username)
	if err != nil {
		if err != db.ErrRecordNotFound {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	//检查传入的密码和db的哈希密码是否一致
	if err = util.CheckPassWord(req.Username, user.HashPassword); err != nil {
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	accessToken, accesspayload, _ := server.tokenMaker.CreateToken(req.Username, server.config.AccessTokenDuration)
	refreshToken, refreshpayload, _ := server.tokenMaker.CreateToken(req.Username, server.config.RefreshTokenDuration)
	sessionparam := &db.CreateSessionParams{
		ID:           refreshpayload.ID,
		Username:     user.Name,
		RefreshToken: refreshToken,
		Useragent:    ctx.Request.UserAgent(),
		ClientIp:     ctx.ClientIP(),
		Isblocked:    false,
		ExpireAt:     refreshpayload.ExpireAt,
		//Random出来的tokenID即为session的ID(前提是refreshToken,因为session里只存储了refreshToken)
	}
	session, err := server.querier.CreateSession(ctx, *sessionparam)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	responce := loginUserResponse{
		SessionID:             session.ID,
		AccessToken:           accessToken,
		AccessTokenExpiresAt:  accesspayload.ExpireAt,
		RefreshToken:          refreshToken,
		RefreshTokenExpiresAt: refreshpayload.ExpireAt,
		Username:              req.Username,
		Password:              req.Hashpassword,
	}
	ctx.JSON(http.StatusOK, responce)

}

type getUserRequest struct {
	ID int64 `uri:"id" binding:"require,min=1"`
}

func (server *Server) getUser(ctx *gin.Context) {
	var req getUserRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		//不用return，因为这个函数用于http协议通信，只需要发送响应即可
	}

	user, err := server.querier.GetUserByID(ctx, req.ID)
	if err != nil {
		if errors.Is(err, db.ErrRecordNotFound) {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}

		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, normalResponce("user", user))
}

// func (server *Server) DeleteUser(ctx *gin.Context) *db.User {
// 	var req getUserRequest
// 	ctx.
// }

//TODO: 实现用户删除(为异步操作，最好两天内可以取消删除)

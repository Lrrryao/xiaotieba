package api

import (
	"errors"
	"fmt"
	"net/http"
	db "xiaotieba/db/sqlc"

	_ "github.com/lib/pq" // 导入PostgreSQL驱动

	//匿名引用是因为不直接调用但需要驱动程序
	"github.com/gin-gonic/gin"
)

type createPostRequest struct {
	User_ID  int64  `json:"user_id" binding:"require,min=1"`
	Titles   string `json:"title" binding:"require"`
	Content  string `json:"content" binding:"require"`
	Username string `json:"username" binding:"require"`
}

// 不用return，因为这个函数用于http协议通信，只需要发送响应即可
func (server *Server) CreatePost(ctx *gin.Context) {
	var req createPostRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
	}
	createpost := db.CreatePostParams{
		UserID:  req.User_ID,
		Content: req.Content,
		Titles:  req.Titles,
	}

	post, err := server.querier.CreatePost(ctx, createpost)
	if err != nil {
		errCode := db.ErrorCode(err)
		if errCode == db.ForeignKeyViolation || errCode == db.UniqueViolation {
			ctx.JSON(http.StatusForbidden, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	id := fmt.Sprintf("%v", post.ID)

	//建立新的policy，赋予该该角色postOwner在domain："/post/"+id里删除帖子和回复的权限
	server.enforcer.AddPolicy("postOwner", "/post/"+id, "delete_post", "/post/"+id)
	server.enforcer.AddPolicy("postOwner", "/post/"+id, "delete_comment", "/post/"+id+"allcomments")
	//赋予该用户 在domain里的角色
	server.enforcer.AddRoleForUserInDomain(req.Username, "postOwner", "/post/"+id)

	ctx.JSON(http.StatusOK, normalResponce("post", post))

}

type getPostRequest struct {
	ID int64 `uri:"id" binding:"require,min=1"`
}

func (server *Server) getPost(ctx *gin.Context) {
	var req getPostRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		//不用return，因为这个函数用于http协议通信，只需要发送响应即可
	}

	post, err := server.querier.GetPost(ctx, req.ID)
	if err != nil {
		if errors.Is(err, db.ErrRecordNotFound) {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}

		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, normalResponce("post", post))
}

type deletePostRequest struct {
	ID       int64  `json:"id" binding:"require"`
	Username string `json:"username" binding:"require"`
	UserID   int64  `json:"userid" binding:"require"`
}

func (server *Server) DeletePost(ctx *gin.Context) {
	var req deletePostRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
	}
	id := fmt.Sprintf("%v", req.ID)
	//验证身份
	is, err := server.enforcer.Enforce(req.Username, "/post/"+id, "delete_post", "/post/"+id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	if !is {
		ctx.JSON(http.StatusUnauthorized, errorResponse(errors.New("you dont have the power")))
		return
	}

	if err := server.querier.DeletePost(ctx, req.ID); err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	ctx.JSON(http.StatusOK, normalResponce("delete successfully", req.ID))

}

//TODO: 实现帖子删除(为异步操作，最好两天内可以取消删除)

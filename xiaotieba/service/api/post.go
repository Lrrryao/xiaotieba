package api

import (
	"errors"
	"net/http"
	db "xiaotieba/db/sqlc"

	_ "github.com/lib/pq" // 导入PostgreSQL驱动
	//匿名引用是因为不直接调用但需要驱动程序
	"github.com/gin-gonic/gin"
)

type createPostRequest struct {
	User_ID int64  `json:"user_id" binding:"require,min=1"`
	Titles  string `json:"title" binding:"require"`
	Content string `json:"content" binding:"require"`
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

// func (server *Server) DeletePost(ctx *gin.Context) *db.Post {
// 	var req getPostRequest
// 	ctx.
// }

//TODO: 实现帖子删除(为异步操作，最好两天内可以取消删除)

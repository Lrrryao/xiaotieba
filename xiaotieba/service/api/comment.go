package api

import (
	"errors"
	"net/http"
	db "xiaotieba/db/sqlc"

	_ "github.com/lib/pq" // 导入CommentgreSQL驱动
	//匿名引用是因为不直接调用但需要驱动程序
	"github.com/gin-gonic/gin"
)

type createCommentRequest struct {
	User_ID int64  `json:"user_id" binding:"require,min=1"`
	Post_ID int64  `json:"post_id" binding:"require,min=1"`
	Content string `json:"content" binding:"require"`
}

// 不用return，因为这个函数用于http协议通信，只需要发送响应即可
func (server *Server) CreateComment(ctx *gin.Context) {
	var req createCommentRequest
	if err := ctx.ShouldBindJSON(req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
	}
	createcomment := db.CreateCommentParams{
		UserID:  req.User_ID,
		Content: req.Content,
		PostID:  req.Post_ID,
	}

	comment, err := server.querier.CreateComment(ctx, createcomment)
	if err != nil {
		errCode := db.ErrorCode(err)
		if errCode == db.ForeignKeyViolation || errCode == db.UniqueViolation {
			ctx.JSON(http.StatusForbidden, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	ctx.JSON(http.StatusOK, normalResponce("comment", comment))

}

type getCommentRequest struct {
	ID int64 `uri:"id" binding:"require,min=1"`
}

func (server *Server) getComment(ctx *gin.Context) {
	var req getCommentRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		//不用return，因为这个函数用于http协议通信，只需要发送响应即可
	}

	comment, err := server.querier.GetComment(ctx, req.ID)
	if err != nil {
		if errors.Is(err, db.ErrRecordNotFound) {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}

		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	ctx.JSON(http.StatusOK, normalResponce("comment", comment))
}

// func (server *Server) DeleteComment(ctx *gin.Context) *db.Comment {
// 	var req getCommentRequest
// 	ctx.
// }

//TODO: 实现帖子删除(为异步操作，最好两天内可以取消删除)

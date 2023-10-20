package api

import (
	"errors"
	"fmt"
	"net/http"
	"time"
	db "xiaotieba/db/sqlc"

	"github.com/gin-gonic/gin"
)

//关于该api的作用的理解
//已知accessToken是短期令牌，负责验证对于服务器资源的请求，
//而refreshToken则是一种长期令牌，用于刷新access令牌，
//由于accessToken有效期过短，我们不需要在session里保存它，否则造成表内存浪费
//由于accessToken刷新时间较短，renewRefreshToken函数用于在验证了refreshToken和session的合法性后，颁发新的accessToken

type renewRefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

type refreshTokenResponse struct {
	AccessToken          string    `json:"access_token"`
	AccessTokenExpiresAt time.Time `json:"access_token_expires_at"`
}

func (server *Server) refreshToken(ctx *gin.Context) {
	var req renewRefreshTokenRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	//通过VerifyToken验证该token是否有效
	refreshpayload, err := server.tokenMaker.VerifyToken(req.RefreshToken)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}
	//通过session验证token里的内容是否真实存在，如ID、是否blocked、refreshToken是否相同
	session, err := server.querier.GetSession(ctx, refreshpayload.ID)
	if err != nil {
		if errors.Is(err, db.ErrRecordNotFound) /*ErrRecordNotFound就是"sql.ErrNoRows"*/ {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}

		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	if session.Isblocked {
		err := fmt.Errorf("session is blocked ")
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
	}
	if refreshpayload.ID != session.ID {
		err := fmt.Errorf("mismatched session ID ")
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
	}
	if req.RefreshToken != session.RefreshToken {
		err := fmt.Errorf("mismatched session ID ")
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
	}

	if time.Now().After(session.ExpireAt) {
		err := fmt.Errorf("expired token")
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
	}
	accessToken, accessPayload, err := server.tokenMaker.CreateToken(refreshpayload.Username, server.config.AccessTokenDuration)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	var Rsp = &refreshTokenResponse{
		AccessToken:          accessToken,
		AccessTokenExpiresAt: accessPayload.ExpireAt,
	}
	ctx.JSON(http.StatusOK, Rsp)

}

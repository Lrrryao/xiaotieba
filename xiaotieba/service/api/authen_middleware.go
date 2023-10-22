package api

import (
	"errors"
	"net/http"
	"strings"
	"xiaotieba/token"

	"github.com/gin-gonic/gin"
)

const (
	authorizationHeaderKey  = "authorization" //用于验证header里是否有value为"authorization"的key
	authorizationTypeBearer = "bearer"        //验证类型为bearer
	authorizationPayloadKey = "authorization_payload"
)

type Payload struct {
}

// 该middleware用于检测请求头中的header中是否有authorization
// 以及authorization中的token是否有效

func authMiddleware(maker token.Maker) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authenHeader := ctx.GetHeader(authorizationHeaderKey)
		if len(authenHeader) == 0 {
			err := errors.New("Authorization not provided")
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(err))
			return
		}
		fields := strings.Fields(authenHeader) //返回一个字符串切片
		if len(fields) < 2 {
			err := errors.New("invalid authorization header format")
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(err))
			return
		}

		authenType := strings.ToLower(fields[0])
		if authenType != authorizationTypeBearer {
			err := errors.New("unsupported authorization type ")
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(err))
			return
		}
		accessToken := fields[1]
		payload, err := maker.VerifyToken(accessToken)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(err))
			return
		}
		ctx.Set(authorizationPayloadKey, payload)
		ctx.Set("request username", payload.Username)
		ctx.Next()
	}
}

//将一个名为 payload 的值关联到一个名为 authorizationPayloadKey 的键上，
//然后将这个键值对存储在一个上下文对象 ctx 中

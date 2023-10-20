package api

import (
	"fmt"
	db "xiaotieba/db/sqlc" //由于包名都是小写，而小写包名只能被内部调用，所以server不能放在目录里面
	"xiaotieba/service/worker"
	"xiaotieba/token"
	"xiaotieba/util"

	"github.com/gin-gonic/gin"
)

type Server struct {
	//结构体或者接口若想被包外调用，名称必为大写
	querier     db.Querier
	config      util.Config
	tokenMaker  token.Maker
	router      *gin.Engine
	distributor worker.TaskDistributor
}

func NewServer(querier db.Querier, config util.Config, taskDistributor worker.TaskDistributor) (*Server, error) {
	tokenMaker, err := token.NewPaseto_Maker(config.TokenSymmetricKey)
	if err != nil {
		if err != nil {
			return nil, fmt.Errorf("cannot create token maker: %w", err)
		}
	}

	server := &Server{
		querier:     querier,
		config:      config,
		tokenMaker:  tokenMaker,
		router:      nil,
		distributor: taskDistributor,
	}

	server.SetupRouter()
	return server, nil

}

func (server *Server) SetupRouter() {
	router := gin.Default()
	router.POST("/api/users", server.SignUpUser)

	//将router这一个engine对象转为GroupRouter对象，后面的函数为中间件
	router.Group("/").Use(authMiddleware(server.tokenMaker))

	router.GET("/api/users:id", server.getUser)
	router.POST("/api/sessions", server.refreshToken)

	//TODO:对handler和url搭建响应路由
	//
	//
	//
	//
	//
	server.router = router
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}

func normalResponce(msg string, body interface{}) gin.H {
	return gin.H{msg: body}
}

//返回一个键值对，key为“error”，值为err.Error()

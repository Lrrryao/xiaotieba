package api

import (
	"fmt"
	db "xiaotieba/db/sqlc" //由于包名都是小写，而小写包名只能被内部调用，所以server不能放在目录里面
	"xiaotieba/service/worker"
	"xiaotieba/token"
	"xiaotieba/util"

	"xiaotieba/cache"

	"github.com/casbin/casbin/v2"
	"github.com/gin-gonic/gin"
)

type Server struct {
	//结构体或者接口若想被包外调用，名称必为大写
	querier     db.Querier
	config      util.Config
	tokenMaker  token.Maker
	router      *gin.Engine
	distributor worker.TaskDistributor
	//enforcer    rbac.Enforcer
	enforcer *casbin.Enforcer
	cache    cache.CacheStore
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
		enforcer:    nil,
	}

	server.SetupRouter()
	return server, nil

}
func NewEnforcer(model, policy string) (*casbin.Enforcer, error) {
	e, err := casbin.NewEnforcer(model, policy)

	e.AddNamedMatchingFunc("g2", "KeyMatch2", util.KeyMatch2)
	e.AddNamedDomainMatchingFunc("g2", "KeyMatch2", util.KeyMatch2)
	return e, err
}
func (server *Server) RunRedisCache(Addr string, Password string, DB int) {
	server.cache = cache.NewRedisCache(Addr, Password, DB)
}

// /
func (server *Server) StartEnforcer(model string, policy string) error {
	var err error
	server.enforcer, err = NewEnforcer(model, policy)
	if err != nil {
		return err
	}
	return nil
}

func (server *Server) SetupRouter() {
	router := gin.Default()
	router.Use(ErrorLogger())
	router.POST("/api/users", server.SignUpUser)

	//将router这一个engine对象转为GroupRouter对象，后面的函数为中间件
	router.Group("/").Use(authMiddleware(server.tokenMaker))

	router.GET("/api/users:id", server.getUser)
	router.GET("/api/posts:id", server.getPost)
	router.DELETE("/api/posts", server.DeletePost)
	router.POST("/api/sessions", server.refreshToken)
	router.POST("/api/login", server.loginUser)
	router.POST("/api/vote", server.vote)

	//TODO:对其他handler和url搭建响应路由
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

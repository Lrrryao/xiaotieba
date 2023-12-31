package main

import (
	"database/sql"
	db "xiaotieba/db/sqlc"
	"xiaotieba/mail"
	"xiaotieba/service/api"
	"xiaotieba/service/worker"
	"xiaotieba/util"

	"github.com/hibiken/asynq"
	"github.com/rs/zerolog/log"
)

func main() {
	//读取配置并将其赋值给config
	config, err := util.LoadConfig(".")

	if err != nil {
		log.Fatal().Err(err).Msg("cannot load config")
	}
	//建立数据库连接
	conn, err := sql.Open(config.DBdriver, config.DBsource)
	if err != nil {
		log.Fatal().Err(err).Msg("cannot connect to db")
	}
	querier := db.New(conn)

	//初始化redis
	//创建生产者
	redisOpt := asynq.RedisClientOpt{
		Addr: config.RedisAddress,
	}
	taskDistributor := worker.NewRedisTaskDistributor(redisOpt)

	//初始化服务器
	server, err := api.NewServer(querier, config, taskDistributor)
	if err != nil {
		log.Fatal().Err(err).Msg("cannot create server")
	}

	//启动reids
	server.RunRedisCache(config.RedisAddress, config.RedisPassword, 0)

	//启动RBAC管理权限模式
	server.StartEnforcer("./model.conf", "./policy.csv")
	server.SetupRouter()

	go runTaskProcessor(redisOpt, querier)
	//思考：为什么runTaskProcessor和worker.NewRedisTaskDistributor，即初始化生产者(distributor)和消费者(processor)时
	//都必须传入相同的asynq.RedisClientOpt? 即redisOps
	//思考结果：因为不同于http，redis中进行任务分发和处理不是通过网络协议，而是在redis上，它们需要相同的redis地址（即端口）
}

func runTaskProcessor(redisOpt asynq.RedisClientOpt, querier db.Querier) {
	config, _ := util.LoadConfig(".")
	mailSender := mail.NewQQSender(config.EMAIL_SENDER_NAME, config.EMAIL_SENDER_ADDRESS, config.EMAIL_SENDER_PASSWORD)
	processor := worker.NewRedisTaskProcessor(redisOpt, querier, mailSender)
	log.Info().Msg("start task processor")
	if err := processor.Start(); err != nil {
		log.Fatal().Err(err).Msg("Failed to run processor server")
	}
}

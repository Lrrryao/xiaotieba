package worker

import (
	"context"
	db "xiaotieba/db/sqlc"
	"xiaotieba/mail"

	"github.com/hibiken/asynq"
	"github.com/rs/zerolog/log"
)

type taskProcessor interface {
	ProcessTaskSendVerifyEmail(ctx context.Context, task *asynq.Task) error
	Start() error
}

// 接收redis异步任务的人
// 由于异步任务可能用到数据库，所以又用到了querier
// processor就相当于web中的服务器端，用于处理taskDistributor发出的任务
type RedisTaskProcessor struct {
	server     *asynq.Server
	querier    db.Querier
	mailSender mail.EmailSender
}

const (
	QueueXiaotieba = "xiaotieba"
	QueueDefault   = "default"
)

// 客户端的异步参数
func NewRedisTaskProcessor(RedisOpt asynq.RedisClientOpt, querier db.Querier, mailSender mail.EmailSender) taskProcessor {
	processor := asynq.NewServer(

		//客户端的异步参数
		RedisOpt,

		//服务器的参数：比如最大并发数，失败任务重试的延迟函数，判断处理函数是否失败的谓词函数，任务队列的名称映射及其优先级，等
		asynq.Config{
			Queues: map[string]int{ //string的key值是队列的名称，int值代表队列优先级
				QueueXiaotieba: 10,
				QueueDefault:   5,
			}, //如果空白表示默认参数

			//注册错误处理函数，即asynq.ErrorHandlerFunc这个type的定义写一个错误处理函数
			//并且将其强制转换为asynq.ErrorHandlerFunc这个类
			ErrorHandler: asynq.ErrorHandlerFunc(func(ctx context.Context, task *asynq.Task, err error) {
				log.Error().Err(err).Str("type", task.Type()).
					Bytes("payload", task.Payload()).Msg("Process task failed")
			}),

			Logger: NewLogger(),
		},
	)
	return &RedisTaskProcessor{
		server:     processor,
		querier:    querier,
		mailSender: mailSender,
	}
}

// 搭建processor的路由
func (processor *RedisTaskProcessor) Start() error {
	mux := asynq.NewServeMux()
	mux.HandleFunc(TaskSendVerifyEmail, processor.ProcessTaskSendVerifyEmail)
	return processor.server.Start(mux)
}

package worker

import (
	"context"

	"github.com/hibiken/asynq"
)

// 创建一个接口将有助于将来模拟测试
type TaskDistributor interface {
	DistributorTaskSendEmail(
		ctx context.Context,
		payload *PayloadSendVerifyEmail,
		opts ...asynq.Option,
	) error
}

// RedisTaskDistributor实现了异步任务的发送到redis的功能
type RedisTaskDistributor struct {
	client *asynq.Client //*asynq.Client
}

func NewRedisTaskDistributor(redisOps asynq.RedisClientOpt) TaskDistributor {
	client := asynq.NewClient(redisOps)
	return &RedisTaskDistributor{
		client: client,
	}
}

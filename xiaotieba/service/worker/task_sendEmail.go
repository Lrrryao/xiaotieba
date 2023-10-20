package worker

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	db "xiaotieba/db/sqlc"
	"xiaotieba/util"

	"github.com/hibiken/asynq"
	"github.com/rs/zerolog/log"
)

const (
	TaskSendVerifyEmail = "task:send_verify_email"
)

// 用于包含想要储存在redis中的所有数据
// 然后worker就会检查到他
// 在这个向用户发邮件这一任务中，只需要一个用户名就可以从数据库得到全部需要的信息
type PayloadSendVerifyEmail struct {
	Username string `json:"username"`
}

func (distributor *RedisTaskDistributor) DistributorTaskSendEmail(
	ctx context.Context,
	payload *PayloadSendVerifyEmail, //由于
	opts ...asynq.Option, //这个字段描述的是你一些你想要配置的参数，用于异步：任务分发（distribute），运行（function）或者重试（retry）
	//...表示可以传入可变数量的参数，即零个、一个、多个opts
) error {
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("fail to marshall payload to json: %v", err)
	}
	task := asynq.NewTask(TaskSendVerifyEmail, jsonPayload, opts...)
	info, err := distributor.client.EnqueueContext(ctx, task)
	if err != nil {
		return fmt.Errorf("failed to enqueue task: %v", err)
	}
	//打印日志
	log.Info().Str("type", task.Type()).Bytes("payload", task.Payload()).
		Str("queue", info.Queue).Int("max_retry", info.MaxRetry).Msg("enqueue task")
	return nil
}

func (processor *RedisTaskProcessor) ProcessTaskSendVerifyEmail(ctx context.Context, task *asynq.Task) error {
	var payload *PayloadSendVerifyEmail
	if err := json.Unmarshal(task.Payload(), payload); err != nil {
		return fmt.Errorf("failed to unmarshall payload: %v", asynq.SkipRetry) //asynq.SkipRetry告诉processor别重试了
	}
	user, err := processor.querier.GetUserByUsername(ctx, payload.Username)
	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("user doesn't exist : %v", asynq.SkipRetry) //如果username根本不存在，则别重试了，重试个屁
		} else {
			return fmt.Errorf("failed to get user : %v", err) //如果存在但是就是没法get到，那就再返回个错误试试
		}
	}
	arg := db.CreateVerifyEmailParams{ //待补充，应该新增一个验证邮箱有效性的代码，或者给该字段设置唯一索引然后检验也可以
		Username:   sql.NullString{String: user.Name, Valid: true},
		Email:      sql.NullString{String: user.Email, Valid: true},
		SecretCode: util.Random(), //生成一个随机码
	}

	verifyEmail, err := processor.querier.CreateVerifyEmail(ctx, arg)
	if err != nil {
		return fmt.Errorf("failed to create verify-email record : %v", err)
	}
	subject := "welcome to xiaotieba"
	content := fmt.Sprintf("welcome to xiaotieba:%s", verifyEmail.Username.String)
	//content这里应该包含一个url, 当用户点进去时，在数据库中的Is_verified字段改为true

	to := []string{verifyEmail.Email.String}

	err = processor.mailSender.SendEmail(subject, content, to, nil, nil, nil)
	if err != nil {
		return fmt.Errorf("failed to send email : %w", err)
	}

	//TODO:不仅要完成邮件，而且要将这个函数和loginUser中的数据库组合成事务。

	log.Info().Str("type", task.Type()).Bytes("payload", task.Payload()).
		Msg("enqueue task")
	return nil
}

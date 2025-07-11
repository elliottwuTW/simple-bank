package worker

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/hibiken/asynq"
	"github.com/simple_bank/database"
	"github.com/simple_bank/util"
)

const TaskSendVerifyEmail = "task:send_verify_email"

// PayloadSendVerifyEmail contains all data of the task that we want to store in MQ
type PayloadSendVerifyEmail struct {
	Username string `json:"username"`
}

func (d *RedisTaskDistributor) DistributeTaskSendVerifyEmail(
	ctx context.Context,
	payload *PayloadSendVerifyEmail,
	opt ...asynq.Option,
) error {
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal payload: %w", err)
	}

	task := asynq.NewTask(TaskSendVerifyEmail, jsonPayload, opt...)
	_, err = d.client.EnqueueContext(ctx, task)
	if err != nil {
		return fmt.Errorf("failed to enqueue task: %w", err)
	}

	// log.Info().Str("type", task.Type()).Bytes("payload", task.Payload()).
	// 	Str("queue", info.Queue).Int("max_retry", info.MaxRetry).Msg("enqueued task")
	return nil
}

func (p *RedisTaskProcessor) ProcessTaskSendVerifyEmail(
	ctx context.Context,
	task *asynq.Task,
) error {
	var payload PayloadSendVerifyEmail
	if err := json.Unmarshal(task.Payload(), &payload); err != nil {
		// 如果連 payload 都 unmarshal 不出來，沒有重試的必要
		// 回傳 asynq.SkipRetry 讓 server 不重試!!!!!!!
		return fmt.Errorf("failed to unmarshal payload: %w", asynq.SkipRetry)
	}

	// 從 DB 取得 user 的 email

	// 建立驗證信資料
	params := database.CreateVerifyEmailParams{
		Username: payload.Username,
		Email:    "elliott10009@gmail.com", /// DB 取得 user 的 email
		Secret:   util.RandomString(32),
	}
	verifyEmail, err := p.db.CreateVerifyEmail(ctx, params)
	if err != nil {
		return fmt.Errorf("failed to create verify email: %w", err)
	}

	subject := "Welcome to Simple Bank"
	verifyUrl := fmt.Sprintf(
		"http://simeple-bank.org/verify-email?id=%s&secret=%s", verifyEmail.ID, verifyEmail.Secret,
	)
	content := fmt.Sprintf(`Hello %s,<br/>
	Thank you for registering with us!<br/>
	Please <a href="%s">click here</a> to verify your email address.<br/>
	`, verifyEmail.Username, verifyUrl)
	to := []string{verifyEmail.Email}
	err = p.mailer.SendEmail(subject, content, to, nil, nil, nil)
	if err != nil {
		return fmt.Errorf("failed to send verify email: %w", err)
	}

	// log.Info().Str("type", task.Type()).Bytes("payload", task.Payload()).
	// 	Str("email", user.Email).Msg("processed task")
	return nil
}

package worker

import (
	"context"

	"github.com/hibiken/asynq"
	"github.com/redis/go-redis/v9"
	"github.com/simple_bank/database"
	"github.com/simple_bank/mail"
)

const (
	QueueCritical = "critical"
	QueueDefault  = "default"
)

type TaskProcessor interface {
	// Start starts the processing server and register handlers
	Start() error

	ProcessTaskSendVerifyEmail(
		ctx context.Context,
		task *asynq.Task,
	) error
}

type RedisTaskProcessor struct {
	server *asynq.Server
	db     database.Database
	mailer mail.EmailSender
}

func NewRedisTaskProcessor(
	redisOpt asynq.RedisClientOpt, db database.Database, mailer mail.EmailSender,
) TaskProcessor {
	logger := NewLogger()
	redis.SetLogger(logger)

	// controls parameters of the asynq Server, such as
	// - the maximum number of concurrent processing of tasks
	// - the retry delay for a failed task
	serverConfig := asynq.Config{
		Queues: map[string]int{
			// critical 佔 10/15 的處理資源，default 佔 5/15 的處理資源
			QueueCritical: 10,
			QueueDefault:  5,
		},
		ErrorHandler: asynq.ErrorHandlerFunc(func(ctx context.Context, task *asynq.Task, err error) {
			// 我們甚至可以傳送通知到 email, slack 等等，回報任務失敗了!
			// log.Error().Err(err).Str("type", task.Type()).
			// 		Bytes("payload", task.Payload()).Msg("process task failed")
		}),
		// 指定 Logger，讓 Asynq 執行遇到問題時的 log，可以符合我們想要的格式!
		Logger: logger,
	}
	server := asynq.NewServer(redisOpt, serverConfig)

	return &RedisTaskProcessor{
		server: server,
		db:     db,
		mailer: mailer,
	}
}

func (p *RedisTaskProcessor) Start() error {
	mux := asynq.NewServeMux()

	mux.HandleFunc(TaskSendVerifyEmail, p.ProcessTaskSendVerifyEmail)

	return p.server.Start(mux)
}

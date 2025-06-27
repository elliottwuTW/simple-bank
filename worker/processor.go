package worker

import (
	"context"

	"github.com/hibiken/asynq"
	"github.com/simple_bank/database"
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
}

func NewRedisTaskProcessor(redisOpt asynq.RedisClientOpt, db database.Database) TaskProcessor {
	// controls parameters of the asynq Server, such as
	// - the maximum number of concurrent processing of tasks
	// - the retry delay for a failed task
	serverConfig := asynq.Config{
		Queues: map[string]int{
			// critical 佔 10/15 的處理資源，default 佔 5/15 的處理資源
			QueueCritical: 10,
			QueueDefault:  5,
		},
	}
	server := asynq.NewServer(redisOpt, serverConfig)

	return &RedisTaskProcessor{
		server: server,
		db:     db,
	}
}

func (p *RedisTaskProcessor) Start() error {
	mux := asynq.NewServeMux()

	mux.HandleFunc(TaskSendVerifyEmail, p.ProcessTaskSendVerifyEmail)

	return p.server.Start(mux)
}

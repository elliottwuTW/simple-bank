package worker

import (
	"context"

	"github.com/hibiken/asynq"
	"github.com/simple_bank/database"
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
	serverConfig := asynq.Config{}
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

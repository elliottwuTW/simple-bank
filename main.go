package main

import (
	"context"
	"log"
	"time"

	"github.com/hibiken/asynq"
	"github.com/simple_bank/api"
	"github.com/simple_bank/config"
	"github.com/simple_bank/database"
	"github.com/simple_bank/mail"
	"github.com/simple_bank/worker"
)

func main() {
	cfg, err := config.LoadConfig(".")
	if err != nil {
		log.Fatal("cannot load config", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	db, err := database.New(ctx, cfg.DB)
	if err != nil {
		log.Fatal("cannot initialize database", err)
	}

	mailer := mail.NewGmailSender(
		cfg.Email.Name, cfg.Email.Address, cfg.Email.Password,
	)

	redisOpt := asynq.RedisClientOpt{Addr: cfg.Redis.Address}
	distributor := worker.NewRedisTaskDistributor(redisOpt)
	// 要使用 goroutine!!!!! 因為裡面 Start() 後，
	// the Asynq server will start and block and keep polling Redis for new tasks.
	go runTaskProcessor(redisOpt, db, mailer)

	server, err := api.NewServer(db, distributor, cfg)
	if err != nil {
		log.Fatal("cannot create server", err)
	}

	err = server.Start(cfg.ServerAddress)
	if err != nil {
		log.Fatal("cannot initialize server")
	}
}

func runTaskProcessor(
	cfg asynq.RedisClientOpt, db database.Database, mailer mail.EmailSender,
) {
	taskProcessor := worker.NewRedisTaskProcessor(cfg, db, mailer)
	// log.Info().Msg("start task processor")
	log.Println("start task processor")

	err := taskProcessor.Start()
	if err != nil {
		// log.Fatal().Err(err).Msg("failed to start task processor")
	}
}

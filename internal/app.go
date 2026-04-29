package internal

import (
	"os"
)

struct App {

}

func New() *App {
	return &App
}



func (a *App) Run() {
		
}


func (a *App) Run() int {
	ctx, cancel := context.WithCancel(context.Background())

	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	redisClient := redis.NewClient(&redis.Options{
		Addr:     a.config.Redis.Addr,
		Password: a.config.Redis.Password,
		DB:       a.config.Redis.DB,
		OnConnect: func(ctx context.Context, cn *redis.Conn) error {
			logger.Info("connected")
			return nil
		},
	})

	redisQueue := queue.NewRedisQueue(redisClient, logger)
	redisLocker := poller.NewRedisLocker(redisClient, logger)

	done := a.handleExit(logger)

	go func() {
		defer cancel()
		<-done
		logger.Info("executing cancel()")
	}()

	p := poller.New(time.Second, logger, redisQueue, redisLocker)

	logger.Info("starting app")
	p.Start(ctx)
	logger.Info("ending app")
	return 0
}

func (a *App) handleExit(logger *slog.Logger) <-chan struct{} {
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, os.Interrupt)
	done := make(chan struct{})
	go func() {
		<-sig
		logger.Info("Handling exit signal")
		done <- struct{}{}
	}()
	return done
}

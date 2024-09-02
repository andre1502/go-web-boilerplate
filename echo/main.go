package main

import (
	"boilerplate/server"
	"boilerplate/utils"
	"boilerplate/utils/config"
	"boilerplate/utils/logger"
	"context"
	"net/http"
	"os"
	"os/signal"
	"time"
)

func main() {
	if len(os.Args) < 3 {
		panic("args error")
	}

	etcdHost := os.Args[1]    // 127.0.0.1:2379
	settingName := os.Args[2] // /app_name/dev

	if utils.IsEmptyString(etcdHost) || utils.IsEmptyString(settingName) {
		panic("args error")
	}

	cfg := config.NewConfig(context.Background(), etcdHost, settingName)

	server := server.NewServer(cfg)

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	go func() {
		if err := server.Echo.Start(server.AppPort); err != http.ErrServerClosed {
			logger.Sugar.Fatal(err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server with a timeout of 30 seconds.
	<-ctx.Done()

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := server.Echo.Shutdown(ctx); err != nil {
		logger.Sugar.Fatal(err)
	}

	logger.Sugar.Info("Server exiting.")
}

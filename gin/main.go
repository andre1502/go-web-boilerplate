package main

import (
	"boilerplate/server"
	"boilerplate/utils"
	"boilerplate/utils/config"
	"boilerplate/utils/logger"
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/gin-contrib/graceful"
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

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	cfg := config.NewConfig(ctx, etcdHost, settingName)

	server := server.NewServer(cfg)

	grc, err := graceful.New(server.Gin, graceful.WithAddr(server.AppPort))
	if err != nil {
		panic(err)
	}

	defer grc.Close()

	go func() {
		if err := grc.RunWithContext(context.Background()); err != nil && err != context.Canceled {
			panic(err)
		}
	}()

	// Listen for the interrupt signal.
	<-ctx.Done()

	if err := grc.Shutdown(context.Background()); err != nil {
		panic(err)
	}

	logger.Sugar.Info("Server exiting.")
}

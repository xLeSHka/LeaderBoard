package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/xLeSHka/LeaderBoard/internal/config"
	"github.com/xLeSHka/LeaderBoard/internal/repository"
	"github.com/xLeSHka/LeaderBoard/internal/service"
	"github.com/xLeSHka/LeaderBoard/internal/transport/handlers"
	"github.com/xLeSHka/LeaderBoard/internal/transport/server"
	cache1 "github.com/xLeSHka/LeaderBoard/pkg/db/cache"
	"github.com/xLeSHka/LeaderBoard/pkg/logger"
)

func main() {
	ctx := context.Background()
	mainLogger := logger.New()
	cfg, err := config.New()
	if err != nil {
		log.Fatal(err)
	}
	redis := cache1.New(cfg.RedisConfig)
	repo := repository.New(redis)
	srv := service.New(repo)
	handlerSrv := handlers.New(srv, mainLogger)
	server, err := server.New(handlerSrv)
	graceCh := make(chan os.Signal, 1)
	signal.Notify(graceCh, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		if err := server.Start(cfg.ServerPort); err != nil {
			mainLogger.Error(err.Error())
		}
	}()
	<-graceCh
	if err := server.Stop(ctx); err != nil {
		mainLogger.Error(err.Error())
	}
	mainLogger.Info("server stopped")
}

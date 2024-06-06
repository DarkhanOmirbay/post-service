package app

import (
	"log/slog"
	grpcapp "post/internal/app/grpc"
	"post/internal/domain/storage"
	"post/internal/services/post"
	"time"
)

type App struct {
	GRPCServer *grpcapp.App
}

func New(log *slog.Logger, grpcPort int, dsn string, tokenTTL time.Duration) *App {
	postStorage, err := storage.NewPostStorage(dsn)
	if err != nil {
		panic(err)
	}
	postService := post.New(log, tokenTTL, postStorage)
	grpcApp := grpcapp.New(log, postService, grpcPort)

	return &App{GRPCServer: grpcApp}
}

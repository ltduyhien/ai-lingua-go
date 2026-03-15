// Package main is the entry point for the ai-lingua-go gRPC server.
package main

import (
	"context"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/ltduyhien/ai-lingua-go/config"
	"github.com/ltduyhien/ai-lingua-go/internal/cache"
	grpchandler "github.com/ltduyhien/ai-lingua-go/internal/grpc"
	"github.com/ltduyhien/ai-lingua-go/internal/translator"
	translationv1 "github.com/ltduyhien/ai-lingua-go/api/gen/translation/v1"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("config: %v", err)
	}
	ctx := context.Background()
	tr, err := translator.New(ctx, cfg.OllamaModel, cfg.OllamaBaseURL)
	if err != nil {
		log.Fatalf("translator: %v", err)
	}
	var cacheSvc grpchandler.Cache
	if cfg.RedisAddr != "" {
		c, err := cache.New(ctx, cfg.RedisAddr, cfg.RedisTTLSeconds)
		if err != nil {
			log.Fatalf("cache: %v", err)
		}
		cacheSvc = c
	}
	srv := grpc.NewServer()
	translationv1.RegisterTranslationServiceServer(srv, grpchandler.NewServer(tr, cacheSvc))
	reflection.Register(srv)
	lis, err := net.Listen("tcp", ":"+cfg.GRPCPort)
	if err != nil {
		log.Fatalf("listen: %v", err)
	}
	go func() {
		if err := srv.Serve(lis); err != nil {
			log.Fatalf("serve: %v", err)
		}
	}()
	log.Printf("server listening on :%s", cfg.GRPCPort)
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	srv.GracefulStop()
	log.Println("server stopped")
}

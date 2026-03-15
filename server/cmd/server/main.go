package main

import (
	"context"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/ltduyhien/ai-lingua-go/config"
	"github.com/ltduyhien/ai-lingua-go/internal/cache"
	grpchandler "github.com/ltduyhien/ai-lingua-go/internal/grpc"
	"github.com/ltduyhien/ai-lingua-go/internal/rest"
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
	grpcSrv := grpc.NewServer()
	translationServer := grpchandler.NewServer(tr, cacheSvc)
	translationv1.RegisterTranslationServiceServer(grpcSrv, translationServer)
	reflection.Register(grpcSrv)
	lis, err := net.Listen("tcp", ":"+cfg.GRPCPort)
	if err != nil {
		log.Fatalf("listen: %v", err)
	}
	go func() {
		if err := grpcSrv.Serve(lis); err != nil {
			log.Fatalf("serve gRPC: %v", err)
		}
	}()
	log.Printf("gRPC listening on :%s", cfg.GRPCPort)

	httpHandler := rest.NewHandler(translationServer)
	httpLis, err := net.Listen("tcp", ":"+cfg.HTTPPort)
	if err != nil {
		log.Fatalf("listen HTTP: %v", err)
	}
	go func() {
		if err := http.Serve(httpLis, httpHandler); err != nil {
			log.Fatalf("serve HTTP: %v", err)
		}
	}()
	log.Printf("HTTP (REST) listening on :%s", cfg.HTTPPort)
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	grpcSrv.GracefulStop()
	_ = httpLis.Close()
	log.Println("server stopped")
}

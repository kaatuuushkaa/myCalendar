package main

import (
	"context"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"
	"log"
	pb "myCalendar/grpc/pb"
	"myCalendar/internal/config"
	"myCalendar/internal/db"
	"myCalendar/internal/jwt"
	"myCalendar/internal/logger"
	"myCalendar/internal/middleware"
	"myCalendar/internal/user"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"
)

func main() {
	//database, err := db.InitDB()
	//if err != nil {
	//	log.Fatalf("Could not connect to DB: %v", err)
	//}
	//
	//key := os.Getenv("JWT_KEY")
	//if key == "" {
	//	log.Fatalf("JWT_KEY not set in .env")
	//}
	cfg, err := config.Load()
	if err != nil {
		log.Fatal("failed to load config", zap.Error(err))
	}

	log := logger.MustNew(cfg.IsDev)
	defer log.Sync()

	database, err := db.New(cfg.DB)
	if err != nil {
		log.Fatal("failed to connect to database", zap.Error(err))
	}
	sqlDB, _ := database.DB()
	defer sqlDB.Close()

	jwtService := jwt.New(cfg.JWT.Secret)
	repo := user.NewRepository(database)
	service := user.NewService(repo, log, jwtService)
	handler := user.NewHandler(service)

	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(middleware.AuthInterceptor(jwtService)))
	pb.RegisterUserServiceServer(grpcServer, handler)
	reflection.Register(grpcServer) //for curl

	lis, err := net.Listen("tcp", ":"+cfg.Server.GRPCPort)
	if err != nil {
		log.Fatal("Failed to listen: %v", zap.Error(err))
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	mux := runtime.NewServeMux(
		runtime.WithIncomingHeaderMatcher(func(key string) (string, bool) {
			if strings.ToLower(key) == "authorization" {
				return key, true
			}
			return runtime.DefaultHeaderMatcher(key)
		}))
	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	}

	err = pb.RegisterUserServiceHandlerFromEndpoint(ctx, mux, "localhost:"+cfg.Server.GRPCPort, opts)
	if err != nil {
		log.Fatal("Failed to register gateway: %v", zap.Error(err))
	}

	httpServer := &http.Server{
		Addr:    ":" + cfg.Server.HTTPPort,
		Handler: mux,
	}

	go func() {
		log.Info("gRPC server started", zap.String("port", cfg.Server.GRPCPort))
		if err := grpcServer.Serve(lis); err != nil {
			log.Error("gRPC server error", zap.Error(err))
		}
	}()

	go func() {
		log.Info("HTTP gateway started", zap.String("port", cfg.Server.HTTPPort))
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Error("HTTP server error", zap.Error(err))
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit

	log.Info("shutting down...")
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownCancel()

	grpcServer.GracefulStop()
	if err := httpServer.Shutdown(shutdownCtx); err != nil {
		log.Error("HTTP server shutdown error", zap.Error(err))
	}
	log.Info("server stopped")
}

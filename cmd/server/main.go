package main

import (
	"context"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/rs/cors"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/encoding/protojson"
	"log"
	pb "myCalendar/grpc/pb"
	"myCalendar/internal/config"
	"myCalendar/internal/jwt"
	"myCalendar/internal/logger"
	"myCalendar/internal/middleware"
	"myCalendar/internal/rpc"
	"myCalendar/internal/rpc/auth"
	"myCalendar/internal/rpc/create_event"
	"myCalendar/internal/rpc/create_user"
	"myCalendar/internal/rpc/delete_event"
	"myCalendar/internal/rpc/delete_user"
	"myCalendar/internal/rpc/get_event"
	"myCalendar/internal/rpc/get_user"
	"myCalendar/internal/rpc/get_user_events"
	"myCalendar/internal/rpc/health"
	"myCalendar/internal/rpc/reset_password"
	"myCalendar/internal/rpc/update_event"
	"myCalendar/internal/rpc/update_user"
	"myCalendar/internal/storage/pgrepo"
	"myCalendar/internal/storage/postgres"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatal("failed to load config", zap.Error(err))
	}

	log := logger.MustNew(cfg.IsDev)
	defer log.Sync()

	database, err := postgres.New(cfg.DB)
	if err != nil {
		log.Fatal("failed to connect to database", zap.Error(err))
	}
	sqlDB, _ := database.DB()
	defer sqlDB.Close()

	userRepo := pgrepo.NewUserRepo(database)
	jwtService := jwt.New(cfg.JWT.Secret)

	healthHandler := health.New()
	createUserHandler := create_user.New(userRepo, log)
	authHandler := auth.New(userRepo, log, jwtService)
	getUserHandler := get_user.New(userRepo, log)
	updateUserHandler := update_user.New(userRepo, log)
	deleteUserHandler := delete_user.New(userRepo, log)
	resetPasswordHandler := reset_password.New(userRepo, log)

	userServer := rpc.NewUserServer(
		healthHandler,
		createUserHandler,
		authHandler,
		getUserHandler,
		updateUserHandler,
		deleteUserHandler,
		resetPasswordHandler,
	)

	eventRepo := pgrepo.NewEventRepo(database)

	createEventHandler := create_event.New(eventRepo, log)
	getEventHandler := get_event.New(eventRepo, log)
	getUserEventsHandler := get_user_events.New(eventRepo, log)
	updateEventHandler := update_event.New(eventRepo, log)
	deleteEventHandler := delete_event.New(eventRepo, log)

	eventServer := rpc.NewEventServer(
		createEventHandler,
		getEventHandler,
		getUserEventsHandler,
		updateEventHandler,
		deleteEventHandler,
	)

	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(middleware.AuthInterceptor(jwtService)))
	pb.RegisterUserServiceServer(grpcServer, userServer)
	pb.RegisterEventServiceServer(grpcServer, eventServer)
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
		}),
		runtime.WithMarshalerOption(runtime.MIMEWildcard, &runtime.JSONPb{
			MarshalOptions: protojson.MarshalOptions{
				UseProtoNames: true,
			},
		}),
	)
	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	}

	err = pb.RegisterUserServiceHandlerFromEndpoint(ctx, mux, "localhost:"+cfg.Server.GRPCPort, opts)
	if err != nil {
		log.Fatal("Failed to register gateway: %v", zap.Error(err))
	}

	err = pb.RegisterEventServiceHandlerFromEndpoint(ctx, mux, "localhost:"+cfg.Server.GRPCPort, opts)
	if err != nil {
		log.Fatal("Failed to register event gateway: %v", zap.Error(err))
	}

	//cors middleware
	c := cors.New(cors.Options{
		AllowedOrigins: []string{
			"http://localhost:3000",
		},

		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Content-Type", "Authorization"},
		ExposedHeaders:   []string{"Content-Length"},
		AllowCredentials: true,
	})

	httpServer := &http.Server{
		Addr:    ":" + cfg.Server.HTTPPort,
		Handler: c.Handler(mux),
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

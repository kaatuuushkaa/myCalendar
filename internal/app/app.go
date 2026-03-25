package app

//
//import (
//	"context"
//	"log"
//	"net"
//	"net/http"
//	"os"
//	"os/signal"
//	"strings"
//	"syscall"
//	"time"
//
//	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
//	"github.com/rs/cors"
//	"go.uber.org/zap"
//	"google.golang.org/grpc"
//	"google.golang.org/grpc/credentials/insecure"
//	"google.golang.org/grpc/reflection"
//	"google.golang.org/protobuf/encoding/protojson"
//
//	pb "myCalendar/grpc/pb"
//	"myCalendar/internal/config"
//	"myCalendar/internal/jwt"
//	"myCalendar/internal/logger"
//	"myCalendar/internal/middleware"
//	"myCalendar/internal/rpc"
//	"myCalendar/internal/rpc/auth"
//	"myCalendar/internal/rpc/create_event"
//	"myCalendar/internal/rpc/create_user"
//	"myCalendar/internal/rpc/delete_event"
//	"myCalendar/internal/rpc/delete_user"
//	"myCalendar/internal/rpc/get_event"
//	"myCalendar/internal/rpc/get_user"
//	"myCalendar/internal/rpc/get_user_by_login"
//	"myCalendar/internal/rpc/get_user_events"
//	"myCalendar/internal/rpc/health"
//	"myCalendar/internal/rpc/reset_password"
//	"myCalendar/internal/rpc/update_event"
//	"myCalendar/internal/rpc/update_user"
//	"myCalendar/internal/storage/pgrepo"
//	"myCalendar/internal/storage/postgres"
//)
//
//// приложение со всеми зависимостями
//type App struct {
//	cfg        *config.Config
//	log        *zap.Logger
//	grpcServer *grpc.Server
//	httpServer *http.Server
//}
//
//// инициализирует все зависимости и возвращает App
//// main.go просто вызывает New() и Run()
//func New() (*App, error) {
//	cfg, err := config.Load()
//	if err != nil {
//		log.Fatal("failed to load config", zap.Error(err))
//	}
//
//	log := logger.MustNew(cfg.IsDev)
//	defer log.Sync()
//
//	database, err := postgres.New(cfg.DB)
//	if err != nil {
//		log.Fatal("failed to connect to database", zap.Error(err))
//	}
//	sqlDB, _ := database.DB()
//	defer sqlDB.Close()
//
//	userRepo := pgrepo.NewUserRepo(database)
//	eventRepo := pgrepo.NewEventRepo(database)
//
//	jwtService := jwt.New(cfg.JWT.Secret)
//
//	userServer := rpc.NewUserServer(
//		health.New(),
//		create_user.New(userRepo, log),
//		auth.New(userRepo, log, jwtService),
//		get_user.New(userRepo, log),
//		update_user.New(userRepo, log),
//		delete_user.New(userRepo, log),
//		reset_password.New(userRepo, log),
//		get_user_by_login.New(userRepo, log),
//	)
//
//	eventServer := rpc.NewEventServer(
//		create_event.New(eventRepo, log),
//		get_event.New(eventRepo, log),
//		get_user_events.New(eventRepo, log),
//		update_event.New(eventRepo, log),
//		delete_event.New(eventRepo, log),
//	)
//
//	grpcServer := grpc.NewServer(
//		grpc.UnaryInterceptor(middleware.AuthInterceptor(jwtService)))
//	pb.RegisterUserServiceServer(grpcServer, userServer)
//	pb.RegisterEventServiceServer(grpcServer, eventServer)
//	reflection.Register(grpcServer) //for curl
//
//	ctx, cancel := context.WithCancel(context.Background())
//	defer cancel()
//	mux := runtime.NewServeMux(
//		runtime.WithIncomingHeaderMatcher(func(key string) (string, bool) {
//			if strings.ToLower(key) == "authorization" {
//				return key, true
//			}
//			return runtime.DefaultHeaderMatcher(key)
//		}),
//		runtime.WithMarshalerOption(runtime.MIMEWildcard, &runtime.JSONPb{
//			MarshalOptions: protojson.MarshalOptions{
//				UseProtoNames: true,
//			},
//		}),
//	)
//
//	//grpcAddr := "localhost:" + cfg.Server.GRPCPort
//	grpcAddr := "0.0.0.0:" + cfg.Server.GRPCPort
//	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}
//
//	if err := pb.RegisterUserServiceHandlerFromEndpoint(ctx, mux, grpcAddr, opts); err != nil {
//		return nil, err
//	}
//	if err := pb.RegisterEventServiceHandlerFromEndpoint(ctx, mux, grpcAddr, opts); err != nil {
//		return nil, err
//	}
//
//	c := cors.New(cors.Options{
//		AllowedOrigins: []string{
//			"http://localhost:3000",
//		},
//
//		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
//		AllowedHeaders:   []string{"Content-Type", "Authorization"},
//		ExposedHeaders:   []string{"Content-Length"},
//		AllowCredentials: true,
//	})
//
//	httpServer := &http.Server{
//		Addr:    ":" + cfg.Server.HTTPPort,
//		Handler: c.Handler(mux),
//	}
//
//	return &App{
//		cfg:        cfg,
//		log:        log,
//		grpcServer: grpcServer,
//		httpServer: httpServer,
//	}, nil
//}
//
//// запускает серверы и ждёт сигнала остановки
//func (a *App) Run() {
//	// запускаем gRPC
//	lis, err := net.Listen("tcp", ":"+a.cfg.Server.GRPCPort)
//	if err != nil {
//		log.Fatal("Failed to listen: %v", zap.Error(err))
//	}
//
//	go func() {
//		a.log.Info("gRPC server started", zap.String("port", a.cfg.Server.GRPCPort))
//		if err := a.grpcServer.Serve(lis); err != nil {
//			a.log.Error("gRPC server error", zap.Error(err))
//		}
//	}()
//
//	go func() {
//		a.log.Info("HTTP gateway started", zap.String("port", a.cfg.Server.HTTPPort))
//		if err := a.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
//			a.log.Error("HTTP server error", zap.Error(err))
//		}
//	}()
//
//	// graceful shutdown
//	quit := make(chan os.Signal, 1)
//	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
//	<-quit
//
//	a.log.Info("shutting down...")
//	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
//	defer shutdownCancel()
//
//	a.grpcServer.GracefulStop()
//	if err := a.httpServer.Shutdown(shutdownCtx); err != nil {
//		a.log.Error("HTTP server shutdown error", zap.Error(err))
//	}
//	a.log.Info("server stopped")
//
//	a.log.Sync() //nolint:errcheck
//}

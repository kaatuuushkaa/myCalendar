package main

import (
	"context"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"
	"log"
	pb "myCalendar/grpc/pb"
	"myCalendar/internal/db"
	"myCalendar/internal/jwt"
	"myCalendar/internal/middleware"
	"myCalendar/internal/user"
	"net"
	"net/http"
	"os"
	"strings"
)

func main() {
	database, err := db.InitDB()
	if err != nil {
		log.Fatalf("Could not connect to DB: %v", err)
	}

	key := os.Getenv("JWT_KEY")
	if key == "" {
		log.Fatalf("JWT_KEY not set in .env")
	}
	jwtService := jwt.New(key)

	repo := user.NewRepository(database)
	service := user.NewService(repo)
	handler := user.NewHandler(service, jwtService)

	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(middleware.AuthInterceptor(jwtService)))
	pb.RegisterUserServiceServer(grpcServer, handler)
	reflection.Register(grpcServer) //for curl

	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	ctx := context.Background()
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

	err = pb.RegisterUserServiceHandlerFromEndpoint(ctx, mux, "localhost:50051", opts)
	if err != nil {
		log.Fatalf("Failed to register gateway: %v", err)
	}

	go func() {
		log.Println("gRPC server on :50051")
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatalf("grpc serve: %v", err)
		}
	}()

	log.Println("Http gateway on :8080")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatalf("http serve: %v", err)
	}
}

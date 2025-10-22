package main

import (
	"context"
	"database/sql"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "github.com/lib/pq"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	applicantsv1 "github.com/Thrun12/golang-assignment/api/proto/v1"
	"github.com/Thrun12/golang-assignment/internal/config"
	"github.com/Thrun12/golang-assignment/internal/db/sqlc"
	"github.com/Thrun12/golang-assignment/internal/middleware"
	"github.com/Thrun12/golang-assignment/internal/server"
	"github.com/Thrun12/golang-assignment/internal/service"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to load config: %v\n", err)
		os.Exit(1)
	}

	// Initialize logger
	log, err := zap.NewDevelopment()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to initialize logger: %v\n", err)
		os.Exit(1)
	}
	defer func() {
		_ = log.Sync()
	}()

	log.Info("starting job applicants API",
		zap.String("version", cfg.ServiceVersion),
		zap.String("environment", cfg.Environment),
		zap.Int("grpc_port", cfg.GRPCPort),
		zap.Int("http_port", cfg.ServerPort),
	)

	// Connect to database
	ctx := context.Background()
	db, err := sql.Open("postgres", cfg.DatabaseURL)
	if err != nil {
		log.Fatal("failed to connect to database",
			zap.Error(err),
		)
	}
	defer db.Close()

	// Configure connection pool for production use
	db.SetMaxOpenConns(25)                 // Maximum number of open connections to the database
	db.SetMaxIdleConns(5)                  // Maximum number of connections in the idle connection pool
	db.SetConnMaxLifetime(5 * time.Minute) // Maximum amount of time a connection may be reused
	db.SetConnMaxIdleTime(10 * time.Minute) // Maximum amount of time a connection may be idle

	// Test database connection
	if err := db.PingContext(ctx); err != nil {
		log.Fatal("failed to ping database",
			zap.Error(err),
		)
	}
	log.Info("successfully connected to database",
		zap.Int("max_open_conns", 25),
		zap.Int("max_idle_conns", 5),
	)

	// Initialize queries and service layers
	queries := sqlc.New(db)
	applicantService := service.NewApplicantService(queries, log)

	// Create gRPC server
	grpcServer := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			middleware.RecoveryInterceptor(log),
			middleware.UnaryServerInterceptor(log),
		),
	)

	// Register gRPC services - service layer implements the gRPC interface directly
	applicantsv1.RegisterApplicantsServiceServer(grpcServer, applicantService)

	// Enable gRPC reflection for debugging
	reflection.Register(grpcServer)

	// Start gRPC server in a goroutine
	grpcListener, err := net.Listen("tcp", fmt.Sprintf(":%d", cfg.GRPCPort))
	if err != nil {
		log.Fatal("failed to create gRPC listener",
			zap.Error(err),
		)
	}

	go func() {
		log.Info("starting gRPC server",
			zap.Int("port", cfg.GRPCPort),
			zap.String("address", fmt.Sprintf("localhost:%d", cfg.GRPCPort)),
		)
		if err := grpcServer.Serve(grpcListener); err != nil {
			log.Fatal("failed to serve gRPC",
				zap.Error(err),
			)
		}
	}()

	// Create HTTP gateway server
	gatewayCtx, gatewayCancel := context.WithCancel(ctx)
	defer gatewayCancel()

	grpcAddress := fmt.Sprintf("localhost:%d", cfg.GRPCPort)
	gatewayHandler, err := server.NewGatewayServer(gatewayCtx, grpcAddress, cfg.GetCORSOrigins(), db, log)
	if err != nil {
		log.Fatal("failed to create gateway server",
			zap.Error(err),
		)
	}

	// Create HTTP server
	httpServer := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.ServerPort),
		Handler:      gatewayHandler,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Start HTTP server in a goroutine
	go func() {
		log.Info("starting HTTP server",
			zap.Int("port", cfg.ServerPort),
			zap.String("docs_url", fmt.Sprintf("http://localhost:%d/docs/", cfg.ServerPort)),
		)
		log.Info("API documentation available",
			zap.String("url", fmt.Sprintf("http://localhost:%d/docs/", cfg.ServerPort)),
			zap.String("description", "Interactive API documentation and testing interface"),
		)
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal("failed to serve HTTP",
				zap.Error(err),
			)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the servers
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Info("shutting down servers...")

	// Graceful shutdown
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer shutdownCancel()

	// Shutdown HTTP server
	if err := httpServer.Shutdown(shutdownCtx); err != nil {
		log.Error("HTTP server shutdown error", zap.Error(err))
	}

	// Shutdown gRPC server
	grpcServer.GracefulStop()

	log.Info("servers stopped")
}

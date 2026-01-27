package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/projuktisheba/pse-api-v1/internal/config"
	"github.com/projuktisheba/pse-api-v1/internal/database"
	"github.com/projuktisheba/pse-api-v1/internal/handler"
	"github.com/projuktisheba/pse-api-v1/internal/repository"
	"github.com/projuktisheba/pse-api-v1/internal/routes"
	"github.com/projuktisheba/pse-api-v1/internal/service"
)

func main() {
	// Create a root context that is cancelled on SIGINT or SIGTERM
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	// Load application configuration (env, DB, JWT, server)
	cfg := config.LoadConfig()

	log.Printf("Config loaded | env=%s port=%d", cfg.App.Env, cfg.Server.Port)

	//--------------------------------
	// Database Connections
	//--------------------------------
	// 1. Connect to PostgreSQL using pgx pool
	pgPool, err := database.ConnectPostgres(cfg.DB.PostgresURL)
	if err != nil {
		log.Fatal("Postgres connect failed:", err)
	}
	defer pgPool.Close()
	log.Println("Postgres connected")

	// 2. Connect to Redis
	rdb := database.ConnectRedis(cfg.DB.RedisURL)
	defer rdb.Close()
	log.Println("Redis connected")

	// setup logger
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)


	//---------------------------------------------
	// Initialize repositories (data access layer)
	//---------------------------------------------
	dbRepository := repository.NewDBRepository(pgPool, rdb)

	//---------------------------------------------
	// Initialize services (business logic layer)
	//---------------------------------------------
	serviceRepository := service.NewServiceRepository(dbRepository, rdb, cfg)

	//---------------------------------------------
	// Initialize HTTP handlers (delivery layer)
	//---------------------------------------------
	handlers := handler.NewHandlerRepository(serviceRepository)

	// setup routes
	registeredRoutes := routes.Routes(cfg, rdb, handlers, infoLog, errorLog)

	// Setup HTTP server with timeouts
	addr := ":" + strconv.Itoa(cfg.Server.Port)
	srv := &http.Server{
		Addr:         addr,
		Handler:      registeredRoutes,
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
		IdleTimeout:  cfg.Server.IdleTimeout,
	}

	// Start HTTP server asynchronously
	go func() {
		log.Printf("%s v%s started on %s [%s]",
			cfg.App.Name,
			cfg.App.Version,
			addr,
			cfg.App.Env,
		)

		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal("Server error:", err)
		}
	}()

	// Wait for shutdown signal (Ctrl+C, SIGTERM)
	<-ctx.Done()
	log.Println("Shutdown signal received")

	// Graceful shutdown with timeout
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Println("Server shutdown failed:", err)
	}

	log.Println("Server stopped gracefully")
}

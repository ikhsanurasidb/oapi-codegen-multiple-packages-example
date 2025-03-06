package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/oapi-codegen-multiple-packages-example/config"
	gen_store "github.com/oapi-codegen-multiple-packages-example/internal/gen/store"
	store_handler "github.com/oapi-codegen-multiple-packages-example/internal/handler/store"
	store_repository "github.com/oapi-codegen-multiple-packages-example/internal/repository/store"
	store_service "github.com/oapi-codegen-multiple-packages-example/internal/service/store"

	"github.com/oapi-codegen-multiple-packages-example/pkg/mysql"
)

func main() {
	// Load configuration
	cfg := config.Get()

	// Initialize MySQL connection
	db, err := mysql.NewConnection(mysql.Config{
		Host:     cfg.MySQL.Host,
		Port:     cfg.MySQL.Port,
		User:     cfg.MySQL.User,
		Password: cfg.MySQL.Password,
		Database: cfg.MySQL.Database,
	})
	if err != nil {
		log.Fatalf("Failed to connect to MySQL: %v", err)
	}
	defer db.Close()

	repo := store_repository.NewRepository(db)

	svc := store_service.NewService(repo)

	storeHandler := store_handler.NewHandler(svc)

	router := gin.Default()

	gen_store.RegisterHandlers(router, storeHandler)

	srv := &http.Server{
		Addr:    ":8080",
		Handler: router,
	}

	// Graceful shutdown
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	// Give the server 5 seconds to finish ongoing requests
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exiting")
}

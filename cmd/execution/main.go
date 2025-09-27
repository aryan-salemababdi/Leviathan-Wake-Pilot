package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"leviathan/leviathan-wake-pilot/internal/config"
	"leviathan/leviathan-wake-pilot/internal/database"
	"leviathan/leviathan-wake-pilot/internal/exchange"
	"leviathan/leviathan-wake-pilot/internal/grpc_server"
	"leviathan/leviathan-wake-pilot/internal/service"

	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using default environment variables.")
	}

	log.Println("Initializing Execution & Risk Brain...")

	cfg := config.Load()

	dbClient, err := database.NewKeyDBClient(cfg.KeyDBAddress)
	if err != nil {
		log.Fatalf("FATAL: Could not connect to KeyDB: %v", err)
	}
	log.Println("Successfully connected to KeyDB.")

	exchangeClient := exchange.NewClient(cfg.ExchangeApiKey, cfg.ExchangeApiSecret)
	log.Println("Exchange client created.")

	executionSvc := service.NewExecutionService(cfg, dbClient, exchangeClient)

	grpcServer := grpc_server.NewGrpcServer(executionSvc)
	go func() {
		if err := grpcServer.Start(cfg.GrpcServerPort); err != nil {
			log.Fatalf("FATAL: Failed to start gRPC server: %v", err)
		}
	}()

	log.Println("Execution & Risk Brain is now RUNNING. Waiting for signals...")

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")
	time.Sleep(2 * time.Second)
	log.Println("Server gracefully stopped.")
}

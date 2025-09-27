package main

import (
	"leviathan/leviathan-wake-pilot/internal/config"
	"leviathan/leviathan-wake-pilot/internal/database"
	"leviathan/leviathan-wake-pilot/internal/exchange"
	"leviathan/leviathan-wake-pilot/internal/grpc_server"
	"leviathan/leviathan-wake-pilot/internal/service"
	"log"
)

func main() {
	log.Println("Initializing Execution & Risk Brain...")

	cfg, err := config.Load("config.json")
	if err != nil {
		log.Fatalf("FATAL: Could not load config: %v", err)
	}

	dbClient, err := database.NewKeyDBClient(cfg.KeyDBAddress)
	if err != nil {
		log.Fatalf("FATAL: Could not connect to KeyDB: %v", err)
	}

	exchangeClient := exchange.NewClient(cfg.ExchangeApiKey, cfg.ExchangeApiSecret)

	executionSvc := service.NewExecutionService(cfg, dbClient, exchangeClient)

	grpcServer := grpc_server.NewGrpcServer(executionSvc)
	go grpcServer.Start(cfg.GrpcServerPort)

	log.Println("Execution & Risk Brain is now RUNNING. Waiting for signals...")

}

package main

import (
	"leviathan/leviathan-wake-pilot/internal/config"
	"leviathan/leviathan-wake-pilot/internal/database"
	"leviathan/leviathan-wake-pilot/internal/exchange"
	"leviathan/leviathan-wake-pilot/internal/grpc_server"
	"leviathan/leviathan-wake-pilot/internal/service"
	"log"
	// ... import‌های دیگر
)

func main() {
	log.Println("Initializing Execution & Risk Brain...")

	// ۱. بارگذاری تنظیمات
	cfg, err := config.Load("config.json")
	if err != nil {
		log.Fatalf("FATAL: Could not load config: %v", err)
	}

	// ۲. اتصال به KeyDB
	dbClient, err := database.NewKeyDBClient(cfg.KeyDBAddress)
	if err != nil {
		log.Fatalf("FATAL: Could not connect to KeyDB: %v", err)
	}

	// ۳. ساخت کلاینت اتصال به صرافی
	exchangeClient := exchange.NewClient(cfg.ExchangeApiKey, cfg.ExchangeApiSecret)

	// ۴. ساخت سرویس اصلی با تزریق وابستگی‌ها
	executionSvc := service.NewExecutionService(cfg, dbClient, exchangeClient)

	// ۵. راه‌اندازی سرور gRPC در یک Goroutine
	grpcServer := grpc_server.NewGrpcServer(executionSvc)
	go grpcServer.Start(cfg.GrpcServerPort) // این تابع باید سرور gRPC را راه‌اندازی کند

	log.Println("Execution & Risk Brain is now RUNNING. Waiting for signals...")

	// ۶. مدیریت خاموش شدن امن
	// ... (کد مدیریت خاموش شدن مانند سرویس قبلی)
}

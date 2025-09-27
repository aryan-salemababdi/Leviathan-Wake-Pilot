package service

import (
	"context"
	"leviathan/leviathan-wake-pilot/internal/exchange"
	"log"

	"github.com/go-redis/redis"
	"honnef.co/go/tools/config"
)

type ExecutionService struct {
	cfg      *config.Config
	dbClient *redis.Client
	exClient *exchange.Client
}

func NewExecutionService(cfg *config.Config, db *redis.Client, ex *exchange.Client) *ExecutionService {
	return &ExecutionService{cfg: cfg, dbClient: db, exClient: ex}
}

func (s *ExecutionService) ProcessSignal(ctx context.Context, signal *proto.WhaleSignal) {
	log.Println("Processing new signal...")
}

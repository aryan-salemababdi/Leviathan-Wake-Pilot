package grpc_server

import (
	"context"
	"fmt"
	"log"
	"net"

	"leviathan/leviathan-wake-pilot/internal/service"
	"leviathan/leviathan-wake-pilot/proto"

	"google.golang.org/grpc"
)

type GrpcServer struct {
	proto.UnimplementedTradeSignalServer
	executionSvc *service.ExecutionService
}

func NewGrpcServer(execSvc *service.ExecutionService) *GrpcServer {
	return &GrpcServer{executionSvc: execSvc}
}

func (s *GrpcServer) SendSignal(ctx context.Context, signal *proto.WhaleSignal) (*proto.Ack, error) {
	log.Printf("Received WhaleSignal: %+v", signal)
	go s.executionSvc.ProcessSignal(ctx, signal)
	return &proto.Ack{Success: true}, nil
}

func (s *GrpcServer) Start(port string) error {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		return fmt.Errorf("failed to listen on port %s: %w", port, err)
	}

	grpcServer := grpc.NewServer()

	proto.RegisterTradeSignalServer(grpcServer, s)

	log.Printf("gRPC server is listening on %s...", port)

	if err := grpcServer.Serve(lis); err != nil {
		return fmt.Errorf("failed to serve gRPC: %w", err)
	}

	return nil
}

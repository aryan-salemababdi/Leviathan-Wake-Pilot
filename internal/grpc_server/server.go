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

// This struct implements the SignalService interface from the .proto file
type GrpcServer struct {
	proto.UnimplementedTradeSignalServer                           // برای سازگاری در آینده
	executionSvc                         *service.ExecutionService // سرویس اصلی اجرا
}

func NewGrpcServer(execSvc *service.ExecutionService) *GrpcServer {
	return &GrpcServer{executionSvc: execSvc}
}

// این متد زمانی که یک سیگنال جدید از Processor دریافت شود، اجرا می‌شود
func (s *GrpcServer) SendSignal(ctx context.Context, signal *proto.WhaleSignal) (*proto.Ack, error) {
	log.Printf("Received WhaleSignal: %+v", signal)
	go s.executionSvc.ProcessSignal(ctx, signal)
	return &proto.Ack{Success: true}, nil
}

// ⭐️ متد Start که فراموش شده بود
func (s *GrpcServer) Start(port string) error {
	// ۱. یک listener روی پورت مشخص شده ایجاد می‌کنیم
	lis, err := net.Listen("tcp", port)
	if err != nil {
		return fmt.Errorf("failed to listen on port %s: %w", port, err)
	}

	// ۲. یک نمونه جدید از سرور gRPC می‌سازیم
	grpcServer := grpc.NewServer()

	// ۳. سرویس خودمان را با سرور gRPC رجیستر می‌کنیم
	proto.RegisterTradeSignalServer(grpcServer, s)

	log.Printf("gRPC server is listening on %s...", port)

	// ۴. سرور را برای پذیرش درخواست‌های ورودی اجرا می‌کنیم
	if err := grpcServer.Serve(lis); err != nil {
		return fmt.Errorf("failed to serve gRPC: %w", err)
	}

	return nil
}

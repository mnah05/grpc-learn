package main

import (
	"context"
	"log"
	"net"
	"time"

	pb "invoiceService/pb/proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

// server is used to implement invoiceService.InvoiceServiceServer.
type server struct {
	pb.UnimplementedInvoiceServiceServer
}

// GenerateInvoice implements invoiceService.InvoiceServiceServer
func (s *server) GenerateInvoice(ctx context.Context, req *pb.InvoiceRequest) (*pb.InvoiceResponse, error) {
	log.Printf("Received invoice request for customer: %s", req.CustomerId)

	var totalAmount float64
	for _, item := range req.Items {
		totalAmount += float64(item.Quantity) * item.UnitPrice
	}

	// Generate a simple invoice ID (in real scenario, use UUID or database sequence)
	invoiceID := "INV-" + time.Now().Format("20060102150405")

	return &pb.InvoiceResponse{
		InvoiceId:     invoiceID,
		TotalAmount:   totalAmount,
		GeneratedAt:   time.Now().Format(time.RFC3339),
	}, nil
}

func main() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterInvoiceServiceServer(s, &server{})
	// Register reflection service on gRPC server.
	reflection.Register(s)
	log.Printf("Server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
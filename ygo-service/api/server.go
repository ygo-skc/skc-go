package api

import (
	"fmt"
	"log"
	"net"

	"github.com/ygo-skc/skc-go/ygo-service/db"
	"github.com/ygo-skc/skc-go/ygo-service/pb"
	"google.golang.org/grpc"
)

var (
	skcDBInterface db.SKCDatabaseAccessObject = db.SKCDAOImplementation{}
)

const (
	port = 9020
)

type Server struct {
	pb.CardServiceServer
}

func RunService() {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()

	// Register the service implementation with the server
	pb.RegisterCardServiceServer(grpcServer, &Server{})

	log.Printf("gRPC server is listening on port %d...", port)
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}

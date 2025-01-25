package api

import (
	"log"
	"net"

	"github.com/ygo-skc/skc-go/skc-db-service/db"
	"github.com/ygo-skc/skc-go/skc-db-service/pb"
	"google.golang.org/grpc"
)

var (
	skcDBInterface db.SKCDatabaseAccessObject = db.SKCDAOImplementation{}
)

type Server struct {
	pb.CardServiceServer
}

func RunService() {
	listener, err := net.Listen("tcp", ":9090")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()

	// Register the service implementation with the server
	pb.RegisterCardServiceServer(grpcServer, &Server{})

	log.Println("gRPC server is listening on port 9090...")
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}

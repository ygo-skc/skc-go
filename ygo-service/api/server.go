package api

import (
	"fmt"
	"log"
	"net"

	"github.com/ygo-skc/skc-go/common/util"
	"github.com/ygo-skc/skc-go/common/ygo"
	"github.com/ygo-skc/skc-go/ygo-service/db"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

var (
	skcDBInterface db.SKCDatabaseAccessObject = db.SKCDAOImplementation{}
)

const (
	port = 9020
)

type Server struct {
	ygo.CardServiceServer
}

func RunService() {
	// combine certs and create TLS creds
	util.CombineCerts("certs")
	if creds, err := credentials.NewServerTLSFromFile("certs/concatenated.crt", "certs/private.key"); err != nil {
		log.Fatalf("Unable to create TLS credentials: %v", err)
	} else {
		// Register the service implementation with the server
		grpcServer := grpc.NewServer(grpc.Creds(creds))
		ygo.RegisterCardServiceServer(grpcServer, &Server{})

		log.Printf("Starting gRPC service on port %d...", port)
		listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
		if err != nil {
			log.Fatalf("Failed to listen: %v", err)
		}
		if err := grpcServer.Serve(listener); err != nil {
			log.Fatalf("Failed to serve: %v", err)
		}
	}
}

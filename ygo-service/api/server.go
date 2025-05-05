package api

import (
	"fmt"
	"log"
	"net"
	"time"

	"github.com/ygo-skc/skc-go/common/util"
	"github.com/ygo-skc/skc-go/common/ygo"
	"github.com/ygo-skc/skc-go/ygo-service/db"
	"github.com/ygo-skc/skc-go/ygo-service/health"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/keepalive"
)

var (
	skcDBInterface db.SKCDatabaseAccessObject = db.SKCDAOImplementation{}
)

const (
	port = 9020
)

type healthServiceServer struct {
	health.HealthServiceServer
}

type ygoServiceServer struct {
	ygo.CardServiceServer
}

func RunService() {
	// combine certs and create TLS creds
	util.CombineCerts("certs")
	if creds, err := credentials.NewServerTLSFromFile("certs/concatenated.crt", "certs/private.key"); err != nil {
		log.Fatalf("Unable to create TLS credentials: %v", err)
	} else {
		// Register the service implementation with the server
		grpcServer := grpc.NewServer(
			grpc.Creds(creds),
			grpc.KeepaliveParams(keepalive.ServerParameters{
				MaxConnectionIdle:     5 * time.Minute,
				MaxConnectionAge:      60 * time.Minute,
				MaxConnectionAgeGrace: 5 * time.Minute,
				Time:                  1 * time.Minute,
			}))
		health.RegisterHealthServiceServer(grpcServer, &healthServiceServer{})
		ygo.RegisterCardServiceServer(grpcServer, &ygoServiceServer{})

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

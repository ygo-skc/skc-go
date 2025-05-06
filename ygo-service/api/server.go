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
			grpc.MaxConcurrentStreams(200),
			grpc.KeepaliveParams(keepalive.ServerParameters{
				MaxConnectionIdle:     30 * time.Minute, // how long a connection can last while idle
				MaxConnectionAge:      4 * time.Hour,    // total time a connection can live for before killed
				MaxConnectionAgeGrace: 15 * time.Second, // time after MaxConnectionAge where connection can finish work
				Time:                  30 * time.Second, // how often to ping client
				Timeout:               1 * time.Second,  // how fast ping should be
			}),
			grpc.KeepaliveEnforcementPolicy(keepalive.EnforcementPolicy{
				MinTime:             30 * time.Second, // prevents clients from sending pings too often
				PermitWithoutStream: true,             // allow pings when no active RPC
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

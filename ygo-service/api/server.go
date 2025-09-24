package api

import (
	"fmt"
	"log"
	"net"
	"time"

	"github.com/ygo-skc/skc-go/common/health"
	"github.com/ygo-skc/skc-go/common/util"
	"github.com/ygo-skc/skc-go/common/ygo"
	"github.com/ygo-skc/skc-go/ygo-service/db"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/keepalive"
)

var (
	cardRepo    db.CardRepository    = db.YGOCardRepository{}
	productRepo db.ProductRepository = db.YGOProductRepository{}
	scoreRepo   db.ScoreRepository   = db.YGOScoreRepository{}
)

const (
	port = 9020
)

type healthServiceServer struct {
	health.HealthServiceServer
}

type ygoCardServiceServer struct {
	ygo.CardServiceServer
}

type ygoProductServiceServer struct {
	ygo.ProductServiceServer
}

type ygoScoreServiceServer struct {
	ygo.ScoreServiceServer
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
			grpc.MaxConcurrentStreams(1024),
			grpc.KeepaliveParams(keepalive.ServerParameters{
				MaxConnectionIdle:     15 * time.Minute,       // how long a connection can last while idle
				MaxConnectionAge:      40 * time.Minute,       // total time a connection can live for before killed
				MaxConnectionAgeGrace: 15 * time.Second,       // time after MaxConnectionAge where connection can finish work
				Time:                  45 * time.Second,       // how often to ping client
				Timeout:               300 * time.Millisecond, // how fast ping should be
			}),
			grpc.KeepaliveEnforcementPolicy(keepalive.EnforcementPolicy{
				MinTime:             45 * time.Second, // prevents clients from sending pings too often
				PermitWithoutStream: true,             // allow pings when no active RPC
			}),
			grpc.ConnectionTimeout(50*time.Millisecond),
			// below are experimental
			grpc.NumStreamWorkers(128),
			grpc.SharedWriteBuffer(true),
		)

		// register services
		health.RegisterHealthServiceServer(grpcServer, &healthServiceServer{})
		ygo.RegisterCardServiceServer(grpcServer, &ygoCardServiceServer{})
		ygo.RegisterProductServiceServer(grpcServer, &ygoProductServiceServer{})
		ygo.RegisterScoreServiceServer(grpcServer, &ygoScoreServiceServer{})

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

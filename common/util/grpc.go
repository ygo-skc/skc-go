package util

import (
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/keepalive"
)

func NewServer(creds credentials.TransportCredentials) *grpc.Server {
	return grpc.NewServer(
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
}

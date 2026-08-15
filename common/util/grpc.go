package util

import (
	"runtime"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/keepalive"
)

func baselineServerOptions(creds credentials.TransportCredentials) []grpc.ServerOption {
	return []grpc.ServerOption{
		grpc.Creds(creds),
		grpc.MaxConcurrentStreams(1024),

		grpc.ReadBufferSize(64 << 10),
		grpc.WriteBufferSize(64 << 10),
		grpc.InitialWindowSize(256 << 10),   // per stream setting
		grpc.InitialConnWindowSize(4 << 20), // this controls how much data is sent for all streams in a connection

		grpc.KeepaliveParams(keepalive.ServerParameters{
			MaxConnectionIdle:     1 * time.Minute,  // how long a connection can last while idle
			MaxConnectionAge:      15 * time.Minute, // total time a connection can live for before killed
			MaxConnectionAgeGrace: 15 * time.Second, // time after MaxConnectionAge where connection can finish work
			Time:                  20 * time.Second, // how often to ping client
			Timeout:               3 * time.Second,  // how fast ping should be
		}),
		grpc.KeepaliveEnforcementPolicy(keepalive.EnforcementPolicy{
			MinTime:             15 * time.Second, // prevents clients from sending pings too often
			PermitWithoutStream: false,            // allow pings when no active RPC - if true conn will probably never close...
		}),
		grpc.ConnectionTimeout(4 * time.Second), // deadline covering the TLS handshake - grpc defaults to 120s

		// below are experimental
		grpc.NumStreamWorkers(uint32(runtime.GOMAXPROCS(0))),
	}
}

// NewServer builds a gRPC server with the shared SKC tuning, leaving grpc-go's default
// message size limits in place. Services that send or receive anything unusual should
// use NewServerWithOptions instead.
func NewServer(creds credentials.TransportCredentials) *grpc.Server {
	return NewServerWithOptions(creds)
}

// NewServerWithOptions builds a gRPC server with the shared SKC tuning, then applies opts
// on top. grpc-go applies server options in order, so anything passed here wins over the
// baseline.
func NewServerWithOptions(creds credentials.TransportCredentials, opts ...grpc.ServerOption) *grpc.Server {
	return grpc.NewServer(append(baselineServerOptions(creds), opts...)...)
}

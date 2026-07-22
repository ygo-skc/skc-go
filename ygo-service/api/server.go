package api

import (
	"fmt"
	"log/slog"
	"net"
	"os"
	"runtime"
	"time"

	"github.com/ygo-skc/skc-go/common/v3/health"
	"github.com/ygo-skc/skc-go/common/v3/util"
	"github.com/ygo-skc/skc-go/common/v3/ygo"
	"github.com/ygo-skc/skc-go/ygo-service/db"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/keepalive"
)

var (
	chicagoLocation *time.Location
)

func init() {
	if location, err := time.LoadLocation("America/Chicago"); err != nil {
		slog.Error("Failed to load Chicago location", slog.Any("err", err))
		os.Exit(1)
	} else {
		chicagoLocation = location
	}
}

var (
	cardRepo            db.CardRepository            = db.YGOCardRepository{}
	productRepo         db.ProductRepository         = db.YGOProductRepository{}
	cardRestrictionRepo db.CardRestrictionRepository = db.YGOCardRestrictionRepository{}
	scoreRepo           db.ScoreRepository           = db.YGOScoreRepository{}
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

type ygoCardRestrictionServiceServer struct {
	ygo.CardRestrictionServiceServer
}

type ygoScoreServiceServer struct {
	ygo.ScoreServiceServer
}

func RunService() {
	util.CombineCerts("certs")
	if creds, err := credentials.NewServerTLSFromFile("certs/concatenated.crt", "certs/private.key"); err != nil {
		slog.Error("Unable to create TLS credentials", slog.Any("err", err))
		os.Exit(1)
	} else {
		grpcServer := grpc.NewServer(
			grpc.Creds(creds),
			grpc.MaxConcurrentStreams(1024),

			grpc.ReadBufferSize(256<<10),
			grpc.WriteBufferSize(256<<10),
			grpc.InitialWindowSize(512<<10),        // per stream setting
			grpc.InitialConnWindowSize(5*512*1024), // if above 3 change, change the middle value here too - this controls how much data is sent for all streams in a connection

			grpc.KeepaliveParams(keepalive.ServerParameters{
				MaxConnectionIdle:     1 * time.Minute,  // how long a connection can last while idle
				MaxConnectionAge:      15 * time.Minute, // total time a connection can live for before killed
				MaxConnectionAgeGrace: 15 * time.Second, // time after MaxConnectionAge where connection can finish work
				Time:                  15 * time.Second, // how often to ping client
				Timeout:               3 * time.Second,  // how fast ping should be
			}),
			grpc.KeepaliveEnforcementPolicy(keepalive.EnforcementPolicy{
				MinTime:             15 * time.Second, // prevents clients from sending pings too often
				PermitWithoutStream: false,            // allow pings when no active RPC - if true conn will probably never close...
			}),
			grpc.ConnectionTimeout(4*time.Second),

			grpc.NumStreamWorkers(uint32(runtime.GOMAXPROCS(0))),
			grpc.SharedWriteBuffer(true),

			grpc.MaxRecvMsgSize(200<<10),
			grpc.MaxSendMsgSize(2<<20),
		)

		health.RegisterHealthServiceServer(grpcServer, &healthServiceServer{})
		ygo.RegisterCardServiceServer(grpcServer, &ygoCardServiceServer{})
		ygo.RegisterProductServiceServer(grpcServer, &ygoProductServiceServer{})
		ygo.RegisterCardRestrictionServiceServer(grpcServer, &ygoCardRestrictionServiceServer{})
		ygo.RegisterScoreServiceServer(grpcServer, &ygoScoreServiceServer{})

		slog.Info("Starting gRPC service", slog.Int("port", port))
		listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
		if err != nil {
			slog.Error("Failed to listen", slog.Any("err", err))
			os.Exit(1)
		}
		if err := grpcServer.Serve(listener); err != nil {
			slog.Error("Failed to serve grpc", slog.Any("err", err))
			os.Exit(1)
		}
	}
}

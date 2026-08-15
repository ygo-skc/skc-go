package api

import (
	"fmt"
	"log/slog"
	"net"
	"os"
	"time"

	"github.com/ygo-skc/skc-go/common/v3/health"
	"github.com/ygo-skc/skc-go/common/v3/util"
	"github.com/ygo-skc/skc-go/common/v3/ygo"
	"github.com/ygo-skc/skc-go/ygo-service/db"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
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
	creds, err := credentials.NewServerTLSFromFile("certs/concatenated.crt", "certs/private.key")
	if err != nil {
		slog.Error("Unable to create TLS credentials", slog.Any("err", err))
		os.Exit(1)
	}

	// shared tuning (streams, buffers, windows, keepalive) lives in util.NewServerWithOptions
	grpcServer := util.NewServerWithOptions(creds,
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

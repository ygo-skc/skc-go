package main

import (
	"context"
	"log"
	"net"
	"net/http"

	pb "github.com/ygo-skc/skc-go/skc-db/api"
	"github.com/ygo-skc/skc-go/skc-db/db"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

var (
	skcDBInterface db.SKCDatabaseAccessObject = db.SKCDAOImplementation{}
)

type Server struct {
	pb.CardServiceServer
}

func (s *Server) GetYGOCard(ctx context.Context, req *pb.YGOCardRequest) (*pb.YGOCardResponse, error) {
	if c, err := skcDBInterface.GetDesiredCardInDBUsingID(ctx, req.CardID); err != nil && err.StatusCode == http.StatusNotFound {
		return nil, status.Errorf(codes.NotFound, "%s", err.Message)
	} else if err != nil {
		return nil, status.Errorf(codes.Internal, "%s", err.Message)
	} else {
		res := &pb.YGOCardResponse{
			CardID:        c.CardID,
			CardColor:     c.CardColor,
			CardName:      c.CardName,
			CardAttribute: c.CardAttribute,
			CardEffect:    c.CardEffect,
		}

		if c.MonsterType != nil {
			res.MonsterType = wrapperspb.String(*c.MonsterType)
		}

		if c.MonsterAttack != nil {
			res.MonsterAttack = wrapperspb.UInt32(uint32(*c.MonsterAttack))
		}

		if c.MonsterDefense != nil {
			res.MonsterDefense = wrapperspb.UInt32(uint32(*c.MonsterDefense))
		}

		return res, nil
	}
}

func listen() {
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

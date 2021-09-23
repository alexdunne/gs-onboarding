package api

import (
	"fmt"
	"net"

	pb "github.com/alexdunne/gs-onboarding/internal/api/protobufs"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

type server struct {
	port   int
	srv    pb.APIServer
	logger *zap.Logger
}

func NewServer(port int, logger *zap.Logger, srv pb.APIServer) *server {
	return &server{
		port:   port,
		srv:    srv,
		logger: logger,
	}
}

func (s *server) Start() error {
	s.logger.Info(fmt.Sprintf("starting server on port %d", s.port))

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", s.port))
	if err != nil {
		return errors.Wrap(err, "failed to listen")
	}

	gs := grpc.NewServer()
	pb.RegisterAPIServer(gs, s.srv)
	if err := gs.Serve(lis); err != nil {
		return errors.Wrap(err, "failed to serve")
	}

	return nil
}

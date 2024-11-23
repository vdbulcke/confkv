package server

import (
	"github.com/vdbulcke/confkv/src/assert"
	"github.com/vdbulcke/confkv/src/controller"
	"github.com/vdbulcke/confkv/src/logger"
	"github.com/vdbulcke/confkv/src/pb"
	"github.com/vdbulcke/confkv/src/scheduler"
	"go.uber.org/zap"
)

const (
	DefaultGRPCPort       = 5000
	DefaultMonitoringPort = 9696
)

type OptionFunc func(s *Server)

func WithLogger(l *zap.Logger) OptionFunc {
	return func(s *Server) {
		s.logger = l
	}
}

func WithGrpcPort(p int) OptionFunc {
	return func(s *Server) {
		s.grpcPort = p
	}
}

func WithMonitoringPort(p int) OptionFunc {
	return func(s *Server) {
		s.monitoringPort = p
	}
}

type Server struct {
	pb.UnimplementedConfKVServer

	ctrl           *controller.Controller
	sched          *scheduler.Scheduler
	monitoringPort int
	grpcPort       int

	logger *zap.Logger
}

func NewServer(ctrl *controller.Controller, s *scheduler.Scheduler, opts ...OptionFunc) (*Server, error) {
	assert.NotNil(ctrl, assert.Exit)
	assert.NotNil(s, assert.Exit)

	srv := &Server{
		ctrl:           ctrl,
		sched:          s,
		monitoringPort: DefaultMonitoringPort,
		grpcPort:       DefaultGRPCPort,
		logger:         logger.GetLogger(false),
	}

	for _, fn := range opts {
		fn(srv)
	}

	return srv, nil
}

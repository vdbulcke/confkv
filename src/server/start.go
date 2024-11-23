package server

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/recovery"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/vdbulcke/confkv/src/pb"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
)

func (s *Server) panicFunc(p any) error {
	return status.Errorf(codes.Unknown, "panic triggered: %v", p)
}

func (s *Server) StartBlocking() error {

	s.logger.Info("starting scheduler")
	s.sched.Start()

	//
	// Start Prometheus
	//
	r := http.NewServeMux()
	r.Handle("/metrics", promhttp.Handler())
	promSrv := &http.Server{
		Addr:    fmt.Sprintf(":%d", s.monitoringPort),
		Handler: r,
	}

	s.logger.Info("starting prometheus server", zap.Int("port", s.monitoringPort))

	// Initializing the server in a goroutine so that
	// it won't block the graceful shutdown handling below
	go func() {
		if err := promSrv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			s.logger.Error("Error starting prometheus server", zap.Error(err))
		}
	}()

	//
	// Start GRPC
	//

	srv := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			recovery.UnaryServerInterceptor(recovery.WithRecoveryHandler(s.panicFunc)),
		),
	)
	grpcAddr := fmt.Sprintf(":%d", s.grpcPort)
	ln, err := net.Listen("tcp", grpcAddr)
	if err != nil {
		return err
	}

	reflection.Register(srv)
	pb.RegisterConfKVServer(srv, s)

	go func() {
		s.logger.Info("starting grpc server", zap.Int("port", s.grpcPort))
		if err := srv.Serve(ln); err != nil && !errors.Is(err, grpc.ErrServerStopped) {
			s.logger.Error("Error starting grpc server", zap.Error(err))
		}
	}()

	//
	// Handle Graceful shutdown
	//

	// Wait for interrupt signal to gracefully shutdown the server with
	// a timeout of 5 seconds.
	quit := make(chan os.Signal, 1)
	// kill (no param) default send syscall.SIGTERM
	// kill -2 is syscall.SIGINT
	// kill -9 is syscall.SIGKILL but can't be catch, so don't need add it
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	s.logger.Info("Shutting down server...")

	// The context is used to inform the server it has 5 seconds to finish
	// the request it is currently handling
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	go s.sched.Shutdown(ctx)
	srv.GracefulStop()

	if err := s.ctrl.Close(); err != nil {
		s.logger.Error("Error controller shutdown", zap.Error(err))
	}

	if err := promSrv.Shutdown(ctx); err != nil {
		s.logger.Error("Error forced to shutdown", zap.Error(err))

	}

	s.logger.Info("Server exiting")

	return nil
}

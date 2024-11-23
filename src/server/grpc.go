package server

import (
	"context"
	"fmt"

	"github.com/vdbulcke/confkv/src/assert"
	"github.com/vdbulcke/confkv/src/logger"
	"github.com/vdbulcke/confkv/src/pb"
	"go.uber.org/zap"
	"google.golang.org/grpc/metadata"
)

func (s *Server) Get(ctx context.Context, in *pb.GetMessage) (*pb.GetResponse, error) {
	assert.NotNil(in, assert.SIGTERM, "missing required grpc GetMessage")

	traceID := getTraceIDFromMetadata(ctx)
	reqCtx := logger.ContextWithTraceID(context.Background(), traceID)
	l := logger.GetLoggerWithContext(reqCtx, s.logger)

	bkt := in.Bucket
	outcome := "failed"
	defer func() {

		l.Debug("GRPC Get", zap.String("bucket", bkt), zap.String("status", outcome))
	}()

	data, err := s.ctrl.Get(reqCtx, bkt, in.Key)
	if err != nil {

		l.Error("Grpc error get key",
			zap.String("bucket", bkt),
			zap.String("key", in.Key),
			zap.Error(err),
		)

		outcome = "failed"
		return nil, fmt.Errorf("GRPC err: %w", err)

	}

	outcome = "success"
	return &pb.GetResponse{
		Value: data,
	}, nil

}

func getTraceIDFromMetadata(ctx context.Context) string {

	traceId := logger.NewTraceID()

	if md, ok := metadata.FromIncomingContext(ctx); ok {

		traces := md["x-trace-id"]
		if len(traces) > 0 {
			traceId = traces[0]
		}
	}

	return traceId
}

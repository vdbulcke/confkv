package controller

import (
	"context"

	"github.com/vdbulcke/confkv/src/logger"
	"go.uber.org/zap"
)

type DirSource struct {
	BucketName string
	Dir        string

	logger *zap.Logger
}

func (s *DirSource) Sync(ctx context.Context) (string, string, error) {
	s.logger.Debug("Dir Sync", logger.GetZapTraceID(ctx))
	return s.Dir, s.BucketName, nil
}

func (s *DirSource) GetBucket() string {
	return s.BucketName
}

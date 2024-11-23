package controller

import (
	"context"
	"fmt"

	"github.com/vdbulcke/confkv/src/config"
	"github.com/vdbulcke/confkv/src/logger"
	"github.com/vdbulcke/confkv/src/storage"
	"go.uber.org/zap"
)

type BucketSource interface {
	// Sync returns (local dir path, bucketname , error)
	Sync(ctx context.Context) (string, string, error)
	GetBucket() string
}

type Controller struct {
	sources []BucketSource

	store  *storage.KVStore
	logger *zap.Logger
}

type OptionFunc func(c *Controller)

func WithLogger(l *zap.Logger) OptionFunc {
	return func(c *Controller) {
		c.logger = l
	}
}

func NewController(db *storage.KVStore, bucketSources map[string]config.Bucket, options ...OptionFunc) (*Controller, error) {

	if db == nil {
		return nil, fmt.Errorf("'db *storage.KVStore' cannot be nil")
	}

	l := logger.GetLogger(false)

	ctrl := &Controller{
		store:  db,
		logger: l,
	}

	for _, fn := range options {
		fn(ctrl)
	}

	size := len(bucketSources)
	sources := make([]BucketSource, 0, size)

	for name, v := range bucketSources {

		switch v.Source {
		case config.SourceDir:

			b := &DirSource{
				BucketName: name,
				Dir:        v.LocalDir,
				logger:     ctrl.logger,
			}

			sources = append(sources, b)

		case config.SourceGit:
			b := &GitSource{
				BucketName:      name,
				LocalDir:        v.LocalDir,
				GitRepo:         v.GitRepoUrl,
				GitRelativePath: v.GitRelativeDir,
				GitUsername:     v.GitUsername,
				GitPasswordEnv:  v.GitPasswordEnvVar,
				logger:          ctrl.logger,
			}

			sources = append(sources, b)
		default:
			return nil, fmt.Errorf("unsupported source %s for bucket: %s", v.Source, name)
		}
	}

	ctrl.sources = sources
	return ctrl, nil
}

func (c *Controller) Close() error {
	return c.store.Close()
}

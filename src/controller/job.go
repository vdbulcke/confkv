package controller

import (
	"context"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/vdbulcke/confkv/src/assert"
	"github.com/vdbulcke/confkv/src/logger"
	"go.uber.org/zap"
)

func (c *Controller) SyncJob(ctx context.Context) error {

	l := logger.GetLoggerWithContext(ctx, c.logger)
	failedBucket := []string{}

	for _, s := range c.sources {
		err := c.SyncSource(ctx, s)
		if err != nil {
			l.Error("error synchronizing source", zap.Error(err))
			failedBucket = append(failedBucket, s.GetBucket())
		}
	}

	if len(failedBucket) > 0 {
		return fmt.Errorf("SyncJobs errors for buckets '%s'", strings.Join(failedBucket, ","))
	}

	return nil
}

func (c *Controller) SyncSource(ctx context.Context, s BucketSource) error {
	assert.NotNil(s, assert.Panic, "SynSource BucketSource is nil")
	l := logger.GetLoggerWithContext(ctx, c.logger)

	outcome := "failed"
	defer func() {
		b := s.GetBucket()
		l.Info("SyncSource ",
			zap.String("status", outcome),
			zap.String("bucket", b),
		)
	}()

	dir, bucketName, err := s.Sync(ctx)
	if err != nil {

		return err
	}

	err = c.SaveDir(ctx, filepath.Clean(dir), bucketName)
	if err != nil {
		return err
	}

	outcome = "success"
	return nil
}

func (c *Controller) SaveDir(ctx context.Context, dir, bucketName string) error {
	assert.NotNil(c.store, assert.SIGTERM, "mussing required controller.store")

	return filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
		c.logger.Debug("WalkDir process", zap.String("file", path))
		// I/O error when reading dir
		if err != nil {
			return err
		}

		if d.IsDir() {
			if strings.HasSuffix(d.Name(), ".git") {
				return fs.SkipDir
			}
			// don't process nested dir
			// WalkDir will process them later
			return nil
		}

		// TODO: add file extension filtering

		f, err := os.Open(path)
		if err != nil {
			c.logger.Error("WalkDir opening file", zap.String("file", path), zap.Error(err))
			return nil // continute processing other files
		}
		defer f.Close()

		data, err := io.ReadAll(f)
		if err != nil {
			c.logger.Error("WalkDir reading file", zap.String("file", path), zap.Error(err))
			return nil // continute processing other files
		}

		key := strings.TrimPrefix(path, dir)
		key = strings.TrimPrefix(key, "/")

		err = c.store.Put(ctx, bucketName, key, data)
		if err != nil {
			c.logger.Error("WalkDir saving file data to bolt", zap.String("file", path), zap.Error(err))
			return nil // continute processing other files
		}

		return nil
	})
}

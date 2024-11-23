package controller

import (
	"context"
	"errors"
	"os"
	"path/filepath"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
	"github.com/vdbulcke/confkv/src/assert"
	"github.com/vdbulcke/confkv/src/logger"
	"go.uber.org/zap"
)

type GitSource struct {
	BucketName string

	LocalDir        string
	GitRepo         string
	GitRelativePath string
	GitUsername     string
	GitPasswordEnv  string

	logger *zap.Logger
}

func (s *GitSource) WithLogger(l *zap.Logger) {
	if l == nil {
		l = logger.GetLogger(true)
	}

	s.logger = l

}

func (s *GitSource) GetBucket() string {
	return s.BucketName
}

func (s *GitSource) Sync(ctx context.Context) (string, string, error) {

	l := logger.GetLoggerWithContext(ctx, s.logger)
	l = l.With(zap.String("bucket", s.BucketName))
	r, err := git.PlainOpen(s.LocalDir)
	if err != nil {

		if !errors.Is(err, git.ErrRepositoryNotExists) {
			return "", "", err
		}
		l.Debug("repo does not exist, clone fresh repo")
		r, err := git.PlainClone(s.LocalDir, false, &git.CloneOptions{
			// The intended use of a GitHub personal access token is in replace of your password
			// because access tokens can easily be revoked.
			// https://help.github.com/articles/creating-a-personal-access-token-for-the-command-line/
			Auth: &http.BasicAuth{
				Username: s.GitUsername, // yes, this can be anything except an empty string
				Password: os.Getenv(s.GitPasswordEnv),
			},
			URL: s.GitRepo,
			// Progress: os.Stdout,
		})
		if err != nil {
			return "", "", err
		}
		ref, err := r.Head()
		assert.ErrNotNil(err, assert.Panic, "cannot get HEAD from git repo")
		l.Info("local git repo", zap.String("HEAD", ref.Name().Short()), zap.String("hash", ref.Hash().String()))

		return filepath.Join(s.LocalDir, s.GitRelativePath), s.BucketName, nil
	}

	wt, err := r.Worktree()
	if err != nil {
		return "", "", err

	}

	l.Debug("pulling repo")
	o := &git.PullOptions{
		// RemoteName: "origin",
		RemoteURL: s.GitRepo,
		Auth: &http.BasicAuth{
			Username: s.GitUsername, // yes, this can be anything except an empty string
			Password: os.Getenv(s.GitPasswordEnv),
		},
		// Progress: os.Stdout,
	}
	err = wt.Pull(o)
	if err != nil {
		if !errors.Is(err, git.NoErrAlreadyUpToDate) {
			return "", "", err
		}
	}

	ref, err := r.Head()
	assert.ErrNotNil(err, assert.Panic, "cannot get HEAD from git repo")

	l.Info("local git repo", zap.String("HEAD", ref.Name().Short()), zap.String("hash", ref.Hash().String()))
	return filepath.Join(s.LocalDir, s.GitRelativePath), s.BucketName, nil
}

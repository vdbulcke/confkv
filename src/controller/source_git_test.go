package controller_test

import (
	"context"
	"log/slog"
	"testing"

	"github.com/vdbulcke/confkv/src/controller"
)

func TestGitSyn(t *testing.T) {

	gs := &controller.GitSource{
		BucketName:      "test",
		LocalDir:        "../../data/git",
		GitRepo:         "https://github.com/vdbulcke/confkv.git",
		GitRelativePath: "./example/config/",
		GitUsername:     "oauth2",
		GitPasswordEnv:  "CONFKV_GIT_PASSWORD",
	}
	gs.WithLogger(nil)

	dir, bucket, err := gs.Sync(context.Background())
	if err != nil {
		t.Fatal(err)
	}

	slog.Info("success", "dir", dir, "bucket", bucket)
}

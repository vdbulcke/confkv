package cmd

import (
	"context"
	"os"

	"github.com/spf13/cobra"
	"github.com/vdbulcke/confkv/src/assert"
	"github.com/vdbulcke/confkv/src/config"
	"github.com/vdbulcke/confkv/src/controller"
	"github.com/vdbulcke/confkv/src/logger"
	"github.com/vdbulcke/confkv/src/storage"
)

// var dir string
var gitRepo string
var gitUsername string
var gitPasswordEnv string
var gitRelativePath string

func init() {
	// bind to root command
	syncCmd.AddCommand(syncGitCmd)

	syncGitCmd.Flags().StringVarP(&configFile, "config", "c", "", "config file")
	//nolint
	syncGitCmd.MarkFlagRequired("config")
	syncGitCmd.Flags().StringVarP(&dir, "dir", "", "", "local dir to sync")
	//nolint
	syncGitCmd.MarkFlagRequired("dir")
	syncGitCmd.Flags().StringVarP(&bucketName, "bucket", "", "", "bucket name")
	//nolint
	syncGitCmd.MarkFlagRequired("bucket")
	syncGitCmd.Flags().StringVarP(&gitRepo, "repo", "", "", "git https repo url")
	//nolint
	syncGitCmd.MarkFlagRequired("repo")
	syncGitCmd.Flags().StringVarP(&gitUsername, "git-username", "", "oauth2", "git username")
	syncGitCmd.Flags().StringVarP(&gitPasswordEnv, "git-password-env", "", "CONFKV_GIT_PASSWORD", "git password env vars")
	syncGitCmd.Flags().StringVarP(&gitRelativePath, "git-relative-path", "", ".", "git repo relative path")

}

var syncGitCmd = &cobra.Command{
	Use:   "git",
	Short: "Commands for manually synchronizing the db from remote git",
	// Long: "",
	Run: func(cmd *cobra.Command, args []string) {

		cfg := config.MustOpen(configFile)
		assert.NotNil(cfg.DB, assert.Exit, "missing required DB config")

		if !cfg.DB.Validate() {
			os.Exit(1)
		}

		singleSource := map[string]config.Bucket{
			bucketName: {
				Source:            config.SourceGit,
				LocalDir:          dir,
				GitRepoUrl:        gitRepo,
				GitRelativeDir:    gitRelativePath,
				GitUsername:       gitUsername,
				GitPasswordEnvVar: gitPasswordEnv,
			},
		}

		l := logger.GetLogger(Debug)

		db, err := storage.NewKVStore(cfg.DB.Dir,
			storage.WithDBName(cfg.DB.Name),
			storage.WithLogger(l),
		)
		assert.ErrNotNil(err, assert.Exit, "require DB connection")

		ctrl, err := controller.NewController(db, singleSource,
			controller.WithLogger(l),
		)
		assert.ErrNotNil(err, assert.Exit, "require valid controller")

		err = ctrl.SyncJob(context.Background())
		assert.ErrNotNil(err, assert.Exit, "sync error")

	},
}

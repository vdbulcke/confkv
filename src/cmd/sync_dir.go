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

var dir string

func init() {
	// bind to root command
	syncCmd.AddCommand(syncDirCmd)

	syncDirCmd.Flags().StringVarP(&configFile, "config", "c", "", "config file")
	//nolint
	syncDirCmd.MarkFlagRequired("config")
	syncDirCmd.Flags().StringVarP(&dir, "dir", "", "", "local dir to sync")
	//nolint
	syncDirCmd.MarkFlagRequired("dir")
	syncDirCmd.Flags().StringVarP(&bucketName, "bucket", "", "", "bucket name")
	//nolint
	syncDirCmd.MarkFlagRequired("bucket")

}

var syncDirCmd = &cobra.Command{
	Use:   "dir",
	Short: "Commands for manually synchronizing the db from local dir",
	// Long: "",
	Run: func(cmd *cobra.Command, args []string) {

		cfg := config.MustOpen(configFile)
		assert.NotNil(cfg.DB, assert.Exit, "missing required DB config")

		if !cfg.DB.Validate() {
			os.Exit(1)
		}

		singleSource := map[string]config.Bucket{
			bucketName: {
				Source:   config.SourceDir,
				LocalDir: dir,
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

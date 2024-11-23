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

func init() {
	// bind to root command
	syncCmd.AddCommand(syncAllCmd)

	syncAllCmd.Flags().StringVarP(&configFile, "config", "c", "", "config file")
	//nolint
	syncAllCmd.MarkFlagRequired("config")

}

var syncAllCmd = &cobra.Command{
	Use:   "all",
	Short: "Commands for manually synchronizing the db from config ",
	// Long: "",
	Run: func(cmd *cobra.Command, args []string) {

		cfg := config.MustOpen(configFile)
		assert.NotNil(cfg.DB, assert.Exit, "missing required DB config")

		if !cfg.Validate() {
			os.Exit(1)
		}

		l := logger.GetLogger(Debug)

		db, err := storage.NewKVStore(cfg.DB.Dir,
			storage.WithDBName(cfg.DB.Name),
			storage.WithLogger(l),
		)
		assert.ErrNotNil(err, assert.Exit, "require DB connection")

		ctrl, err := controller.NewController(db, cfg.Buckets,
			controller.WithLogger(l),
		)
		assert.ErrNotNil(err, assert.Exit, "require valid controller")

		err = ctrl.SyncJob(context.Background())
		assert.ErrNotNil(err, assert.Exit, "sync error")

	},
}

package cmd

import (
	"context"
	"os"

	"github.com/spf13/cobra"
	"github.com/vdbulcke/confkv/src/assert"
	"github.com/vdbulcke/confkv/src/config"
	"github.com/vdbulcke/confkv/src/logger"
	"github.com/vdbulcke/confkv/src/storage"
)

var key string

func init() {
	// bind to root command
	adminDeleteCmd.AddCommand(adminDeleteKeyCmd)

	adminDeleteKeyCmd.Flags().StringVarP(&configFile, "config", "c", "", "config file")
	//nolint
	adminDeleteKeyCmd.MarkFlagRequired("config")
	adminDeleteKeyCmd.Flags().StringVarP(&bucketName, "bucket", "", "", "bucket name")
	//nolint
	adminDeleteKeyCmd.MarkFlagRequired("bucket")
	adminDeleteKeyCmd.Flags().StringVarP(&key, "key", "", "", "key name")
	//nolint
	adminDeleteKeyCmd.MarkFlagRequired("key")
}

var adminDeleteKeyCmd = &cobra.Command{
	Use:   "key",
	Short: "Commands for deleting from db",
	Long: `
		See https://github.com/etcd-io/bbolt/tree/main/cmd/bbolt CLI for 
		extra functionalities
		`,
	Run: func(cmd *cobra.Command, args []string) {

		cfg := config.MustOpen(configFile)
		assert.NotNil(cfg.DB, assert.Exit, "missing required DB config")

		if !cfg.DB.Validate() {
			os.Exit(1)
		}

		l := logger.GetLogger(Debug)

		db, err := storage.NewKVStore(cfg.DB.Dir,
			storage.WithDBName(cfg.DB.Name),
			storage.WithLogger(l),
		)
		assert.ErrNotNil(err, assert.Exit, "require DB connection")

		err = db.Delete(context.Background(), bucketName, key)
		assert.ErrNotNil(err, assert.Exit, "fail to delete name")

	},
}

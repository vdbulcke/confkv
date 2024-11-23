package cmd

import (
	"github.com/spf13/cobra"
	"github.com/vdbulcke/confkv/src/assert"
)

var bucketName string

func init() {
	// bind to root command
	rootCmd.AddCommand(syncCmd)

}

var syncCmd = &cobra.Command{
	Use:   "sync",
	Short: "Commands for manually synchronizing the db",
	// Long: "",
	Run: func(cmd *cobra.Command, args []string) {

		// command does nothing
		err := cmd.Help()
		assert.ErrNotNil(err, assert.Exit)

	},
}

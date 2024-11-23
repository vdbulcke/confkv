package cmd

import (
	"github.com/spf13/cobra"
	"github.com/vdbulcke/confkv/src/assert"
)

func init() {
	// bind to root command
	adminCmd.AddCommand(adminDeleteCmd)

}

var adminDeleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Commands for deleting from db",
	Long: `
		See https://github.com/etcd-io/bbolt/tree/main/cmd/bbolt CLI for 
		extra functionalities
		`,
	Run: func(cmd *cobra.Command, args []string) {

		// command does nothing
		err := cmd.Help()
		assert.ErrNotNil(err, assert.Exit)

	},
}

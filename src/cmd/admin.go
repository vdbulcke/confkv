package cmd

import (
	"github.com/spf13/cobra"
	"github.com/vdbulcke/confkv/src/assert"
)

func init() {
	// bind to root command
	rootCmd.AddCommand(adminCmd)

}

var adminCmd = &cobra.Command{
	Use:   "admin",
	Short: "Commands for administrating the db",
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

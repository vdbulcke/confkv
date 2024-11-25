package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/spf13/cobra"
)

var Debug bool

func init() {

	// add global("persistent") flag
	rootCmd.PersistentFlags().BoolVarP(&Debug, "debug", "d", false, "debug mode enabled")

}

var rootCmd = &cobra.Command{
	Use:   "confkv",
	Short: "confkv  grpc server for config files",

	Run: func(cmd *cobra.Command, args []string) {

		// Root command does nothing
		err := cmd.Help()
		if err != nil {
			log.Fatal(err)
		}
		os.Exit(1)
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

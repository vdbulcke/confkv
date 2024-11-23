package cmd

import (
	"os"

	"github.com/spf13/cobra"
	"github.com/vdbulcke/confkv/src/client"
	"github.com/vdbulcke/confkv/src/logger"
	"go.uber.org/zap"
)

var srvAddr string

var out string

func init() {
	// bind to root command
	rootCmd.AddCommand(clientCmd)

	clientCmd.Flags().StringVarP(&srvAddr, "addr", "", "localhost:5000", "server address")
	clientCmd.Flags().StringVarP(&out, "out", "", "STDOUT", "output (STDOUT for print to terminal)")

	clientCmd.Flags().StringVarP(&bucketName, "bucket", "", "", "bucket name")
	//nolint
	clientCmd.MarkFlagRequired("bucket")
	clientCmd.Flags().StringVarP(&key, "key", "", "", "key name")
	//nolint
	clientCmd.MarkFlagRequired("key")

}

var clientCmd = &cobra.Command{
	Use:   "client",
	Short: "get value from server",
	// Long: "",
	Run: func(cmd *cobra.Command, args []string) {

		l := logger.GetLogger(true)

		client, err := client.NewClient(srvAddr)
		if err != nil {
			l.Error("error creating client", zap.Error(err))
			os.Exit(1)
		}
		defer client.Close()

		data, err := client.Get(bucketName, key)
		if err != nil {
			l.Error("error getting data", zap.Error(err))
			client.Close()
			os.Exit(1)
		}

		if out == "STDOUT" {
			os.Stdout.Write(data)
			return
		}

		f, err := os.OpenFile(out, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0666)
		if err != nil {
			l.Error("error opening output file ", zap.String("file", out), zap.Error(err))
			client.Close()
			os.Exit(1)
		}

		_, err = f.Write(data)
		if err != nil {
			l.Error("error writing output file ", zap.String("file", out), zap.Error(err))
			client.Close()
			os.Exit(1)
		}
	},
}

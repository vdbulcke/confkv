package cmd

import (
	"os"
	"sync"
	"time"

	"github.com/spf13/cobra"
	"github.com/vdbulcke/confkv/src/client"
	"github.com/vdbulcke/confkv/src/logger"
	"go.uber.org/zap"
)

func init() {
	// bind to root command
	rootCmd.AddCommand(benchCmd)

	benchCmd.Flags().StringVarP(&srvAddr, "addr", "", "localhost:5000", "server address")

	benchCmd.Flags().StringVarP(&bucketName, "bucket", "", "", "bucket name")
	//nolint
	benchCmd.MarkFlagRequired("bucket")
	benchCmd.Flags().StringVarP(&key, "key", "", "", "key name")
	//nolint
	benchCmd.MarkFlagRequired("key")

}

var benchCmd = &cobra.Command{
	Use:   "bench",
	Short: "runs 100 requests GET request to the server in //",
	// Long: "",
	Run: func(cmd *cobra.Command, args []string) {

		l := logger.GetLogger(true)

		defer func(start time.Time) {
			l.Info("run completed", zap.Duration("elasped", time.Since(start)))
		}(time.Now())

		client, err := client.NewClient(srvAddr)
		if err != nil {
			l.Error("error creating client", zap.Error(err))
			os.Exit(1)
		}
		defer client.Close()

		var wg sync.WaitGroup

		for i := 1; i <= 100; i++ {

			wg.Add(1)
			go func() {
				outcome := "success"

				defer wg.Done()
				defer func(start time.Time) {
					l.Info("job complete", zap.Int("id", i), zap.Duration("elasped", time.Since(start)), zap.String("outcome", outcome))
				}(time.Now())

				_, err := client.Get(bucketName, key)
				if err != nil {
					outcome = "failure"
				}

			}()

		}

		wg.Wait()

	},
}

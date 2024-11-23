package cmd

import (
	"os"

	"github.com/spf13/cobra"
	"github.com/vdbulcke/confkv/src/config"
	"github.com/vdbulcke/confkv/src/controller"
	"github.com/vdbulcke/confkv/src/logger"
	"github.com/vdbulcke/confkv/src/scheduler"
	"github.com/vdbulcke/confkv/src/server"
	"github.com/vdbulcke/confkv/src/storage"
	"go.uber.org/zap"
)

var port int
var monitoringPort int
var configFile string

func init() {
	// bind to root command
	rootCmd.AddCommand(runCmd)

	runCmd.Flags().StringVarP(&configFile, "config", "c", "", "config file")
	//nolint
	runCmd.MarkFlagRequired("config")

	// add global("persistent") flag
	runCmd.Flags().IntVarP(&monitoringPort, "prometheus-port", "", server.DefaultMonitoringPort, "server prometheus listen port")

	runCmd.Flags().IntVarP(&port, "port", "p", server.DefaultGRPCPort, "server listen port ")

}

var runCmd = &cobra.Command{
	Use:   "run",
	Short: "runs the server",
	// Long: "",
	Run: runServer,
}

func runServer(cmd *cobra.Command, args []string) {

	cfg := config.MustOpen(configFile)

	if !cfg.Validate() {
		os.Exit(1)
	}

	l := logger.GetLogger(Debug)

	db, err := storage.NewKVStore(cfg.DB.Dir,
		storage.WithDBName(cfg.DB.Name),
		storage.WithLogger(l),
	)
	if err != nil {
		l.Error("error creating storage db", zap.Error(err))
		os.Exit(1)
	}

	ctrl, err := controller.NewController(db, cfg.Buckets,
		controller.WithLogger(l),
	)
	if err != nil {
		l.Error("error creating controller", zap.Error(err))
		os.Exit(1)
	}

	sched, err := scheduler.NewScheduler(cfg.Cron.CronTab, ctrl, l)
	if err != nil {
		l.Error("error creating scheduler", zap.Error(err))
		os.Exit(1)
	}

	srv, err := server.NewServer(ctrl, sched,
		server.WithGrpcPort(port),
		server.WithMonitoringPort(monitoringPort),
		server.WithLogger(l),
	)
	if err != nil {
		l.Error("error creating server", zap.Error(err))
		os.Exit(1)
	}

	if err := srv.StartBlocking(); err != nil {
		l.Error("error starting server", zap.Error(err))
		os.Exit(1)
	}

}

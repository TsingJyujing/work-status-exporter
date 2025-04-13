package cmd

import (
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"os"
	"time"
	"work-status-exporter/logging"
	"work-status-exporter/workstatus"
)

var rootCmd = &cobra.Command{
	Use:   "workstatus",
	Short: "Start monitoring server",
}

var (
	httpAddr           string
	verbose            bool
	monitoringInterval time.Duration
)

func init() {
	rootCmd.AddCommand(startMonitoringServerCommand)
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "verbose output")
	startMonitoringServerCommand.PersistentFlags().StringVarP(&httpAddr, "hostAddr", "a", ":9000", "host address, default :9000")
	startMonitoringServerCommand.PersistentFlags().DurationVarP(&monitoringInterval, "monitoringInterval", "i", time.Second*5, "monitoring interval, default 5s")
}

var startMonitoringServerCommand = &cobra.Command{
	Use:   "server",
	Short: "Starting monitoring server",
	Run: func(cmd *cobra.Command, args []string) {
		if verbose {
			logging.Logger.Info("verbose output is enabled")
			logging.Logger.SetLevel(logrus.DebugLevel)
		}
		workstatus.StartPrometheusMetricsServer(httpAddr, monitoringInterval)
	},
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		logging.Logger.WithError(err).Fatal("Error executing root command")
		os.Exit(1)
	}
}

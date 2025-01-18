package main

import (
	goflag "flag"
	"fmt"
	"os"

	"github.com/golang/glog"
	"github.com/gonzojive/heatpump/cmd/heatpump/logstats"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

func main() {
	rootCmd.Execute()
}

var rootCmd = &cobra.Command{
	Use:   "heatpump",
	Short: "heatpump exposes a gRPC service to control a Chiltrix heatpump",
	Long:  "",
	Run: func(cmd *cobra.Command, args []string) {
		glog.Errorf("root command is not intended for execution")
		os.Exit(1)
	},
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of heatpump",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("v0.0.1")
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
	rootCmd.AddCommand(logstats.LogCommand)
	pflag.CommandLine.AddGoFlagSet(goflag.CommandLine)
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		glog.Errorf("exiting program due to error")
		os.Exit(1)
	}
}

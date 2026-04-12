package main

import (
	"context"
	goflag "flag"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

var rootCmd = &cobra.Command{
	Use:   "chiltrix",
	Short: "Unified CLI to manage the Chiltrix CX34 heatpump",
	RunE: func(cmd *cobra.Command, args []string) error {
		return cmd.Help()
	},
}

func main() {
	// Expose standard go log flags to pflag
	pflag.CommandLine.AddGoFlagSet(goflag.CommandLine)

	if err := rootCmd.ExecuteContext(context.Background()); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

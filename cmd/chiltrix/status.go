package main

import (
	"context"
	"fmt"
	"time"

	"github.com/gonzojive/heatpump/cx34"
	"github.com/gonzojive/heatpump/proto/chiltrix"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	statusGrpcAddr string
)

var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Query the status of the local chiltrix readwrite service",
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx, cancel := context.WithTimeout(cmd.Context(), 10*time.Second)
		defer cancel()

		conn, err := grpc.DialContext(ctx, statusGrpcAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			return fmt.Errorf("failed to connect to readwrite service at %q: %w", statusGrpcAddr, err)
		}
		defer conn.Close()

		client := chiltrix.NewReadWriteServiceClient(conn)
		resp, err := client.GetState(ctx, &chiltrix.GetStateRequest{})
		if err != nil {
			return fmt.Errorf("GetState RPC failed: %w", err)
		}

		state, err := cx34.StateFromProto(resp)
		if err != nil {
			return fmt.Errorf("error parsing state from Chiltrix proto: %w", err)
		}

		fmt.Println("--- CX34 Heat Pump Status ---")
		fmt.Println(state.Report(false, nil))
		return nil
	},
}

func init() {
	statusCmd.Flags().StringVar(&statusGrpcAddr, "grpc-address", "localhost:8084", "Address of the local readwrite service to query")
	rootCmd.AddCommand(statusCmd)
}

package grpcspec

import (
	"flag"
	"fmt"
	"strconv"
)

// Example of using ClientParamsFlag with the flag package.
func ExampleClientParamsFlag() {
	// Define a new flag set
	fs := flag.NewFlagSet("example", flag.ContinueOnError)

	// Define the grpc-client flag within the new flag set
	var cpFlag ClientParamsFlag
	fs.Var(&cpFlag, "grpc-client", "gRPC client parameters (e.g. grpc://localhost:50051?insecure=true)")

	// Parse the flags with a predefined argument list
	err := fs.Parse([]string{"-grpc-client=grpc://localhost:50051?insecure=true"})
	if err != nil {
		panic(err)
	}

	if cpFlag.Value != nil {
		fmt.Printf("Address: %s\n", cpFlag.Value.Addr())
		fmt.Printf("Insecure: %s\n", strconv.FormatBool(cpFlag.Value.Insecure()))
		fmt.Printf("URL: %s\n", cpFlag.Value.URL())
	} else {
		fmt.Printf("no value!\n")
	}
	// Output:
	// Address: localhost:50051
	// Insecure: true
	// URL: grpc://localhost:50051?insecure=true
}

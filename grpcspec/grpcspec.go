// Package grpcspec defines a means for creating a gRPC client from a flag.
package grpcspec

import (
	"fmt"
	"net/url"
	"strconv"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// ClientParams is an immutable object that contains parameters that describe
// how to create a gRPC client that connects to a remote server.
type ClientParams struct {
	addr     string
	insecure bool
}

// ParseURL parses a URL like `grpc://address:port?insecure={true|false}` into ClientParams.
func ParseURL(urlStr string) (*ClientParams, error) {
	u, err := url.Parse(urlStr)
	if err != nil {
		return nil, err
	}

	if u.Scheme != "grpc" {
		return nil, fmt.Errorf("invalid scheme: %q", u.Scheme)
	}

	cp := &ClientParams{
		addr: u.Host,
	}

	q := u.Query()
	if insecureStr := q.Get("insecure"); insecureStr != "" {
		insecure, err := strconv.ParseBool(insecureStr)
		if err != nil {
			return nil, fmt.Errorf("invalid insecure value: %q", insecureStr)
		}
		cp.insecure = insecure
	}

	return cp, nil
}

// URL returns a URL that can be parsed back into ClientParams by
// [ParseURL].
func (cp *ClientParams) URL() *url.URL {
	return &url.URL{
		Scheme:   "grpc",
		Host:     cp.addr,
		RawQuery: fmt.Sprintf("insecure=%s", strconv.FormatBool(cp.insecure)),
	}
}

// Addr returns the address of the remote server.
func (cp *ClientParams) Addr() string { return cp.addr }

// Insecure returns true if the connection to the remote server should be insecure.
func (cp *ClientParams) Insecure() bool { return cp.insecure }

// ToBuilder converts ClientParams to a ClientParamsBuilder.
func (cp *ClientParams) ToBuilder() *ClientParamsBuilder {
	return &ClientParamsBuilder{
		addr:     cp.addr,
		insecure: cp.insecure,
	}
}

// ClientParamsBuilder is a builder for ClientParams.
type ClientParamsBuilder struct {
	addr     string
	insecure bool
}

// NewClientParamsBuilder creates a new ClientParamsBuilder.
func NewClientParamsBuilder() *ClientParamsBuilder {
	return &ClientParamsBuilder{}
}

// SetAddr sets the address of the remote server.
func (b *ClientParamsBuilder) SetAddr(addr string) *ClientParamsBuilder {
	b.addr = addr
	return b
}

// SetInsecure sets whether the connection to the remote server should be insecure.
func (b *ClientParamsBuilder) SetInsecure(insecure bool) *ClientParamsBuilder {
	b.insecure = insecure
	return b
}

// Build builds a ClientParams from the builder.
func (b *ClientParamsBuilder) Build() *ClientParams {
	return &ClientParams{
		addr:     b.addr,
		insecure: b.insecure,
	}
}

// ClientParamsFlag is a custom flag type for parsing ClientParams from a URL string.
// It implements the flag.Value interface.
type ClientParamsFlag struct {
	Value *ClientParams
}

// String returns the string representation of the ClientParams, which is the URL string.
func (v *ClientParamsFlag) String() string {
	if v.Value != nil {
		return v.Value.URL().String()
	}
	return ""
}

// Set parses the URL string and sets the ClientParams.
func (v *ClientParamsFlag) Set(s string) error {
	if cp, err := ParseURL(s); err != nil {
		return fmt.Errorf("failed to parse flag: %w", err)
	} else {
		v.Value = cp
	}
	return nil
}

// NewClient creates a new gRPC client for the given service.
// It takes a ClientParams and a function that creates a new client for the specific service.
func NewClient[T any](params *ClientParams, newClientFunc func(conn grpc.ClientConnInterface) T) (T, *grpc.ClientConn, error) {
	var dialOpts []grpc.DialOption
	if params.Insecure() {
		dialOpts = append(dialOpts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	}
	conn, err := grpc.NewClient(params.Addr(), dialOpts...)
	if err != nil {
		var zero T
		return zero, conn, fmt.Errorf("error creating connection: %w", err)
	}
	return newClientFunc(conn), conn, nil
}

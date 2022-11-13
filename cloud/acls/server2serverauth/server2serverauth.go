// Package server2serverauth is used for sending gRPC requests between Google
// Cloud  hosted services (like Cloud Run instances).
package server2serverauth

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/golang/glog"
	"github.com/gonzojive/heatpump/util/must"
	"golang.org/x/oauth2"
	"google.golang.org/api/idtoken"

	"google.golang.org/grpc/credentials"
	grpcmetadata "google.golang.org/grpc/metadata"
)

const authorizationHeader = "Authorization"

type Sender struct {
	audienceURL string
	tokenSource oauth2.TokenSource
}

// NewSender returns a new object used for sending authenticating information
// to another server.
func NewSender(ctx context.Context, audience string) (*Sender, error) {
	// client is a http.Client that automatically adds an "Authorization" header
	// to any requests made.
	tokSource, err := idtoken.NewTokenSource(ctx, audience)
	if err != nil {
		return nil, fmt.Errorf("idtoken.NewClient error: %w", err)
	}

	return &Sender{audience, tokSource}, nil
}

// AddTokenToOutgoingContext returns a new context.Context with an
// ("Authorozation", "Bearer xxxx") header in the gRPC metadata.
func (p *Sender) AddTokenToOutgoingContext(ctx context.Context) (context.Context, error) {
	tok, err := p.tokenSource.Token()
	if err != nil {
		return nil, fmt.Errorf("Token() error: %w", err)
	}
	return grpcmetadata.AppendToOutgoingContext(ctx, authorizationHeader, tok.Type()+" "+tok.AccessToken), nil
}

// PerRPCCredentials returns a new context.Context with an
// ("Authorozation", "Bearer xxxx") header in the gRPC metadata.
func (p *Sender) PerRPCCredentials() credentials.PerRPCCredentials {
	return &perRPCCredentials{
		func(uri string) bool { return true },
		p,
	}
}

type perRPCCredentials struct {
	uriMatcher func(uri string) bool
	s          *Sender
}

// GetRequestMetadata gets the current request metadata, refreshing
// tokens if required. This should be called by the transport layer on
// each request, and the data should be populated in headers or other
// context. If a status code is returned, it will be used as the status
// for the RPC. uri is the URI of the entry point for the request.
// When supported by the underlying implementation, ctx can be used for
// timeout and cancellation. Additionally, RequestInfo data will be
// available via ctx to this call.
// TODO(zhaoq): Define the set of the qualified keys instead of leaving
// it as an arbitrary string.
func (c *perRPCCredentials) GetRequestMetadata(ctx context.Context, uri ...string) (map[string]string, error) {
	glog.Infof("GetRequestMetadata(%s called)", string(must.Value(json.Marshal(uri))))
	if len(uri) != 1 {
		return nil, fmt.Errorf("perRPCCredentials can't deal with len(uri) = %d: %v", len(uri), uri)
	}
	if !c.uriMatcher(uri[0]) {
		return nil, fmt.Errorf("perRPCCredentials not supposed to be use for uri = %q", uri[0])
	}
	tok, err := c.s.tokenSource.Token()
	if err != nil {
		return nil, fmt.Errorf("perRPCCredentials failed to get auth token: %w", err)
	}

	return map[string]string{
		authorizationHeader: tok.Type() + " " + tok.AccessToken,
	}, nil
}

// RequireTransportSecurity indicates whether the credentials requires
// transport security.
func (c *perRPCCredentials) RequireTransportSecurity() bool {
	return true
}

// Validator is used for validating requests coming from other servers.
type Validator struct {
	audience string
	v        *idtoken.Validator
}

// NewValidator constructs a Validator for use by a server that needs to
// validate other servers.
func NewValidator(ctx context.Context, audience string) (*Validator, error) {
	v, err := idtoken.NewValidator(ctx)
	if err != nil {
		return nil, fmt.Errorf("error creating validator for server identified by audience %q: %w", audience, err)
	}
	return &Validator{audience, v}, nil
}

// ValidateFromIncomingContext extracts a validated payload from the gRPC
// metadata of an incoming request.
func (v *Validator) ValidateFromIncomingContext(ctx context.Context) (*IncomingTokenPayload, error) {
	md, ok := grpcmetadata.FromIncomingContext(ctx)
	if !ok {
		return nil, fmt.Errorf("error getting gRPC metadata from context while validating token")
	}
	values := md.Get(authorizationHeader)
	if got, want := len(values), 1; got != want {
		return nil, fmt.Errorf("got %d values for gRPC metadata header %q, want 1", got, authorizationHeader)
	}

	token := strings.TrimPrefix(values[0], "Bearer ")
	if token == values[0] {
		return nil, fmt.Errorf("unexpected %q header value doesn't have 'Bearer ' pefix", authorizationHeader)
	}

	payload, err := idtoken.Validate(ctx, token, v.audience)
	if err != nil {
		return nil, fmt.Errorf("validation error: %w", err)
	}

	email, ok := payload.Claims["email"].(string)
	if !ok {
		return nil, fmt.Errorf("payload is missing 'email' claim")
	}

	if verified, ok := payload.Claims["email_verified"].(bool); !ok {
		return nil, fmt.Errorf("payload is missing email_verified header for %q", email)
	} else if !verified {
		return nil, fmt.Errorf("email %q must be verified, but it isn't", email)
	}

	return &IncomingTokenPayload{payload, email}, nil
}

// IncomingTokenPayload is data embedded in the token that has been validated by
// a trusted authority (typically Google).
type IncomingTokenPayload struct {
	p     *idtoken.Payload
	email string
}

func (p *IncomingTokenPayload) Email() string {
	return p.email
}

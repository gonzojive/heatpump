// Package acls is the main library used by cloud services to authenticate
// clients and check that they have access to named resources.
package acls

import (
	"context"
	"crypto/x509"
	"encoding/base64"
	"fmt"
	"time"

	"github.com/gonzojive/heatpump/cloud/google/cloudconfig"
	"github.com/gonzojive/heatpump/util/must"
	"github.com/samber/lo"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/peer"
	"google.golang.org/grpc/status"

	grpcmetadata "google.golang.org/grpc/metadata"

	_ "embed"
)

// DeviceAccessTokenMetadataKey is the gRPC metadata key used to provide a
// device access token identity secret to the server.
const DeviceAccessTokenMetadataKey = "heatpump-device-access-token"

var (
	//go:embed client-signer-cert-authority-cert.pem
	clientSignerCertAuthorityCertData []byte

	// ClientCACertPool is the certificate pool that should be used by gRPC
	// services as the tls.Config.ClientCAs attribute for authenticating
	// clients.
	//
	// The certificate embedded in this package was generated using the
	// generate_client_cert_signer.sh script.
	ClientCACertPool = must.Compute(func() (*x509.CertPool, error) {
		pool := x509.NewCertPool()
		if !pool.AppendCertsFromPEM(clientSignerCertAuthorityCertData) {
			return nil, fmt.Errorf("invalid client certificate")
		}
		return pool, nil
	})

	hardcodedIdentities = []*Identity{
		{
			tlsClientCertPublicKeys: []string{
				"MIICIjANBgkqhkiG9w0BAQEFAAOCAg8AMIICCgKCAgEA2QdskQ+S6UHxzoB912k8dr615pBkScmBPQtOrE6Fyr/842Wl67YWrwSnoHx225173m0U79yKqK5bAlRh+nfbKBKK6aQP7BJzoocLQhQJl0hYapqJbtl1WgcJ0DafN4AEcc/jT/ETc9lZsiFhhTHnVxaIJ5CNVPr4C18+1dXEneusfpfBOVAzU+wnQFNU7DJJxR7+YoWj1+t6XpvHJ+L6RnAQJzeC0otujx4VCttLVgbctJXwQZp/t/0WaEixsel01n7EdBUfZubQfVIvdKVI+PblYEaXmEnBlE4KndULr071e2rluLXSY22VjKxxBa73eVeEIE0mx58ABV/PWU9kkViON6Oj3U5vqefNk76rZCcKMMbZS0zEZJo1QaUhYbh5x5XCH5ZVJ4YD5qeIDWhTw9+qLZAJ1FMUGmODMjReGQj0eCAzBUpw1swqqo3EqbuSTtgsl6+okM/TuBi4vwzUHYbN8XJEZcwGqMFZZo4e2I2X/Mr8y0rr8ofAseo5thi0bsuv8EMxmqac2LLQCq13W1GtMjmOGthVfm4L74QPp897bO8agqxZJgtM4H72IDBT2ZPXIAUQK0vC/+PRMarBiSBSTpgLZSSFOBbw0wamwu+pA+h3thMwxmhW1Syppe624ASTs7nqTEnwt3q1r/h58J5Qem0MT2hpBWPVGamM0KECAwEAAQ",
			},
			id: "redshouse",
		},
	}
)

func FixmeMainHardcodedIdentity() *Identity { return hardcodedIdentities[0] }

// Identity is information about a user or robot account.
type Identity struct {
	// Base64-encoded public keys of the user. These are used to extract an
	// identity from an authenticated RPC request.
	tlsClientCertPublicKeys []string

	// system-assigned identity
	id string
}

// ID returns a id of the user.
func (i *Identity) ID() string {
	return i.id
}

// String returns a debug representation of the identity.
func (i *Identity) String() string {
	return fmt.Sprintf("[%s]", i.id)
}

// NewService loads a new identity service.
func NewService(params *cloudconfig.Params) *Service {
	pubkeyToIdentity := map[string]*Identity{}

	for _, ident := range hardcodedIdentities {
		for _, pubKey := range ident.tlsClientCertPublicKeys {
			pubkeyToIdentity[pubKey] = ident
		}
	}
	return &Service{pubkeyToIdentity, NewDeviceTokenVerifier(params)}
}

// Service is the service responsible for getting identity information about a
// client like an IoT device.
type Service struct {
	pubkeyToIdentity map[string]*Identity
	tokenVerifier    *DeviceTokenVerifier
}

// IdentityFromContext returns an Identity instance
func (s *Service) IdentityFromContext(ctx context.Context) (*Identity, error) {
	peer, ok := peer.FromContext(ctx)
	if !ok {
		return nil, fmt.Errorf("could not extract peer info from context")
	}

	// If the request uses TLS, we can use the client certificate to authenticate.
	tlsInfo, ok := peer.AuthInfo.(credentials.TLSInfo)
	if !ok {
		return s.identityFromGRPCMetadata(ctx)
	}
	var keys []string
	for _, v := range tlsInfo.State.PeerCertificates {
		pubKey, err := x509.MarshalPKIXPublicKey(v.PublicKey)
		if err != nil {
			return nil, fmt.Errorf("error marshaling public key: %w", err)
		}
		pubKeyStr := base64.RawStdEncoding.EncodeToString(pubKey)
		if ident := s.pubkeyToIdentity[pubKeyStr]; ident != nil {
			return ident, nil
		}
		keys = append(keys, pubKeyStr)
	}
	return nil, fmt.Errorf("failed to identify client with public keys %s",
		lo.Map(keys, func(item string, _ int) string { return fmt.Sprintf("%q", item) }))
}

func (s *Service) identityFromGRPCMetadata(ctx context.Context) (*Identity, error) {
	md, ok := grpcmetadata.FromIncomingContext(ctx)
	if !ok {
		status.Errorf(codes.Unauthenticated, "incoming request is missing metadata")
	}

	metadataEntries := md.Get(DeviceAccessTokenMetadataKey)
	if len(metadataEntries) != 1 {
		return nil, fmt.Errorf("expected 1 gRPC metadata entry for %q, got %d", DeviceAccessTokenMetadataKey, len(metadataEntries))
	}

	time.Sleep(time.Second)
	dat, err := s.tokenVerifier.Verify(TokenString(metadataEntries[0]))
	if err != nil {
		return nil, fmt.Errorf("could not verify device access token: %w", err)
	}

	return &Identity{id: dat.UserID()}, nil
}

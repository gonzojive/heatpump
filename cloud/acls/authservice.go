package acls

import (
	"context"
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/base32"
	"encoding/base64"
	"fmt"
	"time"

	secretmanager "cloud.google.com/go/secretmanager/apiv1"
	"cloud.google.com/go/secretmanager/apiv1/secretmanagerpb"
	"github.com/golang/glog"
	"github.com/golang/protobuf/proto"
	"github.com/gonzojive/heatpump/cloud/google/cloudconfig"
	"golang.org/x/crypto/ssh"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"

	pb "github.com/gonzojive/heatpump/proto/controller"
)

const debugOnlyExtendUnsignedTokens = false // DO NOT SUBMIT

var (
	sigOptions = &rsa.PSSOptions{
		SaltLength: rsa.PSSSaltLengthAuto,
	}

	signatureStringEncoding = base64.StdEncoding
	tokenStringEncoding     = base32.StdEncoding
)

type AuthService struct {
	pb.UnimplementedAuthServiceServer

	projectParams *cloudconfig.Params

	privKey *rsa.PrivateKey
}

func NewAuthService(ctx context.Context, params *cloudconfig.Params) (*AuthService, error) {
	privKey, err := fetchPrivateKey(ctx, params.DeviceAccessTokenSecretVersionName)
	if err != nil {
		return nil, fmt.Errorf("failed to get private key used for signing device access tokens: %w", err)
	}
	publicRsaKey, err := ssh.NewPublicKey(&privKey.PublicKey)
	if err != nil {
		return nil, err
	}

	glog.Infof("Starting auth service using public key for verification:\n%s", string(ssh.MarshalAuthorizedKey(publicRsaKey)))

	return &AuthService{
		projectParams: params,
		privKey:       privKey,
	}, nil
}

func (s *AuthService) ExtendToken(ctx context.Context, req *pb.ExtendTokenRequest) (*pb.ExtendTokenResponse, error) {
	wireFormatToken, err := tokenStringEncoding.DecodeString(req.GetToken())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "bad token - not base32")
	}

	unwrappedToken := &pb.DeviceAccessToken{}
	if err := proto.Unmarshal(wireFormatToken, unwrappedToken); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "bad token - improper format")
	}
	if unwrappedToken.GetUserId() == "" {
		if debugOnlyExtendUnsignedTokens {
			unwrappedToken = &pb.DeviceAccessToken{
				UserId:     hardcodedIdentities[0].id,
				Expiration: timestamppb.New(time.Now().Add(time.Hour)),
			}
			unwrappedToken.Signature, err = s.sign(unwrappedToken)
			if err != nil {
				return nil, fmt.Errorf("failed to sign dev token: %w", err)
			}
		} else {
			return nil, status.Errorf(codes.InvalidArgument, "invalid signature contents")
		}
	}

	if err := (&DeviceTokenVerifier{&s.privKey.PublicKey}).verifyUnwrapped(unwrappedToken); err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "Token verification failed")
	}
	if unwrappedToken.GetExpiration().AsTime().Before(time.Now()) {
		return nil, status.Errorf(codes.InvalidArgument, "Token expired")
	}

	unwrappedToken.Expiration = timestamppb.New(time.Now().Add(time.Hour * 24 * 700))
	sig, err := s.sign(unwrappedToken)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to sign extended token")
	}
	unwrappedToken.Signature = sig

	wireFormatted, err := proto.Marshal(unwrappedToken)
	if err != nil {
		return nil, err
	}

	return &pb.ExtendTokenResponse{
		RefreshedToken: tokenStringEncoding.EncodeToString(wireFormatted),
		Expiration:     unwrappedToken.GetExpiration(),
	}, nil
}

func fetchPrivateKey(ctx context.Context, secretVersionName string) (*rsa.PrivateKey, error) {
	c, err := secretmanager.NewClient(ctx)
	if err != nil {
		return nil, err
	}
	defer c.Close()

	secret, err := c.AccessSecretVersion(ctx, &secretmanagerpb.AccessSecretVersionRequest{
		Name: secretVersionName,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to fetch private key from secret manager: %w", err)
	}

	signer, err := ssh.ParseRawPrivateKey(secret.GetPayload().GetData())
	if err != nil {
		return nil, fmt.Errorf("error parsing private key from secret store: %w", err)
	}
	rsaKey, ok := signer.(*rsa.PrivateKey)
	if !ok {
		return nil, fmt.Errorf("private key is not an RSA key")
	}
	return rsaKey, nil
}

func (s *AuthService) sign(unwrappedToken *pb.DeviceAccessToken) (string, error) {
	digest, err := digestForSignature(unwrappedToken)
	if err != nil {
		return "", err
	}

	sig, err := rsa.SignPSS(rand.Reader, s.privKey, crypto.SHA256, digest, sigOptions)
	if err != nil {
		return "", err
	}

	return signatureStringEncoding.EncodeToString(sig), nil
}

func digestForSignature(unwrappedToken *pb.DeviceAccessToken) (digest []byte, err error) {
	cloned := proto.Clone(unwrappedToken).(*pb.DeviceAccessToken)
	cloned.Signature = ""
	signedData, err := proto.Marshal(cloned)
	if err != nil {
		return nil, err
	}
	hash := sha256.New()
	hash.Write(signedData)
	return hash.Sum(nil), nil
}

// DeviceTokenVerifier verifies device access tokens sent by an IoT device.
type DeviceTokenVerifier struct {
	publicKey *rsa.PublicKey
}

func NewDeviceTokenVerifier(params *cloudconfig.Params) *DeviceTokenVerifier {
	if params.DeviceTokenSigningKey == nil {
		panic("nil DeviceTokenSigningKey")
	}
	return &DeviceTokenVerifier{publicKey: params.DeviceTokenSigningKey}
}

func (dtv *DeviceTokenVerifier) Verify(token TokenString) (*DeviceAccessToken, error) {
	parsed, err := token.parse()
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "bad token: %v", err)
	}
	if err := dtv.verifyUnwrapped(parsed.proto); err != nil {
		return nil, err
	}
	return parsed, err
}

func (dtv *DeviceTokenVerifier) verifyUnwrapped(unwrappedToken *pb.DeviceAccessToken) error {
	sig, err := signatureStringEncoding.DecodeString(unwrappedToken.GetSignature())
	if err != nil {
		return fmt.Errorf("could not base64 decode signature")
	}
	digest, err := digestForSignature(unwrappedToken)
	if err != nil {
		return err
	}

	err = rsa.VerifyPSS(dtv.publicKey, crypto.SHA256, digest, sig, sigOptions)
	if err != nil {
		return fmt.Errorf("signature verification failed")
	}

	return nil
}

// TokenString is a base32-encoded DeviceAccessToken proto. This is delievered
// to the client as an opaque string and so should not contain private data. It is sent in the metadata fields of gRPC requests to authenticate a client.
type TokenString string

// parse returns the proto encoded by this token.
func (ts TokenString) parse() (*DeviceAccessToken, error) {
	wireFormatToken, err := tokenStringEncoding.DecodeString(string(ts))
	if err != nil {
		return nil, fmt.Errorf("bad token - not base32")
	}

	unwrappedToken := &pb.DeviceAccessToken{}
	if err := proto.Unmarshal(wireFormatToken, unwrappedToken); err != nil {
		return nil, fmt.Errorf("bad token - improper format")
	}
	return &DeviceAccessToken{unwrappedToken}, nil
}

// DeviceAccessToken wraps the proto by the same name.
type DeviceAccessToken struct {
	proto *pb.DeviceAccessToken
}

func (dat *DeviceAccessToken) Proto() *pb.DeviceAccessToken {
	return proto.Clone(dat.proto).(*pb.DeviceAccessToken)
}

func (dat *DeviceAccessToken) Expiration() time.Time {
	return dat.proto.GetExpiration().AsTime()
}

func (dat *DeviceAccessToken) UserID() string {
	return dat.proto.GetUserId()
}

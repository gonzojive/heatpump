package main

import (
	"bytes"
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"flag"
	"fmt"
	"strings"

	secretmanager "cloud.google.com/go/secretmanager/apiv1"
	"github.com/golang/glog"
	"github.com/gonzojive/heatpump/cloud/google/cloudconfig"
	"golang.org/x/crypto/ssh"
	"google.golang.org/api/iterator"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	secretmanagerpb "cloud.google.com/go/secretmanager/apiv1/secretmanagerpb"
)

var (
	secretName = flag.String("secret-name", "", "Name of secret to update.")
	overwrite  = flag.Bool("overwrite", false, "True if the secret should be overwritten.")
)

func main() {
	flag.Parse()
	if err := run(context.Background()); err != nil {
		glog.Exitf("%v", err)
	}
}

func run(ctx context.Context) error {
	key, err := newKeyPair()
	if err != nil {
		return fmt.Errorf("error generating key pair: %w", err)
	}
	glog.Infof("about to write private key to secret store for public key %s", string(key.AuthorizedKeysLine))

	// This snippet has been automatically generated and should be regarded as a code template only.
	// It will require modifications to work:
	// - It may require correct/in-range values for request initialization.
	// - It may require specifying regional endpoints when creating the service client as shown in:
	//   https://pkg.go.dev/cloud.google.com/go#hdr-Client_Options
	c, err := secretmanager.NewClient(ctx)
	if err != nil {
		return err
	}
	defer c.Close()

	parent := fmt.Sprintf("projects/%s", cloudconfig.DefaultParams().GCPProject)

	fullSecretName := fmt.Sprintf("%s/secrets/%s", parent, *secretName)
	secret, err := c.CreateSecret(ctx, &secretmanagerpb.CreateSecretRequest{
		Parent:   parent,
		SecretId: *secretName,
		Secret: &secretmanagerpb.Secret{
			Replication: &secretmanagerpb.Replication{
				Replication: &secretmanagerpb.Replication_Automatic_{
					Automatic: &secretmanagerpb.Replication_Automatic{},
				},
			},
		},
	})
	if err != nil {
		if status, _ := status.FromError(err); status.Code() == codes.AlreadyExists && *overwrite {
			glog.Infof("there is already a secret with name %q, overwriting...", fullSecretName)
		} else {
			return fmt.Errorf("error creating secret: %w", err)
		}
	} else {
		if fullSecretName != secret.GetName() {
			return fmt.Errorf("invalid assumption about secret naming convention: %q (generated) != %q (actual)", fullSecretName, secret.Name)
		}
	}

	c.AddSecretVersion(ctx, &secretmanagerpb.AddSecretVersionRequest{
		Parent: fmt.Sprintf("%s/secrets/%s", parent, *secretName),
		Payload: &secretmanagerpb.SecretPayload{
			Data: []byte(key.PrivateKeyPEMFormat),
		},
	})

	it := c.ListSecrets(ctx, &secretmanagerpb.ListSecretsRequest{
		Parent: fmt.Sprintf("projects/%s", cloudconfig.DefaultParams().GCPProject),
	})
	for {
		resp, err := it.Next()
		if errors.Is(err, iterator.Done) {
			break
		}
		if err != nil {
			return fmt.Errorf("error while listing secrets: %v", err)
		}
		glog.Infof("got secret %q: %+v", resp.GetName(), resp)
	}

	return nil
}

type KeyPair struct {
	AuthorizedKeysLine  []byte
	PrivateKey          *rsa.PrivateKey
	PrivateKeyPEMFormat string
}

func newKeyPair() (*KeyPair, error) {
	rsaKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return nil, err
	}
	pemEncoded := &strings.Builder{}
	pem.Encode(pemEncoded, &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(rsaKey),
	})
	rsaPubKey, err := ssh.NewPublicKey(&rsaKey.PublicKey)
	if err != nil {
		return nil, err
	}
	return &KeyPair{
		AuthorizedKeysLine:  bytes.TrimSpace(ssh.MarshalAuthorizedKey(rsaPubKey)),
		PrivateKey:          rsaKey,
		PrivateKeyPEMFormat: pemEncoded.String(),
	}, nil
}

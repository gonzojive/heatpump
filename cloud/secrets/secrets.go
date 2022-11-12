// Package secrets is a utility library for retrieving secrets from Google
// Secret Manager.
package secrets

import (
	"fmt"

	secretmanager "cloud.google.com/go/secretmanager/apiv1"
	"cloud.google.com/go/secretmanager/apiv1/secretmanagerpb"
	"golang.org/x/net/context"
)

// Name of a secret in the secret store.
//
// This is the name of a "secret version" resource like
// "projects/my-project/secrets/my-secret-name/versions/latest"
type Name string

// Access returns the value of a secret given a secret name.
func (n Name) Access(ctx context.Context) ([]byte, error) {
	c, err := secretmanager.NewClient(ctx)
	if err != nil {
		return nil, err
	}
	defer c.Close()

	secret, err := c.AccessSecretVersion(ctx, &secretmanagerpb.AccessSecretVersionRequest{
		Name: string(n),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to fetch private key from secret manager: %w", err)
	}

	return secret.GetPayload().GetData(), nil
}

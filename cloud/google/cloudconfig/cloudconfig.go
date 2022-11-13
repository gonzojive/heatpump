package cloudconfig

import (
	"crypto/rsa"
	"fmt"
	"reflect"

	"github.com/gonzojive/heatpump/util/must"
	"golang.org/x/crypto/ssh"
)

const CommandsTopic = "iot-commands"

type Params struct {
	GCPProject                         string
	DeviceAccessTokenSecretVersionName string
	CommandsSubscriptionName           string
	TokenStoreFirebaseCollectionName   string

	// Used to authenticate device accesss token.
	DeviceTokenSigningKey *rsa.PublicKey
}

func DefaultParams() *Params {
	return &Params{
		GCPProject:                         "heatpump-dev",
		CommandsSubscriptionName:           "queueserver_pull_1",
		DeviceAccessTokenSecretVersionName: "projects/heatpump-dev/secrets/device-token-signer-private-rsa/versions/latest",
		DeviceTokenSigningKey:              must.Value(parseAuthorizedKey("ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQDVxkIyR83R1KwFHqrY+ymjyql5yEX1mTgqMrjfedFQEyybhvcQmVW9S54kgLY7htxa3JCj/0jhxPVB6uzBXvw6NkAjAfjdEIJ8koHl2X1IsMRVWpIGe4jLcF0okktAe9rztKQ/I/hXvbKLF6gEq0GiJkVMhqAxRNX/gCniTSAvIv29aYiQAI9GNWbN7rgm754v0q/UWp0U4akgUtWl/I2zrDX1tztLWz/MZWnHKaO218qTvmVfk4N1Z0iTZ8KQs7mEA06u8CzbEd3P7+j+pjEOS8iBaRHLt76bm7PFp3eQOpRSm7C36BOHmYtIEWr24shU2IQNkmD8jDu4ibFHFFzH")),
		TokenStoreFirebaseCollectionName:   "oauth2tokens-actions",
	}
}

// parseAuthorizedKey returns an *rsa.PublicKey from a line in an
// authorized_keys file.
func parseAuthorizedKey(authorizedKeyLine string) (*rsa.PublicKey, error) {
	sshKey, _, _, _, err := ssh.ParseAuthorizedKey([]byte(authorizedKeyLine))
	if err != nil {
		return nil, err
	}

	// To get back to an *rsa.PublicKey, we need to first upgrade to the
	// ssh.CryptoPublicKey interface
	parsedCryptoKey, ok := sshKey.(ssh.CryptoPublicKey)
	if !ok {
		return nil, fmt.Errorf("public key couldn't be converted to ssh.CryptoPublicKey: %+v", reflect.ValueOf(sshKey))
	}

	// Then, we can call CryptoPublicKey() to get the actual crypto.PublicKey
	rsaKey, ok := parsedCryptoKey.CryptoPublicKey().(*rsa.PublicKey)
	if !ok {
		return nil, fmt.Errorf("public key couldn't be converted to RSA: %+v", reflect.ValueOf(sshKey))
	}
	return rsaKey, nil
}

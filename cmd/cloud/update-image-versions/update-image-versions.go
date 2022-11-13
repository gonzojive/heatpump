// Program update-image-versions generates a JSON file that resolves docker
// image names to digests. This is used as part of the continuous deployment
// process.
package main

import (
	"context"
	"flag"
	"fmt"
	"net/url"

	"github.com/golang/glog"
	"github.com/gonzojive/heatpump/util/cmdutil"
	"github.com/google/go-containerregistry/pkg/authn"
	"github.com/google/go-containerregistry/pkg/name"
	"github.com/google/go-containerregistry/pkg/v1/remote"

	"golang.org/x/oauth2/google"
)

var (
	input  = flag.String("input", "", "Input JSON mapping.")
	output = flag.String("output", "", "Output JSON mapping.")
)

func main() {
	flag.Parse()
	if err := run(context.Background()); err != nil {
		glog.Exitf("%v", err)
	}
}

func run(ctx context.Context) error {
	if *input == "" {
		return fmt.Errorf("empty --input flag")
	}
	if *output == "" {
		return fmt.Errorf("empty --output flag")
	}

	inputMapping, err := cmdutil.ReadJSONFile(*input, &Mapping{})
	if err != nil {
		return err
	}

	outputMapping, err := processMapping(ctx, inputMapping)
	if err != nil {
		return err
	}

	return cmdutil.WriteJSONFile(*output, outputMapping)
}

type Mapping struct {
	RegistryURL string               `json:"registry-url"`
	Images      map[string]ImageSpec `json:"images"`
}

type ImageSpec struct {
	Spec     string `json:"spec,omitempty"`
	Resolved string `json:"resolved"`
}

func processMapping(ctx context.Context, m *Mapping) (*Mapping, error) {
	creds, err := google.FindDefaultCredentials(ctx, "https://www.googleapis.com/auth/cloud-platform")
	if err != nil {
		return nil, fmt.Errorf("failed to find Google default credentials: %w", err)
	}
	tok, err := creds.TokenSource.Token()
	if err != nil {
		return nil, fmt.Errorf("failed to generate oauth token: %w", err)
	}

	authenticator := authn.FromConfig(authn.AuthConfig{
		Username: "oauth2accesstoken",
		Password: tok.AccessToken,
	})

	url, err := url.Parse(m.RegistryURL)
	if err != nil {
		return nil, fmt.Errorf("invalid registry URL %q: %w", m.RegistryURL, err)
	}

	out := *m
	for k, entry := range out.Images {
		// NOTE: This is the best documentation of the terminology related to images that I have found.
		// https://github.com/google/go-containerregistry/blob/main/pkg/v1/remote/README.md
		ref2, err := name.ParseReference(entry.Spec)
		if err != nil {
			return nil, fmt.Errorf("failed to parse %q docker image spec %q: %w", k, entry.Spec, err)
		}

		if want, got := url.Host, ref2.Context().RegistryStr(); want != got {
			return nil, fmt.Errorf("image spec %q has a 'spec' attribute with host %q, but the 'registry-url' hostname is %q; update the input JSON file to make the two hostnames match", k, got, want)
		}

		img, err := remote.Image(ref2, remote.WithAuth(authenticator))
		if err != nil {
			return nil, fmt.Errorf("failed to get %q docker image: %w", k, err)
		}
		manifestID, err := img.Digest()
		if err != nil {
			return nil, fmt.Errorf("failed to get %q manifest digest (i.e. image id): %w", k, err)
		}

		entry.Resolved = ref2.Context().Digest(manifestID.String()).String()
		glog.Infof("resolved %q to %q", k, entry.Resolved)

		out.Images[k] = entry
	}
	return &out, nil
}

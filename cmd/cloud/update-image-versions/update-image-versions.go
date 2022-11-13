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
	"github.com/heroku/docker-registry-client/registry"

	dockerparser "github.com/novln/docker-parser"
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

	url, err := url.Parse(m.RegistryURL)
	if err != nil {
		return nil, fmt.Errorf("invalid registry URL %q: %w", m.RegistryURL, err)
	}

	username := "oauth2accesstoken" // anonymous
	password := tok.AccessToken
	hub, err := registry.New(url.String(), username, password)
	if err != nil {
		return nil, fmt.Errorf("failed to create image registry client: %w", err)
	}

	out := *m
	for k, entry := range out.Images {
		// NOTE: This is the best documentation of the terminology related to images that I have found.
		// https://github.com/google/go-containerregistry/blob/main/pkg/v1/remote/README.md
		ref2, err := name.ParseReference(entry.Spec)
		if err != nil {
			return nil, fmt.Errorf("failed to parse %q docker image spec %q: %w", k, entry.Spec, err)
		}
		img, err := remote.Image(ref2, remote.WithAuth(authn.FromConfig(authn.AuthConfig{
			Username: "oauth2accesstoken",
			Password: tok.AccessToken,
		})))
		if err != nil {
			return nil, fmt.Errorf("failed to get %q docker image: %w", k, err)
		}
		manifestID, err := img.Digest()
		if err != nil {
			return nil, fmt.Errorf("failed to get %q manifest digest (i.e. image id): %w", k, err)
		}
		imageID, err := img.ConfigName()
		if err != nil {
			return nil, fmt.Errorf("failed to get %q config name (i.e. image id): %w", k, err)
		}
		ref, err := dockerparser.Parse(entry.Spec)
		if err != nil {
			return nil, fmt.Errorf("failed to parse %q docker image spec %q: %w", k, entry.Spec, err)
		}
		if want, got := url.Host, ref.Registry(); want != got {
			return nil, fmt.Errorf("image spec %q has a 'spec' attribute with host %q, but the 'registry-url' hostname is %q; update the input JSON file to make the two hostnames match", k, got, want)
		}

		//glog.Infof("%q: parsed remote as reg = %q, repo = %q, name = %q, shortname = %q, tag = %q", k, ref.Registry(), ref.Repository(), ref.Name(), ref.ShortName(), ref.Tag())

		digest, err := hub.ManifestDigest(ref.ShortName(), ref.Tag())
		if err != nil {
			return nil, fmt.Errorf("failed to get digest for entry %q: %w", k, err)
		}
		//glog.Infof("%q got manifest %+v w/digest %q", k, mf, digest)

		glog.Infof("image info for %q:\n  config name = %q\n  manifest id = %q\n  other lib   = %q", k, imageID.String(), manifestID.String(), digest)

		finalDigest := manifestID.String()

		//digest, err := hub.ManifestDigest(ref.ShortName(), ref.Tag())
		entry.Resolved = fmt.Sprintf("%s:%s", ref.Repository(), finalDigest)
		//entry.Resolved = fmt.Sprintf("%s:main", ref.Repository())
		out.Images[k] = entry
	}
	return &out, nil
}

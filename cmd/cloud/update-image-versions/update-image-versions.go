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
	Repository string `json:"repository"`
	Reference  string `json:"reference"`
	Remote     string `json:"remote,omitempty"`
	Resolved   string `json:"resolved"`
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
		digest, err := hub.ManifestDigest(entry.Repository, entry.Reference)
		if err != nil {
			return nil, fmt.Errorf("failed to get digest for entry %q: %w", k, err)
		}
		entry.Resolved = fmt.Sprintf("%s/%s@%s", url.Hostname(), entry.Repository, digest.String())
		out.Images[k] = entry

		if entry.Remote != "" {
			ref, err := dockerparser.Parse(entry.Remote)
			if err != nil {
				return nil, fmt.Errorf("failed to parse %q docker image spec %q: %w", k, entry.Remote, err)
			}
			glog.Infof("parsed remote as reg = %q, repo = %q, name = %q, shortname = %q, tag = %q", ref.Registry(), ref.Repository(), ref.Name(), ref.ShortName(), ref.Tag())
		}
	}
	return &out, nil
}

// Program heatpump-oauth-server is an OAuth 2 server for operating a
// home's fan coil units using Google Asistant.
//
// See https://developers.google.com/assistant/smarthome/overview.
package main

import (
	"context"
	"flag"
	"fmt"
	"net/http"
	"os"
	"strconv"

	"cloud.google.com/go/pubsub"
	"github.com/go-oauth2/oauth2/v4/errors"
	"github.com/go-oauth2/oauth2/v4/manage"
	"github.com/go-oauth2/oauth2/v4/models"
	oauthserver "github.com/go-oauth2/oauth2/v4/server"
	"github.com/go-oauth2/oauth2/v4/store"
	"github.com/golang/glog"
	"github.com/gonzojive/heatpump/cloud/google/cloudconfig"
	"github.com/gonzojive/heatpump/cloud/google/server/fulfilment"
	smarthome "github.com/rmrobinson/google-smart-home-action-go"
)

var (
	port = flag.Int("port", 42598, "Local port to use for HTTP server.")
)

func main() {
	flag.Parse()
	if err := run(context.Background()); err != nil {
		glog.Exitf("%v", err)
	}
}

func run(ctx context.Context) error {
	manager := manage.NewDefaultManager()
	// token memory store
	manager.MustTokenStorage(store.NewMemoryTokenStore())

	// client memory store
	clientStore := store.NewClientStore()
	clientStore.Set("424211", &models.Client{
		ID:     "424211",
		Secret: "423151315asdf",
		// Domain of redirect url.
		Domain: "oauth-redirect.googleusercontent.com",
	})
	manager.MapClientStorage(clientStore)

	oauth2Server := oauthserver.NewDefaultServer(manager)
	oauth2Server.SetAllowGetAccessRequest(true)
	oauth2Server.SetClientInfoHandler(oauthserver.ClientFormHandler)

	oauth2Server.SetInternalErrorHandler(func(err error) (re *errors.Response) {
		glog.Infof("Internal Error:", err.Error())
		return
	})

	oauth2Server.SetResponseErrorHandler(func(re *errors.Response) {
		glog.Infof("Response Error: %s", re.Error.Error())
	})
	oauth2Server.SetUserAuthorizationHandler(func(w http.ResponseWriter, r *http.Request) (userID string, err error) {
		return "4242", nil
	})

	http.HandleFunc("/authorize", func(w http.ResponseWriter, r *http.Request) {
		err := oauth2Server.HandleAuthorizeRequest(w, r)
		glog.Infof("handled authorize request got err = %v", err)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}
	})

	http.HandleFunc("/token", func(w http.ResponseWriter, r *http.Request) {
		oauth2Server.HandleTokenRequest(w, r)
	})

	http.HandleFunc("/hello", func(w http.ResponseWriter, r *http.Request) {
		name := os.Getenv("NAME")
		if name == "" {
			name = "World"
		}
		fmt.Fprintf(w, "Hello %s!\n", name)
	})

	psc, err := pubsub.NewClient(ctx, cloudconfig.DefaultParams().GCPProject)
	if err != nil {
		return fmt.Errorf("failed to initialize pubsub client: %w", err)
	}

	svc := fulfilment.NewService(psc)

	http.Handle(smarthome.GoogleFulfillmentPath, svc.GoogleFulfillmentHandler())

	finalPort := *port

	if got := os.Getenv("PORT"); got != "" {
		x, err := strconv.Atoi(got)
		if err != nil {
			glog.Fatalf("bad PORT env var %q", got)
		}
		finalPort = x
	}

	return http.ListenAndServe(fmt.Sprintf(":%d", finalPort), nil)
}

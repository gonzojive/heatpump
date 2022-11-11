// Package httpendpoint is an OAuth 2 server for operating a home's fan coil
// units using Google Asistant.
//
// See https://developers.google.com/assistant/smarthome/overview and
// specifically
// https://developers.google.com/assistant/smarthome/develop/implement-oauth.
package httpendpoint

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"cloud.google.com/go/firestore"
	"github.com/go-oauth2/oauth2/v4/errors"
	"github.com/go-oauth2/oauth2/v4/manage"
	"github.com/go-oauth2/oauth2/v4/models"
	oauthserver "github.com/go-oauth2/oauth2/v4/server"
	"github.com/go-oauth2/oauth2/v4/store"
	"github.com/golang/glog"
	"github.com/gonzojive/heatpump/cloud/google/cloudconfig"
	"github.com/gonzojive/heatpump/cloud/oauthstore"
	"github.com/gonzojive/heatpump/cloud/secrets"
)

const (
	// A value chosen by us to assign to Google's requester here:
	// https://console.actions.google.com/u/0/project/hydronics-9f50d/smarthomeaccountlinking/
	googleClientID = "424211"

	// A client secret is chosen by us to assign to Google's requester here:
	// https://console.actions.google.com/u/0/project/hydronics-9f50d/smarthomeaccountlinking/
	//
	// This is the name of that secret in the Google Cloud Secret Manager.
	googleClientSecretRef secrets.Name = "projects/heatpump-dev/secrets/google-actions-oauth-client-secret/versions/1"

	// Domain of redirect URL to expect from the Google Actions client.
	//
	// TODO(reddaly): Check https.
)

// AccountLinkingServer provides an HTTP server that performs account linking.
//
// See https://developers.google.com/assistant/smarthome/develop/implement-oauth.
type AccountLinkingServer struct {
	oauthServer *oauthserver.Server
}

// NewAccountLinkingServer creates a use *AccountLinkingServer that uses the
// provided context to communicate with GCP APIs like the secret manager.
func NewAccountLinkingServer(ctx context.Context, cloudParams *cloudconfig.Params) (*AccountLinkingServer, error) {
	secret, err := googleClientSecretRef.Access(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to load secret from secret storage: %w", err)
	}
	googleOAuthClient := &models.Client{
		ID:     googleClientID,
		Secret: string(secret),
		// Domain of redirect url.
		Domain: "oauth-redirect.googleusercontent.com",
	}

	// token memory store
	fsClient, err := createFirestoreClient(ctx, cloudParams.GCPProject)
	if err != nil {
		return nil, err
	}

	tokenStore := oauthstore.NewTokenStorage(fsClient, cloudParams.TokenStoreFirebaseCollectionName)
	// client memory store
	clientStore := store.NewClientStore()
	clientStore.Set(googleOAuthClient.ID, googleOAuthClient)

	manager := manage.NewDefaultManager()
	manager.MapTokenStorage(tokenStore)
	manager.MapClientStorage(clientStore)

	service := &AccountLinkingServer{oauthserver.NewDefaultServer(manager)}

	service.oauthServer.SetAllowGetAccessRequest(true)
	service.oauthServer.SetClientInfoHandler(oauthserver.ClientFormHandler)

	service.oauthServer.SetInternalErrorHandler(func(err error) (re *errors.Response) {
		glog.Infof("Internal Error:", err.Error())
		return
	})

	service.oauthServer.SetResponseErrorHandler(func(re *errors.Response) {
		glog.Infof("Response Error: %s", re.Error.Error())
	})

	service.oauthServer.SetUserAuthorizationHandler(func(rw http.ResponseWriter, req *http.Request) (userID string, err error) {
		// Normally, this would be the place where you would check if the user is logged in and gives his consent.
		// We're simplifying things and just checking if the request includes a valid username and password
		if req.Form.Get("secret") != "warmth" {
			rw.Header().Set("Content-Type", "text/html;charset=UTF-8")
			rw.Write([]byte(`<h1>Login</h1>`))
			rw.Write([]byte(`
			<p>What's the secret password?</p>
			<form method="post">
				<input type="password" name="secret" />
				<input type="submit">
			</form>
		`))
			return
		}
		// Return user id.
		return "4242", nil
	})

	return service, nil
}

// RegisterHandlers registers the handlers for account linking described here
// https://developers.google.com/assistant/smarthome/develop/implement-oauth and
// configured at
// https://console.actions.google.com/u/0/project/hydronics-9f50d/actions/smarthome/.
func (s *AccountLinkingServer) RegisterHandlers(mux *http.ServeMux) {
	mux.Handle("/authorize", http.HandlerFunc(s.handleAuthorizeRequest))
	mux.Handle("/token", http.HandlerFunc(s.handleTokenRequest))

	mux.HandleFunc("/hello", func(w http.ResponseWriter, r *http.Request) {
		name := os.Getenv("NAME")
		if name == "" {
			name = "World"
		}
		fmt.Fprintf(w, "Hello %s!\n", name)
	})
}

func (s *AccountLinkingServer) handleAuthorizeRequest(w http.ResponseWriter, r *http.Request) {
	err := s.oauthServer.HandleAuthorizeRequest(w, r)
	glog.Errorf("handled authorize request got err = %v", err)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
}

func (s *AccountLinkingServer) handleTokenRequest(w http.ResponseWriter, r *http.Request) {
	s.oauthServer.HandleTokenRequest(w, r)
}

func createFirestoreClient(ctx context.Context, projectID string) (*firestore.Client, error) {
	// Sets your Google Cloud Platform project ID.
	client, err := firestore.NewClient(ctx, projectID)
	if err != nil {
		return nil, fmt.Errorf("failed to create FireStore client: %v", err)
	}
	// Close client when done with
	// defer client.Close()
	return client, nil
}

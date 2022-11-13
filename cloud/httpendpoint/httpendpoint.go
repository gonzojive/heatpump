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

	"cloud.google.com/go/pubsub"
	"github.com/gonzojive/heatpump/cloud/google/cloudconfig"
	"github.com/gonzojive/heatpump/cloud/google/server/fulfilment"

	cpb "github.com/gonzojive/heatpump/proto/controller"
)

type Server struct {
}

func NewServer(ctx context.Context, mux *http.ServeMux, cloudParams *cloudconfig.Params, stateService cpb.StateServiceClient) (*Server, error) {
	{
		s, err := NewAccountLinkingServer(ctx, cloudParams)
		if err != nil {
			return nil, err
		}
		s.RegisterHandlers(mux)
	}

	if err := registerFulfillmentHandlers(ctx, mux, cloudParams, stateService); err != nil {
		return nil, fmt.Errorf("error bringing up fulfillment server: %w", err)
	}

	return &Server{}, nil
}

// registerFulfillmentHandlers registers the handlers for action fulfillment
// described at
// https://developers.google.com/assistant/smarthome/concepts/fulfillment-authentication
// and configured at
// https://console.actions.google.com/u/0/project/hydronics-9f50d/actions/smarthome/.
func registerFulfillmentHandlers(ctx context.Context, mux *http.ServeMux, cloudParams *cloudconfig.Params, stateService cpb.StateServiceClient) error {
	psc, err := pubsub.NewClient(ctx, cloudParams.GCPProject)
	if err != nil {
		return fmt.Errorf("failed to initialize pubsub client: %w", err)
	}

	svc := fulfilment.NewService(psc, stateService)

	// Fulfillment path configured at
	// https://console.actions.google.com/u/0/project/hydronics-9f50d/actions/smarthome/
	mux.Handle("/fulfillment", svc.GoogleFulfillmentHandler())
	mux.Handle("/debug", svc.DebugHandler())
	return nil
}

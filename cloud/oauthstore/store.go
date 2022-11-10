package oauthstore

import (
	"context"
	"errors"
	"reflect"
	"time"

	"cloud.google.com/go/firestore"
	"github.com/go-oauth2/oauth2/v4"
	"github.com/go-oauth2/oauth2/v4/models"
)

const (
	keyCode     = "Code"
	keyAccess   = "Access"
	keyRefresh  = "Refresh"
	KeyClientID = "ID"

	timeout = 30 * time.Second
)

// NewTokenStorage returns a new Firestore token store.
// The provided firestore client will never be closed.
func NewTokenStorage(c *firestore.Client, collection string) oauth2.TokenStore {
	fs := &store{c: c, n: collection, t: timeout}
	return &client{c: fs}
}

// NewClientStorage returns a new Firestore token store.
// The provided firestore client will never be closed.
func NewClientStorage(c *firestore.Client, collection string) oauth2.ClientStore {
	fs := &store{c: c, n: collection, t: timeout}
	return &client{c: fs}
}

type client struct {
	c *store
}

func (f *client) Create(ctx context.Context, info oauth2.TokenInfo) error {
	t, err := token(info)
	if err != nil {
		return err
	}
	return f.c.Put(ctx, t)
}

func (f *client) RemoveByCode(ctx context.Context, code string) error {
	return f.c.Del(ctx, keyCode, code)
}

func (f *client) RemoveByAccess(ctx context.Context, access string) error {
	return f.c.Del(ctx, keyAccess, access)
}

func (f *client) RemoveByRefresh(ctx context.Context, refresh string) error {
	return f.c.Del(ctx, keyRefresh, refresh)
}

func (f *client) GetByCode(ctx context.Context, code string) (oauth2.TokenInfo, error) {
	return f.c.Get(ctx, keyCode, code)
}

func (f *client) GetByAccess(ctx context.Context, access string) (oauth2.TokenInfo, error) {
	return f.c.Get(ctx, keyAccess, access)
}

func (f *client) GetByRefresh(ctx context.Context, refresh string) (oauth2.TokenInfo, error) {
	return f.c.Get(ctx, keyRefresh, refresh)
}

func (f *client) GetByID(ctx context.Context, id string) (oauth2.ClientInfo, error) {
	return f.c.GetClient(ctx, KeyClientID, id)
}

// ErrInvalidTokenInfo is returned whenever TokenInfo is either nil or zero/empty.
var ErrInvalidTokenInfo = errors.New("invalid TokenInfo")

func token(info oauth2.TokenInfo) (*models.Token, error) {
	if isNilOrZero(info) {
		return nil, ErrInvalidTokenInfo
	}
	return &models.Token{
		ClientID:         info.GetClientID(),
		UserID:           info.GetUserID(),
		RedirectURI:      info.GetRedirectURI(),
		Scope:            info.GetScope(),
		Code:             info.GetCode(),
		CodeCreateAt:     info.GetCodeCreateAt(),
		CodeExpiresIn:    info.GetCodeExpiresIn(),
		Access:           info.GetAccess(),
		AccessCreateAt:   info.GetAccessCreateAt(),
		AccessExpiresIn:  info.GetAccessExpiresIn(),
		Refresh:          info.GetRefresh(),
		RefreshCreateAt:  info.GetRefreshCreateAt(),
		RefreshExpiresIn: info.GetRefreshExpiresIn(),
	}, nil
}

func isNilOrZero(info oauth2.TokenInfo) bool {
	if info == nil {
		return true
	}
	if v := reflect.ValueOf(info); v.IsNil() {
		return true
	}
	return reflect.DeepEqual(info, info.New())
}

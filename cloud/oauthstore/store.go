/*
Copyright (c) 2019 Tadej Slamic

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
*/

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

// NewTokenStorage returns a new Firestore-backed token store.
// The provided firestore client will never be closed.
func NewTokenStorage(c *firestore.Client, collection string) *TokenStore {
	return &TokenStore{&store{c: c, collection: c.Collection(collection), timeout: timeout}}
}

// NewClientStorage returns a new Firestore token store.
// The provided firestore client will never be closed.
func NewClientStorage(c *firestore.Client, collection string) *ClientStore {
	return &ClientStore{&store{c: c, collection: c.Collection(collection), timeout: timeout}}
}

type ClientStore struct {
	st *store
}

var _ oauth2.ClientStore = (*ClientStore)(nil)

func (s *ClientStore) GetByID(ctx context.Context, id string) (oauth2.ClientInfo, error) {
	return s.st.GetClient(ctx, KeyClientID, id)
}

type TokenStore struct {
	c *store
}

var _ oauth2.TokenStore = (*TokenStore)(nil)

func (f *TokenStore) Create(ctx context.Context, info oauth2.TokenInfo) error {
	t, err := token(info)
	if err != nil {
		return err
	}
	return f.c.Put(ctx, t)
}

func (f *TokenStore) RemoveByCode(ctx context.Context, code string) error {
	return f.c.Del(ctx, keyCode, code)
}

func (f *TokenStore) RemoveByAccess(ctx context.Context, access string) error {
	return f.c.Del(ctx, keyAccess, access)
}

func (f *TokenStore) RemoveByRefresh(ctx context.Context, refresh string) error {
	return f.c.Del(ctx, keyRefresh, refresh)
}

func (f *TokenStore) GetByCode(ctx context.Context, code string) (oauth2.TokenInfo, error) {
	return f.c.Get(ctx, keyCode, code)
}

func (f *TokenStore) GetByAccess(ctx context.Context, access string) (oauth2.TokenInfo, error) {
	return f.c.Get(ctx, keyAccess, access)
}

func (f *TokenStore) GetByRefresh(ctx context.Context, refresh string) (oauth2.TokenInfo, error) {
	return f.c.Get(ctx, keyRefresh, refresh)
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

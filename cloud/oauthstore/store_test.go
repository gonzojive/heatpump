package oauthstore

import (
	"context"
	"testing"

	"cloud.google.com/go/firestore"
	"github.com/go-oauth2/oauth2/v4"
	"github.com/go-oauth2/oauth2/v4/models"
	"github.com/stretchr/testify/assert"
	"google.golang.org/api/iterator"
)

var c *firestore.Client

/*
func TestMain(m *testing.M) {
	project, ok := os.LookupEnv("PROJECT_ID")
	if !ok {
		log.Fatalln("PROJECT_ID env variable is missing")
	}
	ctx := context.Background()
	conf := &firebase.Config{ProjectID: project}
	app, err := firebase.NewApp(ctx, conf)
	if err != nil {
		log.Fatalln(err)
	}
	c, err = app.Firestore(ctx)
	if err != nil {
		log.Fatalln(err)
	}
	os.Exit(func() int {
		defer c.Close()
		return m.Run()
	}())
}
*/

func TestStoreClient(t *testing.T) {
	client := NewTokenStorage(c, "tests")
	type holder struct {
		key string
		get func(context.Context, string) (oauth2.TokenInfo, error)
		del func(context.Context, string) error
	}
	tokens := map[*models.Token]holder{
		{Access: "access"}:   {key: "access", get: client.GetByAccess, del: client.RemoveByAccess},
		{Code: "code"}:       {key: "code", get: client.GetByCode, del: client.RemoveByCode},
		{Refresh: "refresh"}: {key: "refresh", get: client.GetByRefresh, del: client.RemoveByRefresh},
	}
	for i, h := range tokens {
		ctx := context.Background()
		err := client.Create(ctx, i)
		assert.Nil(t, err)

		tok, err := h.get(ctx, h.key)
		assert.Nil(t, err)
		assert.Equal(t, i, tok)

		err = h.del(ctx, h.key)
		assert.Nil(t, err)

		_, err = h.get(ctx, h.key)
		assert.NotNil(t, err)

		err = h.del(ctx, h.key)
		assert.Nil(t, err)
	}
}

func TestNoDocument(t *testing.T) {
	client := NewTokenStorage(c, "tests")
	ctx := context.Background()
	info, err := client.GetByRefresh(ctx, "whoops")
	assert.Nil(t, info)
	assert.Equal(t, iterator.Done, err)
}

func TestIsNilOrZero(t *testing.T) {
	tokens := map[oauth2.TokenInfo]bool{
		nil:                               true,
		&models.Token{}:                   true,
		&models.Token{Access: "access"}:   false,
		&models.Token{Code: "code"}:       false,
		&models.Token{Refresh: "refresh"}: false,
	}
	for tok, expected := range tokens {
		result := isNilOrZero(tok)
		assert.Equal(t, expected, result)
	}
}

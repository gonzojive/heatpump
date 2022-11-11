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
	"time"

	"cloud.google.com/go/firestore"
	"github.com/go-oauth2/oauth2/v4/models"
	"google.golang.org/api/iterator"
)

type store struct {
	c          *firestore.Client
	collection *firestore.CollectionRef
	timeout    time.Duration
}

func (s *store) Put(ctx context.Context, token *models.Token) error {
	ctx, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()
	_, _, err := s.collection.Add(ctx, token)
	return err
}

func runInTransaction[T any](ctx context.Context, client *firestore.Client, fn func(ctx context.Context, txn *firestore.Transaction) (T, error), opts ...firestore.TransactionOption) (T, error) {
	var result T
	if err := client.RunTransaction(ctx, func(ctx context.Context, txn *firestore.Transaction) error {
		var err error
		result, err = fn(ctx, txn)
		return err
	}, opts...); err != nil {
		var zero T
		return zero, err
	}
	return result, nil
}

func (s *store) Get(ctx context.Context, keyPath string, keyValue string) (*models.Token, error) {
	ctx, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	return runInTransaction(ctx, s.c, func(ctx context.Context, txn *firestore.Transaction) (*models.Token, error) {
		iter := txn.Documents(s.collection.Where(keyPath, "==", keyValue).Limit(1))
		doc, err := first(iter)
		if err != nil {
			return nil, err
		}
		info := &models.Token{}
		err = doc.DataTo(info)
		return info, err
	})
}

func (s *store) GetClient(ctx context.Context, key string, keyValue string) (*models.Client, error) {
	ctx, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()
	iter := s.collection.Where(key, "==", keyValue).Limit(1).Documents(ctx)
	doc, err := first(iter)
	if err != nil {
		return nil, err
	}
	info := &models.Client{}
	err = doc.DataTo(info)
	return info, err
}

func (s *store) Del(ctx context.Context, key string, keyValue string) error {
	ctx, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()
	return s.c.RunTransaction(ctx, func(ctx context.Context, tx *firestore.Transaction) error {
		query := s.collection.Where(key, "==", keyValue).Limit(1)
		iter := tx.Documents(query)
		doc, err := first(iter)
		if err != nil {
			if err == iterator.Done || err == ErrDocumentDoesNotExist {
				return nil // Document does not exist - we're done!
			}
			return err
		}
		return tx.Delete(doc.Ref)
	})
}

// ErrDocumentDoesNotExist is returned whenever a Firestore document does not exist.
var ErrDocumentDoesNotExist = errors.New("document does not exist")

func first(iter *firestore.DocumentIterator) (*firestore.DocumentSnapshot, error) {
	doc, err := iter.Next()
	if err == iterator.Done {
		return nil, ErrDocumentDoesNotExist
	}
	if err != nil {
		return nil, err
	}
	if !doc.Exists() {
		return nil, ErrDocumentDoesNotExist
	}
	return doc, nil
}

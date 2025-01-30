// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package sdk

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"sync"
	"time"

	"github.com/benbjohnson/clock"
)

type token struct {
	value     string
	expiresAt time.Time
}

type TokenStore struct {
	lock           sync.Mutex
	clock          clock.Clock
	token          *token
	httpClient     HttpClient
	tokenServerURL *url.URL
}

func NewDefaultTokenStore(tokenServerURL string) (*TokenStore, error) {
	return NewTokenStore(clock.New(), http.DefaultClient, tokenServerURL)
}

func NewTokenStore(clock clock.Clock, httpClient HttpClient, tokenServerURL string) (*TokenStore, error) {
	tsu, err := url.Parse(tokenServerURL)
	if err != nil {
		return nil, err
	}
	return &TokenStore{
		lock:           sync.Mutex{},
		clock:          clock,
		httpClient:     httpClient,
		tokenServerURL: tsu,
	}, nil
}

func (ts *TokenStore) Fetch(ctx context.Context, clientID string, clientSecret string) (string, error) {

	ts.lock.Lock()
	defer ts.lock.Unlock()

	if ts.token != nil && ts.clock.Now().Before(ts.token.expiresAt) {
		return ts.token.value, nil
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, ts.tokenServerURL.String(), nil)
	if err != nil {
		return "", err
	}

	res, err := ts.httpClient.Do(req)
	if err != nil {
		return "", err
	}

	if res.StatusCode != http.StatusOK {
		return "", fmt.Errorf("expected %d status code but got %d", http.StatusOK, res.StatusCode)
	}

	defer res.Body.Close()
	tokenResponse := struct {
		AccessToken string `json:"access_token"`
		ExpiresIn   int64  `json:"expires_in"`
	}{}

	if err := json.NewDecoder(res.Body).Decode(&tokenResponse); err != nil {
		return "", err
	}

	ts.token = &token{
		value:     tokenResponse.AccessToken,
		expiresAt: ts.clock.Now().Add(ts.getExpiryDurationFromExpiresIn(tokenResponse.ExpiresIn)),
	}

	return ts.token.value, nil
}

func (ts *TokenStore) getExpiryDurationFromExpiresIn(expiresIn int64) time.Duration {
	return time.Duration(float64(expiresIn)*0.9) * time.Second
}

// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package sdk_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"tasks-terraform-provider/internal/sdk"

	"github.com/benbjohnson/clock"
	"github.com/stretchr/testify/assert"
)

func TestTokenStore(t *testing.T) {

	ctx := context.Background()

	t.Run("test that it returns a cached token", func(t *testing.T) {
		calls := 0
		mockOAuthServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			calls++
			_, _ = w.Write([]byte(`{
				"access_token": "dummy-access-token",
				"expires_in": 3600
			}`))
		}))
		defer mockOAuthServer.Close()

		cl := clock.NewMock()
		ts, err := sdk.NewTokenStore(cl, http.DefaultClient, mockOAuthServer.URL)
		assert.NoError(t, err)

		token, err := ts.Fetch(ctx, "", "")
		assert.NoError(t, err)
		assert.Equal(t, "dummy-access-token", token)

		token, err = ts.Fetch(ctx, "", "")
		assert.NoError(t, err)
		assert.Equal(t, "dummy-access-token", token)

		assert.Equal(t, 1, calls)
	})

	t.Run("test that it makes multiple calls if token expired", func(t *testing.T) {
		calls := 0
		mockOAuthServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			calls++
			_, _ = w.Write([]byte(`{
				"access_token": "dummy-access-token",
				"expires_in": 3600
			}`))
		}))
		defer mockOAuthServer.Close()

		cl := clock.NewMock()
		ts, err := sdk.NewTokenStore(cl, http.DefaultClient, mockOAuthServer.URL)
		assert.NoError(t, err)

		token, err := ts.Fetch(ctx, "", "")
		assert.NoError(t, err)
		assert.Equal(t, "dummy-access-token", token)

		cl.Add(time.Second * 3600)

		token, err = ts.Fetch(ctx, "", "")
		assert.NoError(t, err)
		assert.Equal(t, "dummy-access-token", token)

		assert.Equal(t, 2, calls)
	})
}

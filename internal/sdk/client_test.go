// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package sdk_test

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"tasks-terraform-provider/internal/sdk"

	"github.com/benbjohnson/clock"
	"github.com/stretchr/testify/assert"
)

func TestTasksClient(t *testing.T) {

	ctx := context.Background()

	t.Run("test that it sends auth header", func(t *testing.T) {
		mockOAuthServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			_, _ = w.Write([]byte(`{
				"access_token": "dummy-access-token",
				"expires_in": 3600
			}`))
		}))
		defer mockOAuthServer.Close()

		mockAPIServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "Bearer dummy-access-token", r.Header.Get("Authorization"))
			assert.True(t, strings.HasSuffix(r.RequestURI, "/org/00000001-00000000-00000000-00000000/project/019010a2-8327-2607-11d7-41bb0a8936d4/task/019011e6-e530-3aca-6cf7-2973387c255d"))
			data := getTestFixture(t, "task.json")
			_, _ = w.Write(data)
		}))
		defer mockAPIServer.Close()

		ts, err := sdk.NewTokenStore(clock.New(), http.DefaultClient, mockOAuthServer.URL)
		assert.NoError(t, err)
		tc := sdk.NewTasksClient(http.DefaultClient, ts, sdk.Config{
			RightbrainAPIHost:   mockAPIServer.URL,
			RightbrainOrgID:     "00000001-00000000-00000000-00000000",
			RightbrainProjectID: "019010a2-8327-2607-11d7-41bb0a8936d4",
		})
		_, err = tc.FetchByID(ctx, "019011e6-e530-3aca-6cf7-2973387c255d")
		assert.NoError(t, err)
	})

	t.Run("test that it can fetch a task", func(t *testing.T) {

		var calls int

		mockOAuthServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			calls++
			_, _ = w.Write([]byte(`{
				"access_token": "dummy-access-token",
				"expires_in": 3600
			}`))
		}))
		defer mockOAuthServer.Close()

		mockAPIServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.True(t, strings.HasSuffix(r.RequestURI, "/org/00000001-00000000-00000000-00000000/project/019010a2-8327-2607-11d7-41bb0a8936d4/task/019011e6-e530-3aca-6cf7-2973387c255d"))
			data := getTestFixture(t, "task.json")
			_, _ = w.Write(data)
		}))
		defer mockAPIServer.Close()

		ts, err := sdk.NewTokenStore(clock.New(), http.DefaultClient, mockOAuthServer.URL)
		assert.NoError(t, err)
		tc := sdk.NewTasksClient(http.DefaultClient, ts, sdk.Config{
			RightbrainAPIHost:   mockAPIServer.URL,
			RightbrainOrgID:     "00000001-00000000-00000000-00000000",
			RightbrainProjectID: "019010a2-8327-2607-11d7-41bb0a8936d4",
		})
		task, err := tc.FetchByID(ctx, "019011e6-e530-3aca-6cf7-2973387c255d")
		assert.NoError(t, err)
		assert.Equal(t, 1, calls)
		assert.Equal(t, "019011e6-e530-3aca-6cf7-2973387c255d", task.ID)
	})

	t.Run("test that it sends a create request", func(t *testing.T) {
		mockOAuthServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			_, _ = w.Write([]byte(`{
				"access_token": "dummy-access-token",
				"expires_in": 3600
			}`))
		}))
		defer mockOAuthServer.Close()

		mockAPIServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, http.MethodPost, r.Method)
			assert.Equal(t, "Bearer dummy-access-token", r.Header.Get("Authorization"))
			assert.True(t, strings.HasSuffix(r.RequestURI, "/org/00000001-00000000-00000000-00000000/project/019010a2-8327-2607-11d7-41bb0a8936d4/task"))
			data := getTestFixture(t, "task.json")
			_, _ = w.Write(data)
		}))
		defer mockAPIServer.Close()

		ts, err := sdk.NewTokenStore(clock.New(), http.DefaultClient, mockOAuthServer.URL)
		assert.NoError(t, err)
		tc := sdk.NewTasksClient(http.DefaultClient, ts, sdk.Config{
			RightbrainAPIHost:   mockAPIServer.URL,
			RightbrainOrgID:     "00000001-00000000-00000000-00000000",
			RightbrainProjectID: "019010a2-8327-2607-11d7-41bb0a8936d4",
		})
		in := &sdk.Task{
			Description: "A task to pre-triage user onboarding before IDV.",
		}
		task, err := tc.Create(ctx, in)
		assert.NoError(t, err)
		assert.Equal(t, "019011e6-e530-3aca-6cf7-2973387c255d", task.ID)
	})

	t.Run("test that it sends an update request", func(t *testing.T) {
		mockOAuthServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			_, _ = w.Write([]byte(`{
				"access_token": "dummy-access-token",
				"expires_in": 3600
			}`))
		}))
		defer mockOAuthServer.Close()

		mockAPIServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, http.MethodPost, r.Method)
			assert.Equal(t, "Bearer dummy-access-token", r.Header.Get("Authorization"))
			assert.True(t, strings.HasSuffix(r.RequestURI, "/org/00000001-00000000-00000000-00000000/project/019010a2-8327-2607-11d7-41bb0a8936d4/task/019011e6-e530-3aca-6cf7-2973387c255d"))
			data := getTestFixture(t, "task.json")
			_, _ = w.Write(data)
		}))
		defer mockAPIServer.Close()

		ts, err := sdk.NewTokenStore(clock.New(), http.DefaultClient, mockOAuthServer.URL)
		assert.NoError(t, err)
		tc := sdk.NewTasksClient(http.DefaultClient, ts, sdk.Config{
			RightbrainAPIHost:   mockAPIServer.URL,
			RightbrainOrgID:     "00000001-00000000-00000000-00000000",
			RightbrainProjectID: "019010a2-8327-2607-11d7-41bb0a8936d4",
		})
		in := &sdk.Task{
			ID:          "019011e6-e530-3aca-6cf7-2973387c255d",
			Description: "A task to pre-triage user onboarding before IDV.",
		}
		task, err := tc.Update(ctx, in)
		assert.NoError(t, err)
		assert.Equal(t, "019011e6-e530-3aca-6cf7-2973387c255d", task.ID)
	})
}

// nolint:unparam
func getTestFixture(t *testing.T, fixture string) []byte {
	data, err := os.ReadFile(fmt.Sprintf("fixtures/%s", fixture))
	assert.NoError(t, err)
	return data
}

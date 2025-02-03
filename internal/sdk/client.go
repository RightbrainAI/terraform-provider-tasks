// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package sdk

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

const (
	DefaultAPIVersion = "v1"
)

func NewTasksClient(log Log, httpClient HttpClient, tokenStore *TokenStore, config Config) *TasksClient {
	return &TasksClient{
		log:        log,
		tokenStore: tokenStore,
		httpClient: httpClient,
		config:     config,
	}
}

type TasksClient struct {
	log        Log
	tokenStore *TokenStore
	httpClient HttpClient
	config     Config
}

func (tc *TasksClient) FetchByID(ctx context.Context, taskID string) (*Task, error) {
	url := fmt.Sprintf("%s/task/%s", tc.getBaseAPIURL(), taskID)
	tc.log.Error("fetching task", "id", taskID, "url", url)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	res, err := tc.DoWithAuth(ctx, req)
	if err != nil {
		return nil, err
	}
	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("cannot fetch task, expected status code %d, but got %d", http.StatusOK, res.StatusCode)
	}
	task := new(Task)

	if err := json.NewDecoder(res.Body).Decode(&task); err != nil {
		return nil, err
	}

	return task, nil
}

func (tc *TasksClient) Create(ctx context.Context, in *CreateTaskRequest) (*Task, error) {
	var data = new(bytes.Buffer)
	if err := json.NewEncoder(data).Encode(&in); err != nil {
		return nil, err
	}
	url := fmt.Sprintf("%s/task", tc.getBaseAPIURL())
	tc.log.Info("creating task", "url", url)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, data)
	if err != nil {
		return nil, err
	}
	res, err := tc.DoWithAuth(ctx, req)
	if err != nil {
		return nil, err
	}
	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("cannot create task, expected status code %d, but got %d", http.StatusOK, res.StatusCode)
	}
	task := new(Task)
	if err := json.NewDecoder(res.Body).Decode(&task); err != nil {
		return nil, err
	}
	return task, nil
}

func (tc *TasksClient) Update(ctx context.Context, in *Task) (*Task, error) {
	var data = new(bytes.Buffer)
	if err := json.NewEncoder(data).Encode(&in); err != nil {
		return nil, err
	}
	url := fmt.Sprintf("%s/task/%s", tc.getBaseAPIURL(), in.ID)
	tc.log.Error("updating task", "id", in.ID, "url", url)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, data)
	if err != nil {
		return nil, err
	}
	res, err := tc.DoWithAuth(ctx, req)
	if err != nil {
		return nil, err
	}
	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("cannot update task, expected status code %d, but got %d", http.StatusOK, res.StatusCode)
	}
	task := new(Task)
	if err := json.NewDecoder(res.Body).Decode(&task); err != nil {
		return nil, err
	}
	return task, nil
}

func (tc *TasksClient) DeleteByID(ctx context.Context, taskID string) error {
	url := fmt.Sprintf("%s/task/%s", tc.getBaseAPIURL(), taskID)
	tc.log.Info("deleting task", "id", taskID, "url", url)

	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, url, nil)
	if err != nil {
		return err
	}
	res, err := tc.DoWithAuth(ctx, req)
	if err != nil {
		return err
	}
	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("cannot delete task, expected status code %d, but got %d", http.StatusOK, res.StatusCode)
	}
	return nil
}

func (tc *TasksClient) DoWithAuth(ctx context.Context, req *http.Request) (*http.Response, error) {
	token, err := tc.tokenStore.Fetch(ctx, tc.config.RightbrainClientID, tc.config.RightbrainClientSecret)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	return tc.httpClient.Do(req)
}

func (tc *TasksClient) getBaseAPIURL() string {
	return fmt.Sprintf("%s/api/%s/org/%s/project/%s", tc.config.RightbrainAPIHost, DefaultAPIVersion, tc.config.RightbrainOrgID, tc.config.RightbrainProjectID)
}

// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package sdk

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	entities "terraform-provider-tasks/internal/sdk/entities"
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

func (tc *TasksClient) Fetch(ctx context.Context, in FetchTaskRequest) (*entities.Task, error) {
	url := fmt.Sprintf("%s/task/%s", tc.getBaseAPIURL(), in.ID)
	tc.log.Info("fetching task", "id", in.ID, "url", url)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		tc.log.Error(err.Error())
		return nil, err
	}
	res, err := tc.DoWithAuth(ctx, req)
	if err != nil {
		tc.log.Error(err.Error())
		return nil, err
	}
	if err := tc.assertStatusCode("cannot fetch task", http.StatusOK, res); err != nil {
		tc.log.Error(err.Error())
		return nil, err
	}

	task := new(entities.Task)

	if err := json.NewDecoder(res.Body).Decode(&task); err != nil {
		tc.log.Error(err.Error())
		return nil, err
	}

	return task, nil
}

func (tc *TasksClient) Create(ctx context.Context, in CreateTaskRequest) (*entities.Task, error) {
	var data = new(bytes.Buffer)
	if err := json.NewEncoder(data).Encode(&in); err != nil {
		return nil, err
	}
	url := fmt.Sprintf("%s/task", tc.getBaseAPIURL())
	tc.log.Info("creating task", "url", url)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, data)
	if err != nil {
		tc.log.Error(err.Error())
		return nil, err
	}
	res, err := tc.DoWithAuth(ctx, req)
	if err != nil {
		tc.log.Error(err.Error())
		return nil, err
	}
	if err := tc.assertStatusCode("cannot create task", http.StatusOK, res); err != nil {
		tc.log.Error(err.Error())
		return nil, err
	}
	task := new(entities.Task)
	if err := json.NewDecoder(res.Body).Decode(&task); err != nil {
		tc.log.Error(err.Error())
		return nil, err
	}
	return task, nil
}

func (tc *TasksClient) Update(ctx context.Context, in UpdateTaskRequest) (*entities.Task, error) {
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
	if err := tc.assertStatusCode("cannot update task", http.StatusOK, res); err != nil {
		tc.log.Error(err.Error())
		return nil, err
	}
	task := new(entities.Task)
	if err := json.NewDecoder(res.Body).Decode(&task); err != nil {
		tc.log.Error(err.Error())
		return nil, err
	}
	if err := tc.markLatestTaskRevisionAsActive(ctx, task); err != nil {
		tc.log.Error(err.Error())
		return nil, err
	}
	return tc.Fetch(ctx, NewFetchTaskRequest(in.ID))
}

func (tc *TasksClient) Delete(ctx context.Context, in DeleteTaskRequest) error {
	url := fmt.Sprintf("%s/task/%s", tc.getBaseAPIURL(), in.ID)
	tc.log.Info("deleting task", "id", in.ID, "url", url)

	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, url, nil)
	if err != nil {
		return err
	}
	res, err := tc.DoWithAuth(ctx, req)
	if err != nil {
		return err
	}
	if err := tc.assertStatusCode("cannot delete task", http.StatusOK, res); err != nil {
		tc.log.Error(err.Error())
		return err
	}
	return nil
}

func (tc *TasksClient) GetAvailableLLMModels(ctx context.Context) ([]entities.Model, error) {
	url := fmt.Sprintf("%s/model", tc.getBaseAPIURL())
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	res, err := tc.DoWithAuth(ctx, req)
	if err != nil {
		return nil, err
	}

	if err := tc.assertStatusCode("cannot obtain model list", http.StatusOK, res); err != nil {
		tc.log.Error(err.Error())
		return nil, err
	}

	defer res.Body.Close()

	var models []entities.Model

	if err := json.NewDecoder(res.Body).Decode(&models); err != nil {
		return nil, err
	}

	return models, nil
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

func (tc *TasksClient) markLatestTaskRevisionAsActive(ctx context.Context, task *entities.Task) error {
	url := fmt.Sprintf("%s/task/%s", tc.getBaseAPIURL(), task.ID)
	tc.log.Info("updating task", "id", task.ID, "url", url)

	rev, err := task.GetLatestRevision()
	if err != nil {
		return err
	}

	data := fmt.Sprintf(`{
		"active_revisions": [
			{
				"weight": 1,
				"task_revision_id": %q
			}
		]
	}`, rev.ID)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader([]byte(data)))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	res, err := tc.DoWithAuth(ctx, req)
	if err != nil {
		return err
	}
	if err := tc.assertStatusCode("cannot make revision active", http.StatusOK, res); err != nil {
		tc.log.Error(err.Error())
		return err
	}
	return nil
}

func (tc *TasksClient) assertStatusCode(prefix string, expected int, res *http.Response) error { //nolint:unparam

	if res.StatusCode == expected {
		return nil
	}

	body, _ := io.ReadAll(res.Body)

	tc.log.Info("status code was not as expected", "expected", expected, "got", res.StatusCode, "body", string(body))

	return fmt.Errorf("%s, expected status code %d but got %d.", prefix, expected, res.StatusCode)
}

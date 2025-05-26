// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package sdk

import (
	entitites "terraform-provider-tasks/internal/sdk/entities"
)

type CreateTaskRequest struct {
	Description     string                      `json:"description"`
	Enabled         bool                        `json:"enabled"`
	ImageRequired   bool                        `json:"image_required"`
	InputProcessors *[]entitites.InputProcessor `json:"input_processors"`
	LLMModelID      string                      `json:"llm_model_id"`
	Name            string                      `json:"name"`
	OutputFormat    map[string]string           `json:"output_format"`
	Public          bool                        `json:"public"`
	SystemPrompt    string                      `json:"system_prompt"`
	UserPrompt      string                      `json:"user_prompt"`
}

func NewCreateTaskRequest() CreateTaskRequest {
	return CreateTaskRequest{
		OutputFormat: make(map[string]string, 0),
	}
}

type UpdateTaskRequest struct {
	ID              string                      `json:"id"`
	Description     string                      `json:"description"`
	Enabled         bool                        `json:"enabled"`
	ImageRequired   bool                        `json:"image_required"`
	InputProcessors *[]entitites.InputProcessor `json:"input_processors"`
	LLMModelID      string                      `json:"llm_model_id"`
	Name            string                      `json:"name"`
	OutputFormat    map[string]string           `json:"output_format"`
	Public          bool                        `json:"public"`
	SystemPrompt    string                      `json:"system_prompt"`
	UserPrompt      string                      `json:"user_prompt"`
}

func NewUpdateTaskRequest(id string) UpdateTaskRequest {
	return UpdateTaskRequest{
		ID:           id,
		OutputFormat: make(map[string]string, 0),
	}
}

type DeleteTaskRequest struct {
	ID string `json:"id"`
}

func NewDeleteTaskRequest(id string) DeleteTaskRequest {
	return DeleteTaskRequest{
		ID: id,
	}
}

type FetchTaskRequest struct {
	ID string `json:"id"`
}

func NewFetchTaskRequest(id string) FetchTaskRequest {
	return FetchTaskRequest{
		ID: id,
	}
}

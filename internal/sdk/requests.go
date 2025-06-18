// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package sdk

import (
	entities "terraform-provider-tasks/internal/sdk/entities"
)

type CreateTaskRequest struct {
	Description     string                                   `json:"description,omitempty"`
	Enabled         bool                                     `json:"enabled,omitempty"`
	ImageRequired   bool                                     `json:"image_required,omitempty"`
	InputProcessors *[]entities.InputProcessor               `json:"input_processors,omitempty"`
	LLMModelID      string                                   `json:"llm_model_id,omitempty"`
	Name            string                                   `json:"name,omitempty"`
	OutputFormat    map[string]entities.OutputFormatExtended `json:"output_format,omitempty"`
	OutputModality  string                                   `json:"output_modality,omitempty"`
	Public          bool                                     `json:"public,omitempty"`
	SystemPrompt    string                                   `json:"system_prompt,omitempty"`
	UserPrompt      string                                   `json:"user_prompt,omitempty"`
}

func NewCreateTaskRequest() CreateTaskRequest {
	return CreateTaskRequest{
		OutputFormat: make(map[string]entities.OutputFormatExtended),
	}
}

type UpdateTaskRequest struct {
	ID              string                                   `json:"id,omitempty"`
	Description     string                                   `json:"description,omitempty"`
	Enabled         bool                                     `json:"enabled,omitempty"`
	ImageRequired   bool                                     `json:"image_required,omitempty"`
	InputProcessors *[]entities.InputProcessor               `json:"input_processors,omitempty"`
	LLMModelID      string                                   `json:"llm_model_id,omitempty"`
	Name            string                                   `json:"name,omitempty"`
	OutputFormat    map[string]entities.OutputFormatExtended `json:"output_format,omitempty"`
	OutputModality  string                                   `json:"output_modality,omitempty"`
	Public          bool                                     `json:"public,omitempty"`
	SystemPrompt    string                                   `json:"system_prompt,omitempty"`
	UserPrompt      string                                   `json:"user_prompt,omitempty"`
}

func NewUpdateTaskRequest(id string) UpdateTaskRequest {
	return UpdateTaskRequest{
		ID:           id,
		OutputFormat: make(map[string]entities.OutputFormatExtended),
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

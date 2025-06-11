// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package sdk

import (
	entities "terraform-provider-tasks/internal/sdk/entities"
)

type BaseTaskProperties struct {
	Description     string                                  `json:"description"`
	Enabled         bool                                    `json:"enabled"`
	ExposedToAgents bool                                    `json:"exposed_to_agents"`
	ImageRequired   bool                                    `json:"image_required"`
	InputProcessors *[]entities.InputProcessor              `json:"input_processors"`
	LLMModelID      string                                  `json:"llm_model_id"`
	Name            string                                  `json:"name"`
	OptimiseImages  bool                                    `json:"optimise_images"`
	OutputFormat    map[string]entities.OutputFormatWrapper `json:"output_format"`
	OutputModality  string                                  `json:"output_modality"`
	Public          bool                                    `json:"public"`
	SystemPrompt    string                                  `json:"system_prompt"`
	UserPrompt      string                                  `json:"user_prompt"`
}

type CreateTaskRequest struct {
	BaseTaskProperties
}

func NewCreateTaskRequest() CreateTaskRequest {
	req := CreateTaskRequest{}
	req.OutputFormat = make(map[string]entities.OutputFormatWrapper, 0)
	return req
}

type UpdateTaskRequest struct {
	BaseTaskProperties
	ID string `json:"id"`
}

func NewUpdateTaskRequest(id string) UpdateTaskRequest {
	req := UpdateTaskRequest{}
	req.ID = id
	req.OutputFormat = make(map[string]entities.OutputFormatWrapper, 0)
	return req
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

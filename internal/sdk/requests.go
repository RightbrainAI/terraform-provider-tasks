// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package sdk

type CreateTaskRequest struct {
	Description   string            `json:"description"`
	Enabled       bool              `json:"enabled"`
	ImageRequired bool              `json:"image_required"`
	LLMModelID    string            `json:"llm_model_id"`
	Name          string            `json:"name"`
	OutputFormat  map[string]string `json:"output_format"`
	Public        bool              `json:"public"`
	SystemPrompt  string            `json:"system_prompt"`
	UserPrompt    string            `json:"user_prompt"`
}

func NewCreateTaskRequest() *CreateTaskRequest {
	return &CreateTaskRequest{
		OutputFormat: make(map[string]string, 0),
	}
}

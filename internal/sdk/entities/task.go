// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package entitites

import "fmt"

// Root represents the overall response structure.
type Task struct {
	Name        string     `json:"name"`
	Enabled     bool       `json:"enabled"`
	Public      bool       `json:"public"`
	ID          string     `json:"id"`
	ProjectID   string     `json:"project_id"`
	Revisions   []Revision `json:"revisions"`
	Description string     `json:"description"`
	AccessToken string     `json:"access_token"`
}

func (t *Task) GetActiveRevision() (*Revision, error) {
	for _, r := range t.Revisions {
		if r.Active {
			return &r, nil
		}
	}
	return nil, fmt.Errorf("could not find active revision for task")
}

func (t *Task) GetLatestRevision() (*Revision, error) {
	if len(t.Revisions) == 0 {
		return nil, fmt.Errorf("could not find latest revision for task")
	}
	return &t.Revisions[0], nil
}

// Revision represents a single revision in the "revisions" array.
type Revision struct {
	SystemPrompt    string            `json:"system_prompt"`
	UserPrompt      string            `json:"user_prompt"`
	LLMModelID      string            `json:"llm_model_id"`
	OutputFormat    OutputFormat      `json:"output_format"`
	ID              string            `json:"id"`
	InputParams     []string          `json:"input_params"`
	TaskForwarderID string            `json:"task_forwarder_id"`
	ImageRequired   bool              `json:"image_required"`
	Active          bool              `json:"active"`
	RAG             RAG               `json:"rag"`
	InputProcessors *[]InputProcessor `json:"input_processors"`
	OutputModality  string            `json:"output_modality"`
}

func (r *Revision) HasInputProcessors() bool {
	return r.InputProcessors != nil && len(*r.InputProcessors) > 0
}

type InputProcessor struct {
	ParamName      string            `json:"param_name"`
	InputProcessor string            `json:"input_processor"`
	Config         map[string]string `json:"config"`
}

// OutputFormat represents the structure of the output format in a revision.
type OutputFormat struct {
	Compliance  string `json:"compliance"`
	Hint        Hint   `json:"hint"`
	Match       Match  `json:"match"`
	Description string `json:"description"`
	Rationale   string `json:"rationale"`
}

// RAG represents the RAG parameters in a revision.
type RAG struct {
	CollectionID string `json:"collection_id"`
	RAGParam     string `json:"rag_param"`
}

type Hint struct {
	Type string `json:"type"` // The type of the hint (e.g., "str")
}

type Match struct {
	Type string `json:"type"` // The type of the match (e.g., "str")
}

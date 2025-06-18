// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package entities

import (
	"encoding/json"
	"fmt"
)

// Revision represents a single revision in the "revisions" array.
type Revision struct {
	Active          bool                           `json:"active,omitempty"`
	ID              string                         `json:"id,omitempty"`
	ImageRequired   bool                           `json:"image_required,omitempty"`
	InputParams     []string                       `json:"input_params,omitempty"`
	InputProcessors *[]InputProcessor              `json:"input_processors,omitempty"`
	LLMModelID      string                         `json:"llm_model_id,omitempty"`
	OptimiseImages  bool                           `json:"optimise_images,omitempty"`
	OutputFormat    map[string]OutputFormatWrapper `json:"output_format,omitempty"`
	OutputModality  string                         `json:"output_modality,omitempty"`
	RAG             RAG                            `json:"rag,omitempty"`
	SystemPrompt    string                         `json:"system_prompt,omitempty"`
	TaskForwarderID string                         `json:"task_forwarder_id,omitempty"`
	UserPrompt      string                         `json:"user_prompt,omitempty"`
}

func (r *Revision) HasInputProcessors() bool {
	return r.InputProcessors != nil && len(*r.InputProcessors) > 0
}

type InputProcessor struct {
	ParamName      string            `json:"param_name,omitempty"`
	InputProcessor string            `json:"input_processor,omitempty"`
	Config         map[string]string `json:"config,omitempty"`
}

type OutputFormatSimple string

func (ofs *OutputFormatSimple) String() string {
	return string(*ofs)
}

// OutputFormat represents the structure of the output format in a revision.
type OutputFormatExtended struct {
	Description string            `json:"description,omitempty"`
	Object      map[string]string `json:"object,omitempty"`
	Options     map[string]string `json:"options,omitempty"`
	Type        string            `json:"type,omitempty"`
	ItemType    string            `json:"item_type,omitempty"`
}

type OutputFormatWrapper struct {
	Extended *OutputFormatExtended
	Simple   *OutputFormatSimple
}

func (ofw OutputFormatWrapper) IsSimple() bool {
	return ofw.Simple != nil
}

func (ofw OutputFormatWrapper) IsExtended() bool {
	return ofw.Extended != nil
}

func (ofw *OutputFormatWrapper) UnmarshalJSON(data []byte) error {
	var simpleOutputFormat *OutputFormatSimple
	if err := json.Unmarshal(data, &simpleOutputFormat); err == nil {
		ofw.Simple = simpleOutputFormat
		return nil
	}
	var extendedOutputFormat *OutputFormatExtended
	if err := json.Unmarshal(data, &extendedOutputFormat); err == nil {
		ofw.Extended = extendedOutputFormat
		return nil
	}
	return fmt.Errorf("cannot unmarshall OutputFormatWrapper")
}

func (ofw *OutputFormatWrapper) MarshalJSON() ([]byte, error) {
	if ofw.IsSimple() {
		return json.Marshal(ofw.Simple)
	}
	if ofw.IsExtended() {
		return json.Marshal(ofw.Extended)
	}
	return nil, fmt.Errorf("cannot marshall OutputFormatWrapper")
}

// RAG represents the RAG parameters in a revision.
type RAG struct {
	CollectionID string `json:"collection_id,omitempty"`
	RAGParam     string `json:"rag_param,omitempty"`
}

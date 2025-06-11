// Copyright (c) HashiCorp, Inc.

package entities

import (
	"encoding/json"
	"fmt"
)

// Revision represents a single revision in the "revisions" array.
type Revision struct {
	Active          bool                           `json:"active"`
	ID              string                         `json:"id"`
	ImageRequired   bool                           `json:"image_required"`
	InputParams     []string                       `json:"input_params"`
	InputProcessors *[]InputProcessor              `json:"input_processors"`
	LLMModelID      string                         `json:"llm_model_id"`
	OptimiseImages  bool                           `json:"optimise_images"`
	OutputFormat    map[string]OutputFormatWrapper `json:"output_format"`
	OutputModality  string                         `json:"output_modality"`
	RAG             RAG                            `json:"rag"`
	SystemPrompt    string                         `json:"system_prompt"`
	TaskForwarderID string                         `json:"task_forwarder_id"`
	UserPrompt      string                         `json:"user_prompt"`
}

func (r *Revision) HasInputProcessors() bool {
	return r.InputProcessors != nil && len(*r.InputProcessors) > 0
}

type InputProcessor struct {
	ParamName      string            `json:"param_name"`
	InputProcessor string            `json:"input_processor"`
	Config         map[string]string `json:"config"`
}

type OutputFormatWrapper struct {
	Complex *OutputFormatComplex
	Simple  *OutputFormatSimple
}

func (ofw *OutputFormatWrapper) IsSimple() bool {
	return ofw.Simple != nil
}

func (ofw *OutputFormatWrapper) IsComplex() bool {
	return ofw.Complex != nil
}

type OutputFormatSimple string

func (ofs *OutputFormatSimple) String() string {
	return string(*ofs)
}

type OutputFormatComplex struct {
	Description string            `json:"description"`
	Object      map[string]string `json:"object"`
	Options     map[string]string `json:"options"`
	Type        string            `json:"type"`
	ItemType    string            `json:"item_type"`
}

func (ofw *OutputFormatWrapper) UnmarshalJSON(data []byte) error {
	var simpleOutputFormat *OutputFormatSimple
	if err := json.Unmarshal(data, &simpleOutputFormat); err == nil {
		ofw.Simple = simpleOutputFormat
		return nil
	}
	var complexOutputFormat *OutputFormatComplex
	if err := json.Unmarshal(data, &complexOutputFormat); err == nil {
		ofw.Complex = complexOutputFormat
		return nil
	}
	return fmt.Errorf("cannot unmarshall OutputFormatWrapper")
}

func (ofw *OutputFormatWrapper) MarshalJSON() ([]byte, error) {
	if ofw.IsSimple() {
		return json.Marshal(ofw.Simple)
	}
	if ofw.IsComplex() {
		return json.Marshal(ofw.Complex)
	}
	return nil, fmt.Errorf("cannot marshall OutputFormatWrapper")
}

type OutputFormatType struct {
	Type string `json:"type"` // The type of the hint (e.g., "str")
}

// RAG represents the RAG parameters in a revision.
type RAG struct {
	CollectionID string `json:"collection_id"`
	RAGParam     string `json:"rag_param"`
}

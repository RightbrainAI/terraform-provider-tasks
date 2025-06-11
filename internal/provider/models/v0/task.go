// Copyright (c) HashiCorp, Inc.

package v0

import (
	"terraform-provider-tasks/internal/sdk/entities"

	"github.com/hashicorp/terraform-plugin-framework/types"
)

type InputProcessorsModel struct {
	InputProcessors []InputProcessorModel `tfsdk:"input_processor"`
}

type InputProcessorModel struct {
	ParamName      types.String            `tfsdk:"param_name"`
	InputProcessor types.String            `tfsdk:"input_processor"`
	Config         map[string]types.String `tfsdk:"config"`
}

// TaskResourceModel describes the resource data model.
type TaskResourceModel struct {
	ActiveRevisionID types.String            `tfsdk:"active_revision_id"`
	Description      types.String            `tfsdk:"description"`
	Enabled          types.Bool              `tfsdk:"enabled"`
	ExposedToAgents  types.Bool              `tfsdk:"exposed_to_agents"`
	ID               types.String            `tfsdk:"id"`
	ImageRequired    types.Bool              `tfsdk:"image_required"`
	InputProcessors  *InputProcessorsModel   `tfsdk:"input_processors"`
	LLMModelID       types.String            `tfsdk:"llm_model_id"`
	Name             types.String            `tfsdk:"name"`
	OptimiseImages   types.Bool              `tfsdk:"optimise_images"`
	OutputFormat     map[string]types.String `tfsdk:"output_format"`
	OutputModality   types.String            `tfsdk:"output_modality"`
	Public           types.Bool              `tfsdk:"public"`
	SystemPrompt     types.String            `tfsdk:"system_prompt"`
	UserPrompt       types.String            `tfsdk:"user_prompt"`
}

func (trm *TaskResourceModel) HasInputProcessors() bool {
	return trm.InputProcessors != nil && len(trm.InputProcessors.InputProcessors) > 0
}

func (trm *TaskResourceModel) PopulateFromTaskEntity(task *entities.Task) error {

	rev, err := task.GetActiveRevision()
	if err != nil {
		return err
	}

	trm.ID = types.StringValue(task.ID)
	trm.Name = types.StringValue(task.Name)
	trm.Enabled = types.BoolValue(task.Enabled)
	trm.ExposedToAgents = types.BoolValue(task.ExposedToAgents)
	trm.Public = types.BoolValue(task.Public)
	trm.Description = types.StringValue(task.Description)

	trm.OptimiseImages = types.BoolValue(rev.OptimiseImages)
	trm.SystemPrompt = types.StringValue(rev.SystemPrompt)
	trm.UserPrompt = types.StringValue(rev.UserPrompt)
	trm.LLMModelID = types.StringValue(rev.LLMModelID)
	trm.ImageRequired = types.BoolValue(rev.ImageRequired)
	trm.OptimiseImages = types.BoolValue(rev.OptimiseImages)

	if rev.HasInputProcessors() {
		trm.InputProcessors = &InputProcessorsModel{
			InputProcessors: make([]InputProcessorModel, len(*rev.InputProcessors)),
		}
		for i, ip := range *rev.InputProcessors {
			trm.InputProcessors.InputProcessors[i] = InputProcessorModel{
				ParamName:      types.StringValue(ip.ParamName),
				InputProcessor: types.StringValue(ip.InputProcessor),
				Config:         make(map[string]types.String, len(ip.Config)),
			}
			for k, v := range ip.Config {
				trm.InputProcessors.InputProcessors[i].Config[k] = types.StringValue(v)
			}
		}
	}

	trm.ActiveRevisionID = types.StringValue(rev.ID)

	return nil
}

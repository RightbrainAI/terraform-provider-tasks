// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"fmt"
	v0_models "terraform-provider-tasks/internal/provider/models/v0"
	v1_models "terraform-provider-tasks/internal/provider/models/v1"
	v1_schemas "terraform-provider-tasks/internal/provider/schemas/v1"
	"terraform-provider-tasks/internal/sdk"
	entities "terraform-provider-tasks/internal/sdk/entities"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.Resource = &TaskResource{}
var _ resource.ResourceWithImportState = &TaskResource{}

func NewTaskResource() resource.Resource {
	return &TaskResource{}
}

// TaskResource defines the resource implementation.
type TaskResource struct {
	client *sdk.TasksClient
}

func (r *TaskResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_task"
}

func (r *TaskResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = v1_schemas.TASK_SCHEMA
}

func (r *TaskResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*sdk.TasksClient)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *sdk.TasksClient, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	r.client = client
}

func (r *TaskResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data v1_models.TaskResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	in := sdk.NewCreateTaskRequest()
	in.Description = data.Description.ValueString()
	in.Enabled = data.Enabled.ValueBool()
	in.ExposedToAgents = data.Enabled.ValueBool()
	in.ImageRequired = data.ImageRequired.ValueBool()
	in.LLMModelID = data.LLMModelID.ValueString()
	in.Name = data.Name.ValueString()
	in.OptimiseImages = data.OptimiseImages.ValueBool()
	in.OutputFormat = data.OutputFormat
	in.OutputModality = data.OutputModality.ValueString()
	in.Public = data.Public.ValueBool()
	in.SystemPrompt = data.SystemPrompt.ValueString()
	in.UserPrompt = data.UserPrompt.ValueString()

	in.InputProcessors = r.FormatInputProcessors(data)

	task, err := r.client.Create(ctx, in)
	if err != nil {
		resp.Diagnostics.AddError(err.Error(), "")
		return
	}

	if err := data.PopulateFromTaskEntity(task); err != nil {
		resp.Diagnostics.AddError(err.Error(), "")
		return
	}

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *TaskResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data v1_models.TaskResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	task, err := r.client.Fetch(ctx, sdk.NewFetchTaskRequest(data.ID.ValueString()))
	if err != nil {
		resp.Diagnostics.AddError(err.Error(), "")
		return
	}

	if err := data.PopulateFromTaskEntity(task); err != nil {
		resp.Diagnostics.AddError(err.Error(), "")
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *TaskResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data v1_models.TaskResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	in := sdk.NewUpdateTaskRequest(data.ID.ValueString())
	in.Description = data.Description.ValueString()
	in.Enabled = data.Enabled.ValueBool()
	in.ExposedToAgents = data.Enabled.ValueBool()
	in.ImageRequired = data.ImageRequired.ValueBool()
	in.LLMModelID = data.LLMModelID.ValueString()
	in.Name = data.Name.ValueString()
	in.OptimiseImages = data.OptimiseImages.ValueBool()
	in.OutputFormat = data.OutputFormat
	in.OutputModality = data.OutputModality.ValueString()
	in.Public = data.Public.ValueBool()
	in.SystemPrompt = data.SystemPrompt.ValueString()
	in.UserPrompt = data.UserPrompt.ValueString()

	in.InputProcessors = r.FormatInputProcessors(data)

	task, err := r.client.Update(ctx, in)
	if err != nil {
		resp.Diagnostics.AddError(err.Error(), "")
		return
	}

	if err := data.PopulateFromTaskEntity(task); err != nil {
		resp.Diagnostics.AddError(err.Error(), "")
		return
	}

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *TaskResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data v1_models.TaskResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.Delete(ctx, sdk.NewDeleteTaskRequest(data.ID.ValueString()))
	if err != nil {
		resp.Diagnostics.AddError(err.Error(), "")
		return
	}
}

func (r *TaskResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func (r *TaskResource) FormatInputProcessors(data v1_models.TaskResourceModel) *[]entities.InputProcessor {
	ips := []entities.InputProcessor{}
	if !data.HasInputProcessors() {
		return &ips
	}
	for _, v := range data.InputProcessors.InputProcessors {
		ip := entities.InputProcessor{
			ParamName:      v.ParamName.ValueString(),
			InputProcessor: v.InputProcessor.ValueString(),
		}
		for k, v := range v.Config {
			if ip.Config == nil {
				ip.Config = make(map[string]string)
			}
			ip.Config[k] = v.ValueString()
		}
		ips = append(ips, ip)
	}
	return &ips
}

func (r *TaskResource) UpgradeState(ctx context.Context) map[int64]resource.StateUpgrader {
	return map[int64]resource.StateUpgrader{
		// State upgrade implementation from 0 (prior state version) to 1 (Schema.Version)
		0: {
			PriorSchema: &v1_schemas.TASK_SCHEMA,
			StateUpgrader: func(ctx context.Context, req resource.UpgradeStateRequest, resp *resource.UpgradeStateResponse) {

				var priorStateData = v0_models.TaskResourceModel{
					OutputFormat: make(map[string]types.String, 0),
				}

				resp.Diagnostics.Append(req.State.Get(ctx, &priorStateData)...)

				if resp.Diagnostics.HasError() {
					return
				}

				// Upgrade logic
				newOutputFormat := make(map[string]entities.OutputFormatWrapper)
				for key, val := range priorStateData.OutputFormat {
					ofs := entities.OutputFormatSimple(val.ValueString())
					newOutputFormat[key] = entities.OutputFormatWrapper{
						Simple: &ofs,
					}
				}

				upgradedStateData := v1_models.TaskResourceModel{
					OutputFormat: newOutputFormat,
				}

				resp.Diagnostics.Append(resp.State.Set(ctx, upgradedStateData)...)
			},
		},
	}
}

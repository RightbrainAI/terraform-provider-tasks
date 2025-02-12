// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"fmt"

	"terraform-provider-tasks/internal/sdk"
	entitites "terraform-provider-tasks/internal/sdk/entities"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
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

// TaskResourceModel describes the resource data model.
type TaskResourceModel struct {
	ID          types.String `tfsdk:"id"`
	Name        types.String `tfsdk:"name"`
	Enabled     types.Bool   `tfsdk:"enabled"`
	Public      types.Bool   `tfsdk:"public"`
	Description types.String `tfsdk:"description"`

	SystemPrompt  types.String            `tfsdk:"system_prompt"`
	UserPrompt    types.String            `tfsdk:"user_prompt"`
	LLMModelID    types.String            `tfsdk:"llm_model_id"`
	ImageRequired types.Bool              `tfsdk:"image_required"`
	OutputFormat  map[string]types.String `tfsdk:"output_format"`

	ActiveRevisionID types.String `tfsdk:"active_revision_id"`
}

func (trm *TaskResourceModel) PopulateFromTaskModel(task *entitites.Task) error {

	rev, err := task.GetActiveRevision()
	if err != nil {
		return err
	}

	trm.ID = types.StringValue(task.ID)
	trm.Name = types.StringValue(task.Name)
	trm.Description = types.StringValue(task.Description)
	trm.ActiveRevisionID = types.StringValue(rev.ID)
	trm.LLMModelID = types.StringValue(rev.LLMModelID)
	trm.SystemPrompt = types.StringValue(rev.SystemPrompt)
	trm.UserPrompt = types.StringValue(rev.UserPrompt)

	return nil
}

func (r *TaskResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_task"
}

func (r *TaskResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "Task resource",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Identifier",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Required:    true,
				Description: "A name or reference for the Task.",
			},
			"description": schema.StringAttribute{
				Optional:    true,
				Description: "A description of the Task.",
			},
			"enabled": schema.BoolAttribute{
				Required:    true,
				Description: "When `true` the Task is active and callable.",
			},
			"public": schema.BoolAttribute{
				Optional:    true,
				Description: "",
			},
			"active_revision_id": schema.StringAttribute{
				Computed: true,
			},

			// revision specific
			"system_prompt": schema.StringAttribute{
				Required:    true,
				Description: "The system prompt that is used to set the LLM context.",
			},
			"user_prompt": schema.StringAttribute{
				Required:    true,
				Description: "The user prompt that is used to set the LLM context.",
			},
			"llm_model_id": schema.StringAttribute{
				Required:    true,
				Description: "The ID of the LLM model to use for the Task.",
			},
			"image_required": schema.BoolAttribute{
				Optional:    true,
				Description: "",
			},
			"output_format": schema.MapAttribute{
				Required:    true,
				ElementType: types.StringType,
			},
		},
	}
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
	var data TaskResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	in := sdk.NewCreateTaskRequest()
	in.Name = data.Name.ValueString()
	in.Description = data.Description.ValueString()
	in.LLMModelID = data.LLMModelID.ValueString()
	in.SystemPrompt = data.SystemPrompt.ValueString()
	in.UserPrompt = data.UserPrompt.ValueString()

	for k, v := range data.OutputFormat {
		in.OutputFormat[k] = v.ValueString()
	}

	task, err := r.client.Create(ctx, in)
	if err != nil {
		resp.Diagnostics.AddError(err.Error(), "")
		return
	}

	if err := data.PopulateFromTaskModel(task); err != nil {
		resp.Diagnostics.AddError(err.Error(), "")
		return
	}

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *TaskResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data TaskResourceModel

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

	if err := data.PopulateFromTaskModel(task); err != nil {
		resp.Diagnostics.AddError(err.Error(), "")
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *TaskResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data TaskResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	in := sdk.NewUpdateTaskRequest(data.ID.ValueString())
	in.Name = data.Name.ValueString()
	in.Description = data.Description.ValueString()
	in.LLMModelID = data.LLMModelID.ValueString()
	in.SystemPrompt = data.SystemPrompt.ValueString()
	in.UserPrompt = data.UserPrompt.ValueString()

	for k, v := range data.OutputFormat {
		in.OutputFormat[k] = v.ValueString()
	}

	task, err := r.client.Update(ctx, in)
	if err != nil {
		resp.Diagnostics.AddError(err.Error(), "")
		return
	}

	if err := data.PopulateFromTaskModel(task); err != nil {
		resp.Diagnostics.AddError(err.Error(), "")
		return
	}

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *TaskResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data TaskResourceModel

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

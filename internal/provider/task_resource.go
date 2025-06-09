// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"fmt"

	"terraform-provider-tasks/internal/sdk"
	entitites "terraform-provider-tasks/internal/sdk/entities"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
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
	ActiveRevisionID types.String                             `tfsdk:"active_revision_id"`
	Description      types.String                             `tfsdk:"description"`
	Enabled          types.Bool                               `tfsdk:"enabled"`
	ExposedToAgents  types.Bool                               `tfsdk:"exposed_to_agents"`
	ID               types.String                             `tfsdk:"id"`
	ImageRequired    types.Bool                               `tfsdk:"image_required"`
	InputProcessors  *InputProcessorsModel                    `tfsdk:"input_processors"`
	LLMModelID       types.String                             `tfsdk:"llm_model_id"`
	Name             types.String                             `tfsdk:"name"`
	OptimiseImages   types.Bool                               `tfsdk:"optimise_images"`
	OutputFormat     map[string]entitites.OutputFormatWrapper `tfsdk:"output_format"`
	OutputModality   types.String                             `tfsdk:"output_modality"`
	Public           types.Bool                               `tfsdk:"public"`
	SystemPrompt     types.String                             `tfsdk:"system_prompt"`
	UserPrompt       types.String                             `tfsdk:"user_prompt"`
}

func (trm *TaskResourceModel) HasInputProcessors() bool {
	return trm.InputProcessors != nil && len(trm.InputProcessors.InputProcessors) > 0
}

func (trm *TaskResourceModel) PopulateFromTaskEntity(task *entitites.Task) error {

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
				Default:     booldefault.StaticBool(false),
				Computed:    true,
			},
			"exposed_to_agents": schema.BoolAttribute{
				Optional:    true,
				Description: "",
				Default:     booldefault.StaticBool(false),
				Computed:    true,
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
				Description: "When true it requires an image to be sent in the Task Run request.",
				Default:     booldefault.StaticBool(false),
				Computed:    true,
			},
			"optimise_images": schema.BoolAttribute{
				Optional:    true,
				Description: "When true (default) images will be automatically optimised before processing. Set to false to disable lossy image optimisation.",
				Default:     booldefault.StaticBool(true),
				Computed:    true,
			},
			"output_modality": schema.StringAttribute{
				Optional:    true,
				Description: "Specifies the output modality of the task. Can be 'json' or 'image'",
				Default:     stringdefault.StaticString("json"),
				Computed:    true,
				Validators: []validator.String{
					stringvalidator.OneOf("json", "image"),
				},
			},
			"output_format": schema.MapNestedAttribute{
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"type": schema.StringAttribute{
							Required: true,
						},
						"description": schema.StringAttribute{
							Optional: true,
						},
						"item_type": schema.StringAttribute{
							Optional: true,
						},
						"object": schema.MapAttribute{
							Optional:    true,
							ElementType: types.StringType,
						},
						"options": schema.MapAttribute{
							Optional:    true,
							ElementType: types.StringType,
						},
					},
				},
			},
		},
		Blocks: map[string]schema.Block{
			"input_processors": schema.SingleNestedBlock{
				Blocks: map[string]schema.Block{
					"input_processor": schema.ListNestedBlock{
						NestedObject: schema.NestedBlockObject{
							Attributes: map[string]schema.Attribute{
								"param_name": schema.StringAttribute{
									Required: true,
								},
								"input_processor": schema.StringAttribute{
									Required: true,
								},
								"config": schema.MapAttribute{
									Optional:    true,
									ElementType: types.StringType,
								},
							},
						},
					},
				},
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

	if err := data.PopulateFromTaskEntity(task); err != nil {
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

func (r *TaskResource) FormatInputProcessors(data TaskResourceModel) *[]entitites.InputProcessor {
	ips := []entitites.InputProcessor{}
	if !data.HasInputProcessors() {
		return &ips
	}
	for _, v := range data.InputProcessors.InputProcessors {
		ip := entitites.InputProcessor{
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

// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"fmt"

	"terraform-provider-tasks/internal/sdk"
	entities "terraform-provider-tasks/internal/sdk/entities"

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
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

const TASK_SCHEMA_VERSION = 0

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

type InputProcessorModelCollection struct {
	InputProcessors []InputProcessorModel `tfsdk:"input_processor"`
}

type InputProcessorModel struct {
	ParamName      types.String            `tfsdk:"param_name"`
	InputProcessor types.String            `tfsdk:"input_processor"`
	Config         map[string]types.String `tfsdk:"config"`
}

type OutputFormatModelCollection struct {
	OutputFormats []OutputFormatModel `tfsdk:"output_format"`
}

// OutputFormat represents the structure of the output format in a revision.
type OutputFormatModel struct {
	Description types.String            `tfsdk:"description"`
	Object      map[string]types.String `tfsdk:"object"`
	Options     map[string]types.String `tfsdk:"options"`
	Type        types.String            `tfsdk:"type"`
	Name        types.String            `tfsdk:"name"`
	ItemType    types.String            `tfsdk:"item_type"`
}

// TaskResourceModel describes the resource data model.
type TaskResourceModel struct {
	ID              types.String `tfsdk:"id"`
	Name            types.String `tfsdk:"name"`
	Enabled         types.Bool   `tfsdk:"enabled"`
	Public          types.Bool   `tfsdk:"public"`
	Description     types.String `tfsdk:"description"`
	ExposedToAgents types.Bool   `tfsdk:"exposed_to_agents"`

	SystemPrompt    types.String                   `tfsdk:"system_prompt"`
	UserPrompt      types.String                   `tfsdk:"user_prompt"`
	LLMModelID      types.String                   `tfsdk:"llm_model_id"`
	ImageRequired   types.Bool                     `tfsdk:"image_required"`
	OutputFormat    map[string]types.String        `tfsdk:"output_format"`
	OutputFormats   *OutputFormatModelCollection   `tfsdk:"output_formats"`
	OutputModality  types.String                   `tfsdk:"output_modality"`
	InputProcessors *InputProcessorModelCollection `tfsdk:"input_processors"`
	OptimiseImages  types.Bool                     `tfsdk:"optimise_images"`

	ActiveRevisionID types.String `tfsdk:"active_revision_id"`
}

func (trm *TaskResourceModel) HasInputProcessors() bool {
	return trm.InputProcessors != nil && len(trm.InputProcessors.InputProcessors) > 0
}

func (trm *TaskResourceModel) HasOutputFormats() bool {
	return trm.OutputFormats != nil && len(trm.OutputFormats.OutputFormats) > 0
}

func (trm *TaskResourceModel) PopulateFromTaskEntity(task *entities.Task) error {

	rev, err := task.GetActiveRevision()
	if err != nil {
		return err
	}

	trm.ID = types.StringValue(task.ID)
	trm.Name = types.StringValue(task.Name)
	trm.Enabled = types.BoolValue(task.Enabled)
	trm.Public = types.BoolValue(task.Public)
	trm.Description = types.StringValue(task.Description)

	trm.ExposedToAgents = types.BoolValue(task.ExposedToAgents)
	trm.ImageRequired = types.BoolValue(rev.ImageRequired)
	trm.LLMModelID = types.StringValue(rev.LLMModelID)
	trm.OptimiseImages = types.BoolValue(rev.OptimiseImages)
	trm.OutputModality = types.StringValue(rev.OutputModality)
	trm.SystemPrompt = types.StringValue(rev.SystemPrompt)
	trm.UserPrompt = types.StringValue(rev.UserPrompt)
	trm.ActiveRevisionID = types.StringValue(rev.ID)

	if rev.HasInputProcessors() {
		trm.InputProcessors = &InputProcessorModelCollection{
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

	trm.OutputFormat = make(map[string]types.String)
	trm.OutputFormats = &OutputFormatModelCollection{
		OutputFormats: make([]OutputFormatModel, 0),
	}

	for k, v := range rev.OutputFormat {
		if v.IsSimple() {
			trm.OutputFormat[k] = types.StringValue(v.Simple.String())
			trm.OutputFormats.OutputFormats = append(trm.OutputFormats.OutputFormats, OutputFormatModel{
				Name: types.StringValue(k),
				Type: types.StringValue(v.Simple.String()),
			})
		}
		if v.IsExtended() {
			of := OutputFormatModel{
				Type:        types.StringValue(v.Extended.Type),
				Description: types.StringValue(v.Extended.Description),
				Object:      make(map[string]types.String, len(v.Extended.Object)),
				Options:     make(map[string]types.String, len(v.Extended.Options)),
				ItemType:    types.StringValue(v.Extended.ItemType),
			}
			for k, v := range v.Extended.Object {
				of.Object[k] = types.StringValue(v)
			}
			for k, v := range v.Extended.Options {
				of.Options[k] = types.StringValue(v)
			}
			trm.OutputFormats.OutputFormats = append(trm.OutputFormats.OutputFormats, of)
		}
	}

	return nil
}

func (r *TaskResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_task"
}

func (r *TaskResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{

		MarkdownDescription: "Task resource",
		Version:             TASK_SCHEMA_VERSION,

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
			"active_revision_id": schema.StringAttribute{
				Computed: true,
			},
			"exposed_to_agents": schema.BoolAttribute{
				Optional:    true,
				Description: "",
				Default:     booldefault.StaticBool(false),
				Computed:    true,
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
				Default:     booldefault.StaticBool(false),
				Computed:    true,
			},
			"output_format": schema.MapAttribute{
				Required:    true,
				ElementType: types.StringType,
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
		},
		Blocks: map[string]schema.Block{
			"output_formats": schema.SingleNestedBlock{
				Blocks: map[string]schema.Block{
					"output_format": schema.ListNestedBlock{
						NestedObject: schema.NestedBlockObject{
							Attributes: map[string]schema.Attribute{
								"name": schema.StringAttribute{
									Required: true,
								},
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
			},
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
		tflog.Error(ctx, "Error creating task resource")
		return
	}

	in := sdk.NewCreateTaskRequest()
	in.Name = data.Name.ValueString()
	in.Description = data.Description.ValueString()
	in.LLMModelID = data.LLMModelID.ValueString()
	in.SystemPrompt = data.SystemPrompt.ValueString()
	in.UserPrompt = data.UserPrompt.ValueString()
	in.Enabled = data.Enabled.ValueBool()
	in.Public = data.Public.ValueBool()
	in.ImageRequired = data.ImageRequired.ValueBool()
	in.OutputModality = data.OutputModality.ValueString()

	in.InputProcessors = r.FormatInputProcessors(data)
	in.OutputFormat = r.FormatOutputFormat(data)

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
		tflog.Error(ctx, "Error reading task resource")
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
		tflog.Error(ctx, "Error updating task resource")
		return
	}

	in := sdk.NewUpdateTaskRequest(data.ID.ValueString())
	in.Name = data.Name.ValueString()
	in.Description = data.Description.ValueString()
	in.LLMModelID = data.LLMModelID.ValueString()
	in.SystemPrompt = data.SystemPrompt.ValueString()
	in.UserPrompt = data.UserPrompt.ValueString()
	in.Enabled = data.Enabled.ValueBool()
	in.Public = data.Public.ValueBool()
	in.ImageRequired = data.ImageRequired.ValueBool()
	in.OutputModality = data.OutputModality.ValueString()

	in.InputProcessors = r.FormatInputProcessors(data)
	in.OutputFormat = r.FormatOutputFormat(data)

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
		tflog.Error(ctx, "Error deleting task resource")
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

// FormatInputProcessors converts the TaskResourceModel's input processors data into the
// format expected by the SDK.
func (r *TaskResource) FormatInputProcessors(data TaskResourceModel) *[]entities.InputProcessor {
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

// FormatOutputFormat converts the TaskResourceModel's output format data into the
// format expected by the SDK.
//
// It handles both simple and extended output formats, ensuring that the
// output format is correctly structured for the Task entity.
func (r *TaskResource) FormatOutputFormat(data TaskResourceModel) map[string]entities.OutputFormatExtended {
	ofw := make(map[string]entities.OutputFormatExtended)
	for k, v := range data.OutputFormat {
		if v.ValueString() == "" {
			ofw[k] = entities.OutputFormatExtended{
				Type: v.ValueString(),
			}
		}
	}
	if data.HasOutputFormats() {
		for _, v := range data.OutputFormats.OutputFormats {
			ofw[v.Name.ValueString()] = entities.OutputFormatExtended{
				Description: v.Description.ValueString(),
				Object:      make(map[string]string, len(v.Object)),
				Options:     make(map[string]string, len(v.Options)),
				Type:        v.Type.ValueString(),
				ItemType:    v.ItemType.ValueString(),
			}
			for k, obv := range v.Object {
				ofw[v.Name.ValueString()].Object[k] = obv.ValueString()
			}
			for k, obv := range v.Options {
				ofw[v.Name.ValueString()].Options[k] = obv.ValueString()
			}
		}
	}
	return ofw
}

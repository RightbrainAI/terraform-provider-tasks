// Copyright (c) HashiCorp, Inc.

package v0

import (
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var TASK_SCHEMA = schema.Schema{

	MarkdownDescription: "Task resource",
	Version:             1,

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
			Required: true,
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

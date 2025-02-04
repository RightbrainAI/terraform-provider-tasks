// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"fmt"
	"terraform-provider-tasks/internal/sdk"
	entitites "terraform-provider-tasks/internal/sdk/entities"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ datasource.DataSource = &LLMModelDataSource{}

func NewLLMModelDataSource() datasource.DataSource {
	return &LLMModelDataSource{}
}

// LLMModelDataSource defines the data source implementation.
type LLMModelDataSource struct {
	client *sdk.TasksClient
}

// LLMModelDataSourceModel describes the data source data model.
type LLMModelDataSourceModel struct {
	ID types.String `tfsdk:"id"`

	Alias          types.String `tfsdk:"alias"`
	Description    types.String `tfsdk:"description"`
	Name           types.String `tfsdk:"name"`
	Provider       types.String `tfsdk:"model_provider"`
	SupportsVision types.Bool   `tfsdk:"supports_vision"`
}

func (d *LLMModelDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_model"
}

func (d *LLMModelDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "LLM Model data source",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "LLMModel identifier",
				Computed:            true,
			},
			"alias": schema.StringAttribute{
				Optional:    true,
				Description: "",
			},
			"description": schema.StringAttribute{
				Optional:    true,
				Description: "",
			},
			"name": schema.StringAttribute{
				Required:    true,
				Description: "",
			},
			"model_provider": schema.StringAttribute{
				Optional:    true,
				Description: "",
			},
			"supports_vision": schema.BoolAttribute{
				Optional:    true,
				Description: "",
			},
		},
	}
}

func (d *LLMModelDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

	d.client = client
}

func (d *LLMModelDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data LLMModelDataSourceModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	models, err := d.client.GetAvailableLLMModels(ctx)
	if err != nil {
		resp.Diagnostics.AddError(err.Error(), "")
		return
	}

	var found *entitites.Model

	for _, model := range models {
		if model.Name == data.Name.ValueString() {
			found = &model
		}
	}

	if found == nil {
		resp.Diagnostics.AddError(fmt.Sprintf("cannot find model %s", data.Name), "")
		return
	}

	data.ID = types.StringValue(found.ID)
	data.Alias = types.StringValue(found.Alias)
	data.Description = types.StringValue(found.Description)
	data.Name = types.StringValue(found.Name)
	data.Provider = types.StringValue(found.Provider)
	data.SupportsVision = types.BoolValue(found.SupportsVision)

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

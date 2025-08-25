package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ datasource.DataSource = &ScenarioDataSource{}

func NewScenarioDataSource() datasource.DataSource {
	return &ScenarioDataSource{}
}

// ScenarioDataSource defines the data source implementation.
type ScenarioDataSource struct {
	client *MakeAPIClient
}

// ScenarioDataSourceModel describes the data source data model.
type ScenarioDataSourceModel struct {
	Id          types.String `tfsdk:"id"`
	Name        types.String `tfsdk:"name"`
	Description types.String `tfsdk:"description"`
	Active      types.Bool   `tfsdk:"active"`
	TeamId      types.String `tfsdk:"team_id"`
}

func (d *ScenarioDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_scenario"
}

func (d *ScenarioDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "Make.com scenario data source",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "Scenario identifier",
				Required:            true,
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "Name of the scenario",
				Computed:            true,
			},
			"description": schema.StringAttribute{
				MarkdownDescription: "Description of the scenario",
				Computed:            true,
			},
			"active": schema.BoolAttribute{
				MarkdownDescription: "Whether the scenario is active",
				Computed:            true,
			},
			"team_id": schema.StringAttribute{
				MarkdownDescription: "Team ID where the scenario belongs",
				Computed:            true,
			},
		},
	}
}

func (d *ScenarioDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*MakeAPIClient)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *MakeAPIClient, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	d.client = client
}

func (d *ScenarioDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data ScenarioDataSourceModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Get the scenario from the API
	scenario, err := d.client.GetScenario(ctx, data.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read scenario, got error: %s", err))
		return
	}

	// Map API response to Terraform state
	data.Id = types.StringValue(scenario.ID)
	data.Name = types.StringValue(scenario.Name)
	data.Active = types.BoolValue(scenario.Active)

	if scenario.Description != "" {
		data.Description = types.StringValue(scenario.Description)
	} else {
		data.Description = types.StringNull()
	}

	if scenario.TeamID != "" {
		data.TeamId = types.StringValue(scenario.TeamID)
	} else {
		data.TeamId = types.StringNull()
	}

	// Write logs using the tflog package
	tflog.Trace(ctx, "read a scenario data source")

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

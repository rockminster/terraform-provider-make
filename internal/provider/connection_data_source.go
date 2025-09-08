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
var _ datasource.DataSource = &ConnectionDataSource{}

func NewConnectionDataSource() datasource.DataSource {
	return &ConnectionDataSource{}
}

// ConnectionDataSource defines the data source implementation.
type ConnectionDataSource struct {
	client *MakeAPIClient
}

// ConnectionDataSourceModel describes the data source data model.
type ConnectionDataSourceModel struct {
	Id       types.String `tfsdk:"id"`
	Name     types.String `tfsdk:"name"`
	AppName  types.String `tfsdk:"app_name"`
	TeamId   types.String `tfsdk:"team_id"`
	Verified types.Bool   `tfsdk:"verified"`
	Settings types.Map    `tfsdk:"settings"`
}

func (d *ConnectionDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_connection"
}

func (d *ConnectionDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "Make.com connection data source",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "Connection identifier",
				Required:            true,
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "Name of the connection",
				Computed:            true,
			},
			"app_name": schema.StringAttribute{
				MarkdownDescription: "Name of the app for this connection",
				Computed:            true,
			},
			"team_id": schema.StringAttribute{
				MarkdownDescription: "Team ID where the connection belongs",
				Computed:            true,
			},
			"verified": schema.BoolAttribute{
				MarkdownDescription: "Whether the connection is verified",
				Computed:            true,
			},
			"settings": schema.MapAttribute{
				MarkdownDescription: "Advanced settings for the connection",
				Computed:            true,
				ElementType:         types.StringType,
			},
		},
	}
}

func (d *ConnectionDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *ConnectionDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data ConnectionDataSourceModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Get the connection from the API
	connection, err := d.client.GetConnection(ctx, data.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read connection, got error: %s", err))
		return
	}

	// Map API response to Terraform state
	data.Id = types.StringValue(connection.ID)
	data.Name = types.StringValue(connection.Name)
	data.AppName = types.StringValue(connection.AppName)
	data.Verified = types.BoolValue(connection.Verified)

	if connection.TeamID != "" {
		data.TeamId = types.StringValue(connection.TeamID)
	} else {
		data.TeamId = types.StringNull()
	}

	if len(connection.Settings) > 0 {
		data.Settings = types.MapValueMust(types.StringType, convertSettingsToStringMap(connection.Settings))
	} else {
		data.Settings = types.MapNull(types.StringType)
	}

	// Write logs using the tflog package
	tflog.Trace(ctx, "read a connection data source")

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

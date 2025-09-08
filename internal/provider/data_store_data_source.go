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
var _ datasource.DataSource = &DataStoreDataSource{}

func NewDataStoreDataSource() datasource.DataSource {
	return &DataStoreDataSource{}
}

// DataStoreDataSource defines the data source implementation.
type DataStoreDataSource struct {
	client *MakeAPIClient
}

// DataStoreDataSourceModel describes the data source data model.
type DataStoreDataSourceModel struct {
	Id          types.String `tfsdk:"id"`
	Name        types.String `tfsdk:"name"`
	Description types.String `tfsdk:"description"`
	TeamId      types.String `tfsdk:"team_id"`
}

func (d *DataStoreDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_data_store"
}

func (d *DataStoreDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Make.com data store data source",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "Data store identifier",
				Required:            true,
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "Name of the data store",
				Computed:            true,
			},
			"description": schema.StringAttribute{
				MarkdownDescription: "Description of the data store",
				Computed:            true,
			},
			"team_id": schema.StringAttribute{
				MarkdownDescription: "Team ID where the data store belongs",
				Computed:            true,
			},
		},
	}
}

func (d *DataStoreDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *DataStoreDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data DataStoreDataSourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	ds, err := d.client.GetDataStore(ctx, data.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read data store, got error: %s", err))
		return
	}

	data.Id = types.StringValue(ds.ID)
	data.Name = types.StringValue(ds.Name)
	if ds.Description == "" {
		data.Description = types.StringNull()
	} else {
		data.Description = types.StringValue(ds.Description)
	}
	if ds.TeamID == "" {
		data.TeamId = types.StringNull()
	} else {
		data.TeamId = types.StringValue(ds.TeamID)
	}

	tflog.Trace(ctx, "read a data store data source")

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

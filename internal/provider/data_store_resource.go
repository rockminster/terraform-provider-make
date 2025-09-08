package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.Resource = &DataStoreResource{}
var _ resource.ResourceWithImportState = &DataStoreResource{}

func NewDataStoreResource() resource.Resource {
	return &DataStoreResource{}
}

// DataStoreResource defines the resource implementation.
type DataStoreResource struct {
	client *MakeAPIClient
}

// DataStoreResourceModel describes the resource data model.
type DataStoreResourceModel struct {
	Id          types.String `tfsdk:"id"`
	Name        types.String `tfsdk:"name"`
	Description types.String `tfsdk:"description"`
	TeamId      types.String `tfsdk:"team_id"`
}

func (r *DataStoreResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_data_store"
}

func (r *DataStoreResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Make.com data store resource",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Data store identifier",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "Name of the data store",
				Required:            true,
			},
			"description": schema.StringAttribute{
				MarkdownDescription: "Description of the data store",
				Optional:            true,
			},
			"team_id": schema.StringAttribute{
				MarkdownDescription: "Team ID where the data store belongs",
				Optional:            true,
			},
		},
	}
}

func (r *DataStoreResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*MakeAPIClient)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *MakeAPIClient, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	r.client = client
}

func (r *DataStoreResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data DataStoreResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	apiReq := DataStoreRequest{
		Name: data.Name.ValueString(),
	}

	if !data.Description.IsNull() {
		apiReq.Description = data.Description.ValueString()
	}

	if !data.TeamId.IsNull() {
		apiReq.TeamID = data.TeamId.ValueString()
	}

	ds, err := r.client.CreateDataStore(ctx, apiReq)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create data store, got error: %s", err))
		return
	}

	data.Id = types.StringValue(ds.ID)
	data.Name = types.StringValue(ds.Name)

	if ds.Description != "" {
		data.Description = types.StringValue(ds.Description)
	}

	if ds.TeamID != "" {
		data.TeamId = types.StringValue(ds.TeamID)
	}

	tflog.Trace(ctx, "created a data store resource")

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *DataStoreResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data DataStoreResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	ds, err := r.client.GetDataStore(ctx, data.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read data store, got error: %s", err))
		return
	}

	data.Id = types.StringValue(ds.ID)
	data.Name = types.StringValue(ds.Name)

	if ds.Description != "" {
		data.Description = types.StringValue(ds.Description)
	} else {
		data.Description = types.StringNull()
	}

	if ds.TeamID != "" {
		data.TeamId = types.StringValue(ds.TeamID)
	} else {
		data.TeamId = types.StringNull()
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *DataStoreResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data DataStoreResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	apiReq := DataStoreRequest{
		Name: data.Name.ValueString(),
	}

	if !data.Description.IsNull() {
		apiReq.Description = data.Description.ValueString()
	}

	if !data.TeamId.IsNull() {
		apiReq.TeamID = data.TeamId.ValueString()
	}

	ds, err := r.client.UpdateDataStore(ctx, data.Id.ValueString(), apiReq)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to update data store, got error: %s", err))
		return
	}

	data.Id = types.StringValue(ds.ID)
	data.Name = types.StringValue(ds.Name)

	if ds.Description != "" {
		data.Description = types.StringValue(ds.Description)
	} else {
		data.Description = types.StringNull()
	}

	if ds.TeamID != "" {
		data.TeamId = types.StringValue(ds.TeamID)
	} else {
		data.TeamId = types.StringNull()
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *DataStoreResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data DataStoreResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.DeleteDataStore(ctx, data.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete data store, got error: %s", err))
		return
	}
}

func (r *DataStoreResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), req.ID)...)
}

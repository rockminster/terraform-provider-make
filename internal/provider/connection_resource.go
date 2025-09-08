package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.Resource = &ConnectionResource{}
var _ resource.ResourceWithImportState = &ConnectionResource{}

func NewConnectionResource() resource.Resource {
	return &ConnectionResource{}
}

// ConnectionResource defines the resource implementation.
type ConnectionResource struct {
	client *MakeAPIClient
}

// ConnectionResourceModel describes the resource data model.
type ConnectionResourceModel struct {
	Id       types.String `tfsdk:"id"`
	Name     types.String `tfsdk:"name"`
	AppName  types.String `tfsdk:"app_name"`
	TeamId   types.String `tfsdk:"team_id"`
	Settings types.Map    `tfsdk:"settings"`
	Verified types.Bool   `tfsdk:"verified"`
}

func (r *ConnectionResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_connection"
}

func (r *ConnectionResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "Make.com connection resource",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Connection identifier",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "Name of the connection",
				Required:            true,
			},
			"app_name": schema.StringAttribute{
				MarkdownDescription: "Name of the app for this connection (e.g., 'gmail', 'slack')",
				Required:            true,
			},
			"team_id": schema.StringAttribute{
				MarkdownDescription: "Team ID where the connection belongs",
				Optional:            true,
			},
			"settings": schema.MapAttribute{
				MarkdownDescription: "Advanced settings for the connection",
				Optional:            true,
				ElementType:         types.StringType,
			},
			"verified": schema.BoolAttribute{
				MarkdownDescription: "Whether the connection is verified",
				Computed:            true,
			},
		},
	}
}

func (r *ConnectionResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *ConnectionResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data ConnectionResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Prepare the API request
	apiReq := ConnectionRequest{
		Name:    data.Name.ValueString(),
		AppName: data.AppName.ValueString(),
	}

	if !data.TeamId.IsNull() {
		apiReq.TeamID = data.TeamId.ValueString()
	}

	if !data.Settings.IsNull() {
		var settingsMap map[string]string
		resp.Diagnostics.Append(data.Settings.ElementsAs(ctx, &settingsMap, false)...)
		if resp.Diagnostics.HasError() {
			return
		}
		apiReq.Settings = make(map[string]interface{}, len(settingsMap))
		for k, v := range settingsMap {
			apiReq.Settings[k] = v
		}
	}

	// Create the connection via API
	connection, err := r.client.CreateConnection(ctx, apiReq)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create connection, got error: %s", err))
		return
	}

	// Map response to Terraform state
	data.Id = types.StringValue(connection.ID)
	data.Name = types.StringValue(connection.Name)
	data.AppName = types.StringValue(connection.AppName)
	data.Verified = types.BoolValue(connection.Verified)

	if connection.TeamID != "" {
		data.TeamId = types.StringValue(connection.TeamID)
	}

	if len(connection.Settings) > 0 {
		settingsVals := make(map[string]attr.Value, len(connection.Settings))
		for k, v := range connection.Settings {
			settingsVals[k] = types.StringValue(fmt.Sprintf("%v", v))
		}
		data.Settings = types.MapValueMust(types.StringType, settingsVals)
	}

	// Write logs using the tflog package
	tflog.Trace(ctx, "created a connection resource")

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *ConnectionResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data ConnectionResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Get the connection from the API
	connection, err := r.client.GetConnection(ctx, data.Id.ValueString())
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
		settingsVals := make(map[string]attr.Value, len(connection.Settings))
		for k, v := range connection.Settings {
			settingsVals[k] = types.StringValue(fmt.Sprintf("%v", v))
		}
		data.Settings = types.MapValueMust(types.StringType, settingsVals)
	} else {
		data.Settings = types.MapNull(types.StringType)
	}

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *ConnectionResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data ConnectionResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Prepare the API request
	apiReq := ConnectionRequest{
		Name:    data.Name.ValueString(),
		AppName: data.AppName.ValueString(),
	}

	if !data.TeamId.IsNull() {
		apiReq.TeamID = data.TeamId.ValueString()
	}

	if !data.Settings.IsNull() {
		var settingsMap map[string]string
		resp.Diagnostics.Append(data.Settings.ElementsAs(ctx, &settingsMap, false)...)
		if resp.Diagnostics.HasError() {
			return
		}
		apiReq.Settings = make(map[string]interface{}, len(settingsMap))
		for k, v := range settingsMap {
			apiReq.Settings[k] = v
		}
	}

	// Update the connection via API
	connection, err := r.client.UpdateConnection(ctx, data.Id.ValueString(), apiReq)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to update connection, got error: %s", err))
		return
	}

	// Map response to Terraform state
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
		settingsVals := make(map[string]attr.Value, len(connection.Settings))
		for k, v := range connection.Settings {
			settingsVals[k] = types.StringValue(fmt.Sprintf("%v", v))
		}
		data.Settings = types.MapValueMust(types.StringType, settingsVals)
	} else {
		data.Settings = types.MapNull(types.StringType)
	}

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *ConnectionResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data ConnectionResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Delete the connection via API
	err := r.client.DeleteConnection(ctx, data.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete connection, got error: %s", err))
		return
	}
}

func (r *ConnectionResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to id attribute
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), req.ID)...)
}

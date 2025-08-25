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
var _ resource.Resource = &ScenarioResource{}
var _ resource.ResourceWithImportState = &ScenarioResource{}

func NewScenarioResource() resource.Resource {
	return &ScenarioResource{}
}

// ScenarioResource defines the resource implementation.
type ScenarioResource struct {
	client *MakeAPIClient
}

// ScenarioResourceModel describes the resource data model.
type ScenarioResourceModel struct {
	Id          types.String `tfsdk:"id"`
	Name        types.String `tfsdk:"name"`
	Description types.String `tfsdk:"description"`
	Active      types.Bool   `tfsdk:"active"`
	TeamId      types.String `tfsdk:"team_id"`
}

func (r *ScenarioResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_scenario"
}

func (r *ScenarioResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "Make.com scenario resource",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Scenario identifier",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "Name of the scenario",
				Required:            true,
			},
			"description": schema.StringAttribute{
				MarkdownDescription: "Description of the scenario",
				Optional:            true,
			},
			"active": schema.BoolAttribute{
				MarkdownDescription: "Whether the scenario is active",
				Optional:            true,
			},
			"team_id": schema.StringAttribute{
				MarkdownDescription: "Team ID where the scenario belongs",
				Optional:            true,
			},
		},
	}
}

func (r *ScenarioResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *ScenarioResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data ScenarioResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Prepare the API request
	apiReq := ScenarioRequest{
		Name:   data.Name.ValueString(),
		Active: data.Active.ValueBool(),
	}

	if !data.Description.IsNull() {
		apiReq.Description = data.Description.ValueString()
	}

	if !data.TeamId.IsNull() {
		apiReq.TeamID = data.TeamId.ValueString()
	}

	// Create the scenario via API
	scenario, err := r.client.CreateScenario(ctx, apiReq)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create scenario, got error: %s", err))
		return
	}

	// Map response to Terraform state
	data.Id = types.StringValue(scenario.ID)
	data.Name = types.StringValue(scenario.Name)
	data.Active = types.BoolValue(scenario.Active)

	if scenario.Description != "" {
		data.Description = types.StringValue(scenario.Description)
	}

	if scenario.TeamID != "" {
		data.TeamId = types.StringValue(scenario.TeamID)
	}

	// Write logs using the tflog package
	tflog.Trace(ctx, "created a scenario resource")

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *ScenarioResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data ScenarioResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Get the scenario from the API
	scenario, err := r.client.GetScenario(ctx, data.Id.ValueString())
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

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *ScenarioResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data ScenarioResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Prepare the API request
	apiReq := ScenarioRequest{
		Name:   data.Name.ValueString(),
		Active: data.Active.ValueBool(),
	}

	if !data.Description.IsNull() {
		apiReq.Description = data.Description.ValueString()
	}

	if !data.TeamId.IsNull() {
		apiReq.TeamID = data.TeamId.ValueString()
	}

	// Update the scenario via API
	scenario, err := r.client.UpdateScenario(ctx, data.Id.ValueString(), apiReq)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to update scenario, got error: %s", err))
		return
	}

	// Map response to Terraform state
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

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *ScenarioResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data ScenarioResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Delete the scenario via API
	err := r.client.DeleteScenario(ctx, data.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete scenario, got error: %s", err))
		return
	}
}

func (r *ScenarioResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to id attribute
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), req.ID)...)
}

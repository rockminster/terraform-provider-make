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
var _ resource.Resource = &WebhookResource{}
var _ resource.ResourceWithImportState = &WebhookResource{}

func NewWebhookResource() resource.Resource {
	return &WebhookResource{}
}

// WebhookResource defines the resource implementation.
type WebhookResource struct {
	client *MakeAPIClient
}

// WebhookResourceModel describes the resource data model.
type WebhookResourceModel struct {
	Id     types.String `tfsdk:"id"`
	Name   types.String `tfsdk:"name"`
	URL    types.String `tfsdk:"url"`
	TeamId types.String `tfsdk:"team_id"`
	Active types.Bool   `tfsdk:"active"`
}

func (r *WebhookResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_webhook"
}

func (r *WebhookResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "Make.com webhook resource",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Webhook identifier",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "Name of the webhook",
				Required:            true,
			},
			"url": schema.StringAttribute{
				MarkdownDescription: "URL endpoint for the webhook",
				Computed:            true,
			},
			"team_id": schema.StringAttribute{
				MarkdownDescription: "Team ID where the webhook belongs",
				Optional:            true,
			},
			"active": schema.BoolAttribute{
				MarkdownDescription: "Whether the webhook is active",
				Optional:            true,
			},
		},
	}
}

func (r *WebhookResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *WebhookResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data WebhookResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Prepare the API request
	apiReq := WebhookRequest{
		Name:   data.Name.ValueString(),
		Active: data.Active.ValueBool(),
	}

	if !data.TeamId.IsNull() {
		apiReq.TeamID = data.TeamId.ValueString()
	}

	// Create the webhook via API
	webhook, err := r.client.CreateWebhook(ctx, apiReq)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create webhook, got error: %s", err))
		return
	}

	// Map response to Terraform state
	data.Id = types.StringValue(webhook.ID)
	data.Name = types.StringValue(webhook.Name)
	data.URL = types.StringValue(webhook.URL)
	data.Active = types.BoolValue(webhook.Active)

	if webhook.TeamID != "" {
		data.TeamId = types.StringValue(webhook.TeamID)
	}

	// Write logs using the tflog package
	tflog.Trace(ctx, "created a webhook resource")

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *WebhookResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data WebhookResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Get the webhook from the API
	webhook, err := r.client.GetWebhook(ctx, data.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read webhook, got error: %s", err))
		return
	}

	// Map API response to Terraform state
	data.Id = types.StringValue(webhook.ID)
	data.Name = types.StringValue(webhook.Name)
	data.URL = types.StringValue(webhook.URL)
	data.Active = types.BoolValue(webhook.Active)

	if webhook.TeamID != "" {
		data.TeamId = types.StringValue(webhook.TeamID)
	} else {
		data.TeamId = types.StringNull()
	}

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *WebhookResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data WebhookResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Prepare the API request
	apiReq := WebhookRequest{
		Name:   data.Name.ValueString(),
		Active: data.Active.ValueBool(),
	}

	if !data.TeamId.IsNull() {
		apiReq.TeamID = data.TeamId.ValueString()
	}

	// Update the webhook via API
	webhook, err := r.client.UpdateWebhook(ctx, data.Id.ValueString(), apiReq)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to update webhook, got error: %s", err))
		return
	}

	// Map response to Terraform state
	data.Id = types.StringValue(webhook.ID)
	data.Name = types.StringValue(webhook.Name)
	data.URL = types.StringValue(webhook.URL)
	data.Active = types.BoolValue(webhook.Active)

	if webhook.TeamID != "" {
		data.TeamId = types.StringValue(webhook.TeamID)
	} else {
		data.TeamId = types.StringNull()
	}

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *WebhookResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data WebhookResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Delete the webhook via API
	err := r.client.DeleteWebhook(ctx, data.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete webhook, got error: %s", err))
		return
	}
}

func (r *WebhookResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to id attribute
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), req.ID)...)
}

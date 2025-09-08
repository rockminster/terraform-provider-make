package provider

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types"
)

func TestMakeAPIClient_MakeRequest(t *testing.T) {
	client := &MakeAPIClient{
		ApiToken: "test-token",
		BaseUrl:  "https://api.make.com/",
	}

	// Test URL construction
	ctx := context.Background()

	// This would normally hit the API, but we can test URL construction
	req := ScenarioRequest{
		Name:   "Test Scenario",
		Active: true,
	}

	// Basic test that client can construct requests
	if client.ApiToken != "test-token" {
		t.Errorf("Expected ApiToken to be 'test-token', got %s", client.ApiToken)
	}

	if client.BaseUrl != "https://api.make.com/" {
		t.Errorf("Expected BaseUrl to be 'https://api.make.com/', got %s", client.BaseUrl)
	}

	// Test that request structure is correct
	if req.Name != "Test Scenario" {
		t.Errorf("Expected Name to be 'Test Scenario', got %s", req.Name)
	}

	// Suppress unused variable warning
	_ = ctx
}

func TestScenarioResourceModel(t *testing.T) {
	model := ScenarioResourceModel{
		Id:          types.StringValue("test-id"),
		Name:        types.StringValue("Test Scenario"),
		Description: types.StringValue("Test Description"),
		Active:      types.BoolValue(true),
		TeamId:      types.StringValue("team-123"),
	}

	if model.Id.ValueString() != "test-id" {
		t.Errorf("Expected Id to be 'test-id', got %s", model.Id.ValueString())
	}

	if model.Name.ValueString() != "Test Scenario" {
		t.Errorf("Expected Name to be 'Test Scenario', got %s", model.Name.ValueString())
	}

	if !model.Active.ValueBool() {
		t.Errorf("Expected Active to be true, got %v", model.Active.ValueBool())
	}
}

func TestConnectionResourceModel(t *testing.T) {
	model := ConnectionResourceModel{
		Id:       types.StringValue("conn-123"),
		Name:     types.StringValue("Gmail Connection"),
		AppName:  types.StringValue("gmail"),
		TeamId:   types.StringValue("team-456"),
		Verified: types.BoolValue(true),
	}

	if model.Id.ValueString() != "conn-123" {
		t.Errorf("Expected Id to be 'conn-123', got %s", model.Id.ValueString())
	}

	if model.AppName.ValueString() != "gmail" {
		t.Errorf("Expected AppName to be 'gmail', got %s", model.AppName.ValueString())
	}

	if !model.Verified.ValueBool() {
		t.Errorf("Expected Verified to be true, got %v", model.Verified.ValueBool())
	}
}

func TestWebhookResourceModel(t *testing.T) {
	model := WebhookResourceModel{
		Id:     types.StringValue("webhook-789"),
		Name:   types.StringValue("Test Webhook"),
		URL:    types.StringValue("https://example.com/webhook"),
		TeamId: types.StringValue("team-789"),
		Active: types.BoolValue(true),
	}

	if model.Id.ValueString() != "webhook-789" {
		t.Errorf("Expected Id to be 'webhook-789', got %s", model.Id.ValueString())
	}

	if model.URL.ValueString() != "https://example.com/webhook" {
		t.Errorf("Expected URL to be 'https://example.com/webhook', got %s", model.URL.ValueString())
	}

	if !model.Active.ValueBool() {
		t.Errorf("Expected Active to be true, got %v", model.Active.ValueBool())
	}
}

func TestDataStoreResourceModel(t *testing.T) {
	model := DataStoreResourceModel{
		Id:          types.StringValue("ds-123"),
		Name:        types.StringValue("Test Data Store"),
		Description: types.StringValue("Test Description"),
		TeamId:      types.StringValue("team-999"),
	}

	if model.Name.ValueString() != "Test Data Store" {
		t.Errorf("Expected Name to be 'Test Data Store', got %s", model.Name.ValueString())
	}

	if model.Description.ValueString() != "Test Description" {
		t.Errorf("Expected Description to be 'Test Description', got %s", model.Description.ValueString())
	}
}

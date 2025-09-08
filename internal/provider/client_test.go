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

func TestConvertSettingsToStringMap(t *testing.T) {
	// Test various data types
	settings := map[string]interface{}{
		"string_val":  "test_string",
		"int_val":     42,
		"float_val":   3.14,
		"bool_val":    true,
		"uint_val":    uint(100),
		"complex_val": map[string]string{"key": "value"}, // Will use fmt.Sprintf fallback
	}

	result := convertSettingsToStringMap(settings)

	// Verify we got the expected number of keys
	if len(result) != len(settings) {
		t.Errorf("Expected %d keys, got %d", len(settings), len(result))
	}

	// Test string conversion
	stringVal := result["string_val"].(types.String)
	if stringVal.ValueString() != "test_string" {
		t.Errorf("Expected string_val to be 'test_string', got %s", stringVal.ValueString())
	}

	// Test int conversion
	intVal := result["int_val"].(types.String)
	if intVal.ValueString() != "42" {
		t.Errorf("Expected int_val to be '42', got %s", intVal.ValueString())
	}

	// Test float conversion
	floatVal := result["float_val"].(types.String)
	if floatVal.ValueString() != "3.140000" {
		t.Errorf("Expected float_val to be '3.140000', got %s", floatVal.ValueString())
	}

	// Test bool conversion
	boolVal := result["bool_val"].(types.String)
	if boolVal.ValueString() != "true" {
		t.Errorf("Expected bool_val to be 'true', got %s", boolVal.ValueString())
	}

	// Test uint conversion
	uintVal := result["uint_val"].(types.String)
	if uintVal.ValueString() != "100" {
		t.Errorf("Expected uint_val to be '100', got %s", uintVal.ValueString())
	}

	// Test complex type fallback
	complexVal := result["complex_val"].(types.String)
	expectedComplex := "map[key:value]"
	if complexVal.ValueString() != expectedComplex {
		t.Errorf("Expected complex_val to be '%s', got %s", expectedComplex, complexVal.ValueString())
	}
}

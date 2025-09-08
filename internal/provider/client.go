package provider

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"path"
)

// ScenarioResponse represents a Make.com scenario from the API
type ScenarioResponse struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	Active      bool   `json:"is_active"`
	TeamID      string `json:"team_id,omitempty"`
}

// ScenarioRequest represents the request payload for creating/updating scenarios
type ScenarioRequest struct {
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	Active      bool   `json:"is_active"`
	TeamID      string `json:"team_id,omitempty"`
}

// ErrorResponse represents an error response from Make.com API
type ErrorResponse struct {
	Error   string `json:"error,omitempty"`
	Message string `json:"message,omitempty"`
	Code    int    `json:"code,omitempty"`
}

// MakeRequest performs a HTTP request to the Make.com API
func (c *MakeAPIClient) MakeRequest(ctx context.Context, method, endpoint string, body interface{}) (*http.Response, error) {
	// Construct the full URL
	baseURL, err := url.Parse(c.BaseUrl)
	if err != nil {
		return nil, fmt.Errorf("invalid base URL: %w", err)
	}

	baseURL.Path = path.Join(baseURL.Path, endpoint)

	var reqBody io.Reader
	if body != nil {
		jsonData, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal request body: %w", err)
		}
		reqBody = bytes.NewBuffer(jsonData)
	}

	req, err := http.NewRequestWithContext(ctx, method, baseURL.String(), reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Set headers
	req.Header.Set("Authorization", "Token "+c.ApiToken)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	// Perform the request
	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to perform request: %w", err)
	}

	return resp, nil
}

// HandleErrorResponse processes error responses from the API
func (c *MakeAPIClient) HandleErrorResponse(resp *http.Response) error {
	defer func() { _ = resp.Body.Close() }()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, resp.Status)
	}

	var errorResp ErrorResponse
	if err := json.Unmarshal(body, &errorResp); err != nil {
		return fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(body))
	}

	message := errorResp.Message
	if message == "" {
		message = errorResp.Error
	}
	if message == "" {
		message = string(body)
	}

	return fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, message)
}

// CreateScenario creates a new scenario in Make.com
func (c *MakeAPIClient) CreateScenario(ctx context.Context, req ScenarioRequest) (*ScenarioResponse, error) {
	resp, err := c.MakeRequest(ctx, "POST", "v2/scenarios", req)
	if err != nil {
		return nil, err
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode >= 400 {
		return nil, c.HandleErrorResponse(resp)
	}

	var scenario ScenarioResponse
	if err := json.NewDecoder(resp.Body).Decode(&scenario); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &scenario, nil
}

// GetScenario retrieves a scenario by ID from Make.com
func (c *MakeAPIClient) GetScenario(ctx context.Context, id string) (*ScenarioResponse, error) {
	endpoint := fmt.Sprintf("v2/scenarios/%s", id)
	resp, err := c.MakeRequest(ctx, "GET", endpoint, nil)
	if err != nil {
		return nil, err
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode == 404 {
		return nil, fmt.Errorf("scenario with ID %s not found", id)
	}

	if resp.StatusCode >= 400 {
		return nil, c.HandleErrorResponse(resp)
	}

	var scenario ScenarioResponse
	if err := json.NewDecoder(resp.Body).Decode(&scenario); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &scenario, nil
}

// UpdateScenario updates an existing scenario in Make.com
func (c *MakeAPIClient) UpdateScenario(ctx context.Context, id string, req ScenarioRequest) (*ScenarioResponse, error) {
	endpoint := fmt.Sprintf("v2/scenarios/%s", id)
	resp, err := c.MakeRequest(ctx, "PUT", endpoint, req)
	if err != nil {
		return nil, err
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode == 404 {
		return nil, fmt.Errorf("scenario with ID %s not found", id)
	}

	if resp.StatusCode >= 400 {
		return nil, c.HandleErrorResponse(resp)
	}

	var scenario ScenarioResponse
	if err := json.NewDecoder(resp.Body).Decode(&scenario); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &scenario, nil
}

// DeleteScenario deletes a scenario from Make.com
func (c *MakeAPIClient) DeleteScenario(ctx context.Context, id string) error {
	endpoint := fmt.Sprintf("v2/scenarios/%s", id)
	resp, err := c.MakeRequest(ctx, "DELETE", endpoint, nil)
	if err != nil {
		return err
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode == 404 {
		// Already deleted or doesn't exist
		return nil
	}

	if resp.StatusCode >= 400 {
		return c.HandleErrorResponse(resp)
	}

	return nil
}

// ConnectionResponse represents a Make.com connection from the API
type ConnectionResponse struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	AppName  string `json:"app_name"`
	TeamID   string `json:"team_id,omitempty"`
	Verified bool   `json:"verified"`
}

// ConnectionRequest represents the request payload for creating connections
type ConnectionRequest struct {
	Name     string                 `json:"name"`
	AppName  string                 `json:"app_name"`
	TeamID   string                 `json:"team_id,omitempty"`
	Settings map[string]interface{} `json:"settings,omitempty"`
}

// CreateConnection creates a new connection in Make.com
func (c *MakeAPIClient) CreateConnection(ctx context.Context, req ConnectionRequest) (*ConnectionResponse, error) {
	resp, err := c.MakeRequest(ctx, "POST", "v2/connections", req)
	if err != nil {
		return nil, err
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode >= 400 {
		return nil, c.HandleErrorResponse(resp)
	}

	var connection ConnectionResponse
	if err := json.NewDecoder(resp.Body).Decode(&connection); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &connection, nil
}

// GetConnection retrieves a connection by ID from Make.com
func (c *MakeAPIClient) GetConnection(ctx context.Context, id string) (*ConnectionResponse, error) {
	endpoint := fmt.Sprintf("v2/connections/%s", id)
	resp, err := c.MakeRequest(ctx, "GET", endpoint, nil)
	if err != nil {
		return nil, err
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode == 404 {
		return nil, fmt.Errorf("connection with ID %s not found", id)
	}

	if resp.StatusCode >= 400 {
		return nil, c.HandleErrorResponse(resp)
	}

	var connection ConnectionResponse
	if err := json.NewDecoder(resp.Body).Decode(&connection); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &connection, nil
}

// UpdateConnection updates an existing connection in Make.com
func (c *MakeAPIClient) UpdateConnection(ctx context.Context, id string, req ConnectionRequest) (*ConnectionResponse, error) {
	endpoint := fmt.Sprintf("v2/connections/%s", id)
	resp, err := c.MakeRequest(ctx, "PUT", endpoint, req)
	if err != nil {
		return nil, err
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode == 404 {
		return nil, fmt.Errorf("connection with ID %s not found", id)
	}

	if resp.StatusCode >= 400 {
		return nil, c.HandleErrorResponse(resp)
	}

	var connection ConnectionResponse
	if err := json.NewDecoder(resp.Body).Decode(&connection); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &connection, nil
}

// DeleteConnection deletes a connection from Make.com
func (c *MakeAPIClient) DeleteConnection(ctx context.Context, id string) error {
	endpoint := fmt.Sprintf("v2/connections/%s", id)
	resp, err := c.MakeRequest(ctx, "DELETE", endpoint, nil)
	if err != nil {
		return err
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode == 404 {
		// Already deleted or doesn't exist
		return nil
	}

	if resp.StatusCode >= 400 {
		return c.HandleErrorResponse(resp)
	}

	return nil
}

// WebhookResponse represents a Make.com webhook from the API
type WebhookResponse struct {
	ID       string                 `json:"id"`
	Name     string                 `json:"name"`
	URL      string                 `json:"url"`
	TeamID   string                 `json:"team_id,omitempty"`
	Active   bool                   `json:"active"`
	Settings map[string]interface{} `json:"settings,omitempty"`
}

// WebhookRequest represents the request payload for creating/updating webhooks
type WebhookRequest struct {
	Name     string                 `json:"name"`
	URL      string                 `json:"url"`
	TeamID   string                 `json:"team_id,omitempty"`
	Active   bool                   `json:"active"`
	Settings map[string]interface{} `json:"settings,omitempty"`
}

// CreateWebhook creates a new webhook in Make.com
func (c *MakeAPIClient) CreateWebhook(ctx context.Context, req WebhookRequest) (*WebhookResponse, error) {
	resp, err := c.MakeRequest(ctx, "POST", "v2/webhooks", req)
	if err != nil {
		return nil, err
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode >= 400 {
		return nil, c.HandleErrorResponse(resp)
	}

	var webhook WebhookResponse
	if err := json.NewDecoder(resp.Body).Decode(&webhook); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &webhook, nil
}

// GetWebhook retrieves a webhook by ID from Make.com
func (c *MakeAPIClient) GetWebhook(ctx context.Context, id string) (*WebhookResponse, error) {
	endpoint := fmt.Sprintf("v2/webhooks/%s", id)
	resp, err := c.MakeRequest(ctx, "GET", endpoint, nil)
	if err != nil {
		return nil, err
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode == 404 {
		return nil, fmt.Errorf("webhook with ID %s not found", id)
	}

	if resp.StatusCode >= 400 {
		return nil, c.HandleErrorResponse(resp)
	}

	var webhook WebhookResponse
	if err := json.NewDecoder(resp.Body).Decode(&webhook); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &webhook, nil
}

// UpdateWebhook updates an existing webhook in Make.com
func (c *MakeAPIClient) UpdateWebhook(ctx context.Context, id string, req WebhookRequest) (*WebhookResponse, error) {
	endpoint := fmt.Sprintf("v2/webhooks/%s", id)
	resp, err := c.MakeRequest(ctx, "PUT", endpoint, req)
	if err != nil {
		return nil, err
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode == 404 {
		return nil, fmt.Errorf("webhook with ID %s not found", id)
	}

	if resp.StatusCode >= 400 {
		return nil, c.HandleErrorResponse(resp)
	}

	var webhook WebhookResponse
	if err := json.NewDecoder(resp.Body).Decode(&webhook); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &webhook, nil
}

// DeleteWebhook deletes a webhook from Make.com
func (c *MakeAPIClient) DeleteWebhook(ctx context.Context, id string) error {
	endpoint := fmt.Sprintf("v2/webhooks/%s", id)
	resp, err := c.MakeRequest(ctx, "DELETE", endpoint, nil)
	if err != nil {
		return err
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode == 404 {
		// Already deleted or doesn't exist
		return nil
	}

	if resp.StatusCode >= 400 {
		return c.HandleErrorResponse(resp)
	}

	return nil
}

// TeamResponse represents a Make.com team from the API
type TeamResponse struct {
	ID             string `json:"id"`
	Name           string `json:"name"`
	OrganizationID string `json:"organization_id,omitempty"`
}

// TeamRequest represents the request payload for creating/updating teams
type TeamRequest struct {
	Name           string `json:"name"`
	OrganizationID string `json:"organization_id,omitempty"`
}

// CreateTeam creates a new team in Make.com
func (c *MakeAPIClient) CreateTeam(ctx context.Context, req TeamRequest) (*TeamResponse, error) {
	resp, err := c.MakeRequest(ctx, "POST", "v2/teams", req)
	if err != nil {
		return nil, err
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode >= 400 {
		return nil, c.HandleErrorResponse(resp)
	}

	var team TeamResponse
	if err := json.NewDecoder(resp.Body).Decode(&team); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &team, nil
}

// GetTeam retrieves a team by ID from Make.com
func (c *MakeAPIClient) GetTeam(ctx context.Context, id string) (*TeamResponse, error) {
	endpoint := fmt.Sprintf("v2/teams/%s", id)
	resp, err := c.MakeRequest(ctx, "GET", endpoint, nil)
	if err != nil {
		return nil, err
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode == 404 {
		return nil, fmt.Errorf("team with ID %s not found", id)
	}

	if resp.StatusCode >= 400 {
		return nil, c.HandleErrorResponse(resp)
	}

	var team TeamResponse
	if err := json.NewDecoder(resp.Body).Decode(&team); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &team, nil
}

// UpdateTeam updates an existing team in Make.com
func (c *MakeAPIClient) UpdateTeam(ctx context.Context, id string, req TeamRequest) (*TeamResponse, error) {
	endpoint := fmt.Sprintf("v2/teams/%s", id)
	resp, err := c.MakeRequest(ctx, "PUT", endpoint, req)
	if err != nil {
		return nil, err
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode == 404 {
		return nil, fmt.Errorf("team with ID %s not found", id)
	}

	if resp.StatusCode >= 400 {
		return nil, c.HandleErrorResponse(resp)
	}

	var team TeamResponse
	if err := json.NewDecoder(resp.Body).Decode(&team); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &team, nil
}

// DeleteTeam deletes a team from Make.com
func (c *MakeAPIClient) DeleteTeam(ctx context.Context, id string) error {
	endpoint := fmt.Sprintf("v2/teams/%s", id)
	resp, err := c.MakeRequest(ctx, "DELETE", endpoint, nil)
	if err != nil {
		return err
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode == 404 {
		// Already deleted or doesn't exist
		return nil
	}

	if resp.StatusCode >= 400 {
		return c.HandleErrorResponse(resp)
	}

	return nil
}

// OrganizationResponse represents a Make.com organization from the API
type OrganizationResponse struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// OrganizationRequest represents the request payload for creating/updating organizations
type OrganizationRequest struct {
	Name string `json:"name"`
}

// CreateOrganization creates a new organization in Make.com
func (c *MakeAPIClient) CreateOrganization(ctx context.Context, req OrganizationRequest) (*OrganizationResponse, error) {
	resp, err := c.MakeRequest(ctx, "POST", "v2/organizations", req)
	if err != nil {
		return nil, err
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode >= 400 {
		return nil, c.HandleErrorResponse(resp)
	}

	var org OrganizationResponse
	if err := json.NewDecoder(resp.Body).Decode(&org); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &org, nil
}

// GetOrganization retrieves an organization by ID from Make.com
func (c *MakeAPIClient) GetOrganization(ctx context.Context, id string) (*OrganizationResponse, error) {
	endpoint := fmt.Sprintf("v2/organizations/%s", id)
	resp, err := c.MakeRequest(ctx, "GET", endpoint, nil)
	if err != nil {
		return nil, err
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode == 404 {
		return nil, fmt.Errorf("organization with ID %s not found", id)
	}

	if resp.StatusCode >= 400 {
		return nil, c.HandleErrorResponse(resp)
	}

	var org OrganizationResponse
	if err := json.NewDecoder(resp.Body).Decode(&org); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &org, nil
}

// UpdateOrganization updates an existing organization in Make.com
func (c *MakeAPIClient) UpdateOrganization(ctx context.Context, id string, req OrganizationRequest) (*OrganizationResponse, error) {
	endpoint := fmt.Sprintf("v2/organizations/%s", id)
	resp, err := c.MakeRequest(ctx, "PUT", endpoint, req)
	if err != nil {
		return nil, err
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode == 404 {
		return nil, fmt.Errorf("organization with ID %s not found", id)
	}

	if resp.StatusCode >= 400 {
		return nil, c.HandleErrorResponse(resp)
	}

	var org OrganizationResponse
	if err := json.NewDecoder(resp.Body).Decode(&org); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &org, nil
}

// DeleteOrganization deletes an organization from Make.com
func (c *MakeAPIClient) DeleteOrganization(ctx context.Context, id string) error {
	endpoint := fmt.Sprintf("v2/organizations/%s", id)
	resp, err := c.MakeRequest(ctx, "DELETE", endpoint, nil)
	if err != nil {
		return err
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode == 404 {
		// Already deleted or doesn't exist
		return nil
	}

	if resp.StatusCode >= 400 {
		return c.HandleErrorResponse(resp)
	}

	return nil
}

// DataStoreResponse represents a Make.com data store from the API
type DataStoreResponse struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	TeamID      string `json:"team_id,omitempty"`
}

// DataStoreRequest represents the request payload for creating/updating data stores
type DataStoreRequest struct {
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	TeamID      string `json:"team_id,omitempty"`
}

// CreateDataStore creates a new data store in Make.com
func (c *MakeAPIClient) CreateDataStore(ctx context.Context, req DataStoreRequest) (*DataStoreResponse, error) {
	resp, err := c.MakeRequest(ctx, "POST", "v2/data-stores", req)
	if err != nil {
		return nil, err
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode >= 400 {
		return nil, c.HandleErrorResponse(resp)
	}

	var ds DataStoreResponse
	if err := json.NewDecoder(resp.Body).Decode(&ds); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &ds, nil
}

// GetDataStore retrieves a data store by ID from Make.com
func (c *MakeAPIClient) GetDataStore(ctx context.Context, id string) (*DataStoreResponse, error) {
	endpoint := fmt.Sprintf("v2/data-stores/%s", id)
	resp, err := c.MakeRequest(ctx, "GET", endpoint, nil)
	if err != nil {
		return nil, err
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode == 404 {
		return nil, fmt.Errorf("data store with ID %s not found", id)
	}

	if resp.StatusCode >= 400 {
		return nil, c.HandleErrorResponse(resp)
	}

	var ds DataStoreResponse
	if err := json.NewDecoder(resp.Body).Decode(&ds); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &ds, nil
}

// UpdateDataStore updates an existing data store in Make.com
func (c *MakeAPIClient) UpdateDataStore(ctx context.Context, id string, req DataStoreRequest) (*DataStoreResponse, error) {
	endpoint := fmt.Sprintf("v2/data-stores/%s", id)
	resp, err := c.MakeRequest(ctx, "PUT", endpoint, req)
	if err != nil {
		return nil, err
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode == 404 {
		return nil, fmt.Errorf("data store with ID %s not found", id)
	}

	if resp.StatusCode >= 400 {
		return nil, c.HandleErrorResponse(resp)
	}

	var ds DataStoreResponse
	if err := json.NewDecoder(resp.Body).Decode(&ds); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &ds, nil
}

// DeleteDataStore deletes a data store from Make.com
func (c *MakeAPIClient) DeleteDataStore(ctx context.Context, id string) error {
	endpoint := fmt.Sprintf("v2/data-stores/%s", id)
	resp, err := c.MakeRequest(ctx, "DELETE", endpoint, nil)
	if err != nil {
		return err
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode == 404 {
		// Already deleted or doesn't exist
		return nil
	}

	if resp.StatusCode >= 400 {
		return c.HandleErrorResponse(resp)
	}

	return nil
}

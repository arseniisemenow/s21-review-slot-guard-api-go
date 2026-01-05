package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"
)

const (
	// DefaultBaseURL is the base URL for the 21-school API
	DefaultBaseURL = "https://platform.21-school.ru"
	// DefaultAuthURL is the base URL for the auth service
	DefaultAuthURL = "https://auth.21-school.ru"
	// GraphQLPath is the path for GraphQL endpoint
	GraphQLPath = "/services/graphql"
	// AuthTokenPath is the path for getting auth token
	AuthTokenPath = "/auth/realms/EduPowerKeycloak/protocol/openid-connect/token"
)

// AuthConfig holds authentication configuration
type AuthConfig struct {
	Login    string
	Password string
}

// TokenResponse is the response from the auth endpoint
type TokenResponse struct {
	AccessToken      string `json:"access_token"`
	TokenType        string `json:"token_type"`
	ExpiresIn        int    `json:"expires_in"`
	RefreshToken     string `json:"refresh_token"`
	IDToken          string `json:"id_token"`
	NotBeforePolicy  int    `json:"not-before-policy"`
	SessionState     string `json:"session_state"`
	Scope            string `json:"scope"`
}

// Client is the 21-school API client
type Client struct {
	httpClient  *http.Client
	baseURL     string
	authURL     string
	authConfig  *AuthConfig
	token       *TokenResponse
	tokenExpiry time.Time
	// Additional headers required for some API operations
	schoolID      string
	userRole      string
	eduProductID  string
	eduOrgUnitID  string
}

// ClientOption is a function that configures a Client
type ClientOption func(*Client)

// WithHTTPClient sets a custom HTTP client
func WithHTTPClient(client *http.Client) ClientOption {
	return func(c *Client) {
		c.httpClient = client
	}
}

// WithBaseURL sets a custom base URL
func WithBaseURL(url string) ClientOption {
	return func(c *Client) {
		c.baseURL = url
	}
}

// WithAuthURL sets a custom auth URL
func WithAuthURL(url string) ClientOption {
	return func(c *Client) {
		c.authURL = url
	}
}

// WithSchoolID sets the school ID header
func WithSchoolID(schoolID string) ClientOption {
	return func(c *Client) {
		c.schoolID = schoolID
	}
}

// WithUserRole sets the user role header
func WithUserRole(role string) ClientOption {
	return func(c *Client) {
		c.userRole = role
	}
}

// WithEduProductID sets the edu product ID header
func WithEduProductID(id string) ClientOption {
	return func(c *Client) {
		c.eduProductID = id
	}
}

// WithEduOrgUnitID sets the edu org unit ID header
func WithEduOrgUnitID(id string) ClientOption {
	return func(c *Client) {
		c.eduOrgUnitID = id
	}
}

// NewClient creates a new API client
func NewClient(authConfig *AuthConfig, opts ...ClientOption) *Client {
	c := &Client{
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
		baseURL:    DefaultBaseURL,
		authURL:    DefaultAuthURL,
		authConfig: authConfig,
	}

	for _, opt := range opts {
		opt(c)
	}

	return c
}

// Authenticate gets a new access token using credentials
func (c *Client) Authenticate(ctx context.Context) (*TokenResponse, error) {
	if c.authConfig == nil {
		return nil, fmt.Errorf("auth config not set")
	}

	authURL := c.authURL + AuthTokenPath
	data := fmt.Sprintf("client_id=s21-open-api&username=%s&password=%s&grant_type=password",
		c.authConfig.Login, c.authConfig.Password)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, authURL, strings.NewReader(data))
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("do request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("auth failed with status %d: %s", resp.StatusCode, string(body))
	}

	var tokenResp TokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&tokenResp); err != nil {
		return nil, fmt.Errorf("decode response: %w", err)
	}

	c.token = &tokenResp
	c.tokenExpiry = time.Now().Add(time.Duration(tokenResp.ExpiresIn) * time.Second)

	return &tokenResp, nil
}

// ensureToken makes sure we have a valid token
func (c *Client) ensureToken(ctx context.Context) error {
	if c.token == nil || time.Until(c.tokenExpiry) < time.Minute {
		_, err := c.Authenticate(ctx)
		return err
	}
	return nil
}

// GraphQLRequest represents a GraphQL request
type GraphQLRequest struct {
	OperationName string                 `json:"operationName,omitempty"`
	Variables     map[string]interface{} `json:"variables,omitempty"`
	Query         string                 `json:"query"`
}

// GraphQLResponse represents a GraphQL response
type GraphQLResponse struct {
	Data   json.RawMessage `json:"data"`
	Errors []GraphQLError  `json:"errors,omitempty"`
}

// GraphQLError represents a GraphQL error
type GraphQLError struct {
	Message string `json:"message"`
	Path    []interface{} `json:"path,omitempty"`
	Extensions map[string]interface{} `json:"extensions,omitempty"`
}

// Do executes a GraphQL request
func (c *Client) Do(ctx context.Context, req *GraphQLRequest, resp interface{}) error {
	if err := c.ensureToken(ctx); err != nil {
		return fmt.Errorf("ensure token: %w", err)
	}

	body, err := json.Marshal(req)
	if err != nil {
		return fmt.Errorf("marshal request: %w", err)
	}

	url := c.baseURL + GraphQLPath
	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("create request: %w", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", "Bearer "+c.token.AccessToken)

	// Add additional headers if set
	if c.schoolID != "" {
		httpReq.Header.Set("schoolid", c.schoolID)
	}
	if c.userRole != "" {
		httpReq.Header.Set("userrole", c.userRole)
	}
	if c.eduProductID != "" {
		httpReq.Header.Set("x-edu-product-id", c.eduProductID)
	}
	if c.eduOrgUnitID != "" {
		httpReq.Header.Set("x-edu-org-unit-id", c.eduOrgUnitID)
	}

	// Debug logging
	if os.Getenv("DEBUG") != "" {
		fmt.Fprintf(os.Stderr, "Request URL: %s\n", url)
		fmt.Fprintf(os.Stderr, "Request Body: %s\n", string(body))
	}

	httpResp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return fmt.Errorf("do request: %w", err)
	}
	defer httpResp.Body.Close()

	if httpResp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(httpResp.Body)
	 errMsg := fmt.Sprintf("request failed with status %d", httpResp.StatusCode)
		if len(body) > 0 {
			errMsg += fmt.Sprintf(": %s", string(body))
		} else {
			errMsg += " (empty response body)"
		}
		if os.Getenv("DEBUG") != "" {
			fmt.Fprintf(os.Stderr, "Response Headers: %v\n", httpResp.Header)
			fmt.Fprintf(os.Stderr, "Response Body Length: %d\n", len(body))
		}
		return fmt.Errorf("%s", errMsg)
	}

	var gqlResp GraphQLResponse
	if err := json.NewDecoder(httpResp.Body).Decode(&gqlResp); err != nil {
		return fmt.Errorf("decode response: %w", err)
	}

	if len(gqlResp.Errors) > 0 {
		return fmt.Errorf("graphql errors: %v", gqlResp.Errors)
	}

	if resp != nil {
		if err := json.Unmarshal(gqlResp.Data, resp); err != nil {
			return fmt.Errorf("unmarshal data: %w", err)
		}
	}

	return nil
}

// SetToken allows setting a token manually (useful for testing)
func (c *Client) SetToken(token *TokenResponse, expiry time.Time) {
	c.token = token
	c.tokenExpiry = expiry
}

// GetToken returns the current token (useful for testing)
func (c *Client) GetToken() *TokenResponse {
	return c.token
}

// SetBaseURL sets the base URL (useful for testing)
func (c *Client) SetBaseURL(url string) {
	c.baseURL = url
}

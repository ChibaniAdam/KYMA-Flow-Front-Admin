package ldap

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/sirupsen/logrus"
)

// Client represents a client for the LDAP Manager service
type Client struct {
	baseURL    string
	httpClient *http.Client
	logger     *logrus.Logger
}

// User represents a user from LDAP Manager service
type User struct {
	UID          string   `json:"uid"`
	CN           string   `json:"cn"`
	SN           string   `json:"sn"`
	GivenName    string   `json:"givenName"`
	Mail         string   `json:"mail"`
	Department   string   `json:"department"`
	UIDNumber    int      `json:"uidNumber"`
	GIDNumber    int      `json:"gidNumber"`
	HomeDir      string   `json:"homeDirectory"`
	Repositories []string `json:"repositories"`
	DN           string   `json:"dn"`
}

// Department represents a department from LDAP Manager service
type Department struct {
	OU           string   `json:"ou"`
	Description  string   `json:"description"`
	Manager      string   `json:"manager,omitempty"`
	Members      []string `json:"members"`
	Repositories []string `json:"repositories"`
	DN           string   `json:"dn"`
}

// NewClient creates a new LDAP Manager client
func NewClient(baseURL string, timeout time.Duration, logger *logrus.Logger) *Client {
	return &Client{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: timeout,
		},
		logger: logger,
	}
}

// doGraphQLRequest performs a GraphQL request to LDAP Manager
func (c *Client) doGraphQLRequest(ctx context.Context, query string, token string) (map[string]interface{}, error) {
	requestBody := map[string]interface{}{
		"query": query,
	}

	bodyBytes, err := json.Marshal(requestBody)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	url := fmt.Sprintf("%s/graphql", c.baseURL)
	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewReader(bodyBytes))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	if token != "" {
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	}

	c.logger.WithFields(logrus.Fields{
		"url": url,
	}).Debug("Making LDAP Manager request")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		c.logger.WithFields(logrus.Fields{
			"status": resp.StatusCode,
			"body":   string(body),
		}).Error("LDAP Manager API error")
		return nil, fmt.Errorf("LDAP Manager API error: status %d", resp.StatusCode)
	}

	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	// Check for GraphQL errors
	if errors, ok := result["errors"].([]interface{}); ok && len(errors) > 0 {
		return nil, fmt.Errorf("GraphQL error: %v", errors)
	}

	return result, nil
}

// GetUser gets a user by UID from LDAP Manager
func (c *Client) GetUser(ctx context.Context, uid string, token string) (*User, error) {
	query := fmt.Sprintf(`
		query {
			user(uid: "%s") {
				uid
				cn
				sn
				givenName
				mail
				department
				uidNumber
				gidNumber
				homeDirectory
				repositories
				dn
			}
		}
	`, uid)

	result, err := c.doGraphQLRequest(ctx, query, token)
	if err != nil {
		return nil, err
	}

	data, ok := result["data"].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid response format")
	}

	userMap, ok := data["user"].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("user not found")
	}

	userBytes, err := json.Marshal(userMap)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal user: %w", err)
	}

	var user User
	if err := json.Unmarshal(userBytes, &user); err != nil {
		return nil, fmt.Errorf("failed to unmarshal user: %w", err)
	}

	return &user, nil
}

// GetDepartment gets a department by OU from LDAP Manager
func (c *Client) GetDepartment(ctx context.Context, ou string, token string) (*Department, error) {
	query := fmt.Sprintf(`
		query {
			department(ou: "%s") {
				ou
				description
				manager
				members
				repositories
				dn
			}
		}
	`, ou)

	result, err := c.doGraphQLRequest(ctx, query, token)
	if err != nil {
		return nil, err
	}

	data, ok := result["data"].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid response format")
	}

	deptMap, ok := data["department"].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("department not found")
	}

	deptBytes, err := json.Marshal(deptMap)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal department: %w", err)
	}

	var dept Department
	if err := json.Unmarshal(deptBytes, &dept); err != nil {
		return nil, fmt.Errorf("failed to unmarshal department: %w", err)
	}

	return &dept, nil
}

// HealthCheck checks if LDAP Manager service is accessible
func (c *Client) HealthCheck(ctx context.Context) error {
	url := fmt.Sprintf("%s/health", c.baseURL)
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return fmt.Errorf("failed to create health check request: %w", err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("health check failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("LDAP Manager unhealthy: status %d", resp.StatusCode)
	}

	return nil
}

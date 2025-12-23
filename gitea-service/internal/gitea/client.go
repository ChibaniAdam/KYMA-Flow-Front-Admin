package gitea

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/sirupsen/logrus"
)

// Client represents a Gitea API client
type Client struct {
	baseURL    string
	token      string
	httpClient *http.Client
	logger     *logrus.Logger
}

// Repository represents a Gitea repository
type Repository struct {
	ID            int64     `json:"id"`
	Name          string    `json:"name"`
	FullName      string    `json:"full_name"`
	Description   string    `json:"description"`
	Private       bool      `json:"private"`
	Fork          bool      `json:"fork"`
	HTMLURL       string    `json:"html_url"`
	SSHURL        string    `json:"ssh_url"`
	CloneURL      string    `json:"clone_url"`
	DefaultBranch string    `json:"default_branch"`
	Language      string    `json:"language"`
	Stars         int       `json:"stars_count"`
	Forks         int       `json:"forks_count"`
	Size          int       `json:"size"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
	Owner         Owner     `json:"owner"`
}

// Owner represents the repository owner
type Owner struct {
	ID       int64  `json:"id"`
	Login    string `json:"login"`
	FullName string `json:"full_name"`
	Email    string `json:"email"`
	AvatarURL string `json:"avatar_url"`
}

// NewClient creates a new Gitea API client
func NewClient(baseURL, token string, logger *logrus.Logger) *Client {
	return &Client{
		baseURL: baseURL,
		token:   token,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
		logger: logger,
	}
}

// doRequest performs an HTTP request to Gitea API
func (c *Client) doRequest(method, path string) ([]byte, error) {
	url := fmt.Sprintf("%s/api/v1%s", c.baseURL, path)

	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", fmt.Sprintf("token %s", c.token))
	req.Header.Set("Content-Type", "application/json")

	c.logger.WithFields(logrus.Fields{
		"method": method,
		"url":    url,
	}).Debug("Making Gitea API request")

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
		}).Error("Gitea API error")
		return nil, fmt.Errorf("gitea API error: %s (status: %d)", string(body), resp.StatusCode)
	}

	return body, nil
}

// ListRepositories lists all repositories accessible by the admin token
func (c *Client) ListRepositories() ([]*Repository, error) {
	// Use /repos/search endpoint to get all repositories
	body, err := c.doRequest("GET", "/repos/search?limit=1000")
	if err != nil {
		return nil, err
	}

	var result struct {
		Data []*Repository `json:"data"`
	}

	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	c.logger.WithField("count", len(result.Data)).Info("Fetched repositories from Gitea")
	return result.Data, nil
}

// GetRepository gets a specific repository by owner and name
func (c *Client) GetRepository(owner, name string) (*Repository, error) {
	path := fmt.Sprintf("/repos/%s/%s", owner, name)
	body, err := c.doRequest("GET", path)
	if err != nil {
		return nil, err
	}

	var repo Repository
	if err := json.Unmarshal(body, &repo); err != nil {
		return nil, fmt.Errorf("failed to parse repository: %w", err)
	}

	return &repo, nil
}

// SearchRepositories searches repositories by query
func (c *Client) SearchRepositories(query string, limit int) ([]*Repository, error) {
	if limit <= 0 {
		limit = 50
	}

	path := fmt.Sprintf("/repos/search?q=%s&limit=%d", query, limit)
	body, err := c.doRequest("GET", path)
	if err != nil {
		return nil, err
	}

	var result struct {
		Data []*Repository `json:"data"`
	}

	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return result.Data, nil
}

// HealthCheck checks if Gitea API is accessible
func (c *Client) HealthCheck() error {
	_, err := c.doRequest("GET", "/version")
	return err
}

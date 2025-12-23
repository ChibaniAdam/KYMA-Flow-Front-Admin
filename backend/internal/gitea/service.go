package gitea

import (
	"context"
	"fmt"
	"strings"

	"github.com/devplatform/ldap-manager/internal/ldap"
	"github.com/devplatform/ldap-manager/internal/models"
	"github.com/sirupsen/logrus"
)

// Service provides repository access control based on LDAP attributes
type Service struct {
	client  *Client
	ldapMgr *ldap.Manager
	logger  *logrus.Logger
}

// NewService creates a new Gitea service
func NewService(client *Client, ldapMgr *ldap.Manager, logger *logrus.Logger) *Service {
	return &Service{
		client:  client,
		ldapMgr: ldapMgr,
		logger:  logger,
	}
}

// GetUserRepositories gets all repositories accessible by a user
// User can access a repository if:
// 1. It's in their personal LDAP githubRepository attribute, OR
// 2. It's in their department's githubRepository attribute
func (s *Service) GetUserRepositories(ctx context.Context, user *models.User) ([]*Repository, error) {
	// Get all repositories from Gitea
	allRepos, err := s.client.ListRepositories()
	if err != nil {
		return nil, fmt.Errorf("failed to list repositories: %w", err)
	}

	// Get user's allowed repository names
	allowedRepoNames := s.getUserAllowedRepos(ctx, user)

	// Filter repositories
	userRepos := make([]*Repository, 0)
	for _, repo := range allRepos {
		if s.isRepoAllowed(repo, allowedRepoNames) {
			userRepos = append(userRepos, repo)
		}
	}

	s.logger.WithFields(logrus.Fields{
		"uid":          user.UID,
		"total_repos":  len(allRepos),
		"user_repos":   len(userRepos),
	}).Info("Filtered user repositories")

	return userRepos, nil
}

// GetRepository gets a specific repository if user has access
func (s *Service) GetRepository(ctx context.Context, user *models.User, owner, name string) (*Repository, error) {
	// Get repository from Gitea
	repo, err := s.client.GetRepository(owner, name)
	if err != nil {
		return nil, err
	}

	// Check if user has access
	allowedRepoNames := s.getUserAllowedRepos(ctx, user)
	if !s.isRepoAllowed(repo, allowedRepoNames) {
		return nil, fmt.Errorf("access denied to repository: %s/%s", owner, name)
	}

	return repo, nil
}

// SearchUserRepositories searches repositories accessible by user
func (s *Service) SearchUserRepositories(ctx context.Context, user *models.User, query string, limit int) ([]*Repository, error) {
	// Get user's accessible repos first
	userRepos, err := s.GetUserRepositories(ctx, user)
	if err != nil {
		return nil, err
	}

	// Filter by search query
	if query == "" {
		// Return all user repos if no query
		if limit > 0 && len(userRepos) > limit {
			return userRepos[:limit], nil
		}
		return userRepos, nil
	}

	// Search within user's accessible repos
	queryLower := strings.ToLower(query)
	results := make([]*Repository, 0)
	for _, repo := range userRepos {
		if s.matchesQuery(repo, queryLower) {
			results = append(results, repo)
			if limit > 0 && len(results) >= limit {
				break
			}
		}
	}

	return results, nil
}

// getUserAllowedRepos gets the list of repository names/patterns the user can access
func (s *Service) getUserAllowedRepos(ctx context.Context, user *models.User) map[string]bool {
	allowedRepos := make(map[string]bool)

	// Add user's personal repositories
	for _, repo := range user.Repositories {
		allowedRepos[s.normalizeRepoName(repo)] = true
	}

	// Add department repositories
	if user.Department != "" {
		dept, err := s.ldapMgr.GetDepartment(ctx, user.Department)
		if err == nil {
			for _, repo := range dept.Repositories {
				allowedRepos[s.normalizeRepoName(repo)] = true
			}
		} else {
			s.logger.WithError(err).Warn("Failed to get department repositories")
		}
	}

	return allowedRepos
}

// isRepoAllowed checks if a repository is in the allowed list
func (s *Service) isRepoAllowed(repo *Repository, allowedRepos map[string]bool) bool {
	// Check full name (owner/repo)
	if allowedRepos[strings.ToLower(repo.FullName)] {
		return true
	}

	// Check by name only
	if allowedRepos[strings.ToLower(repo.Name)] {
		return true
	}

	return false
}

// normalizeRepoName normalizes a repository name for comparison
// Handles formats like:
// - "owner/repo"
// - "https://github.com/owner/repo"
// - "https://gitea.example.com/owner/repo"
// - "repo" (just the name)
func (s *Service) normalizeRepoName(repoName string) string {
	// Remove trailing slashes
	repoName = strings.TrimSuffix(repoName, "/")

	// If it's a URL, extract owner/repo part
	if strings.HasPrefix(repoName, "http://") || strings.HasPrefix(repoName, "https://") {
		parts := strings.Split(repoName, "/")
		if len(parts) >= 5 {
			// Extract owner/repo from URL (last two parts)
			repoName = fmt.Sprintf("%s/%s", parts[len(parts)-2], parts[len(parts)-1])
		}
	}

	return strings.ToLower(repoName)
}

// matchesQuery checks if a repository matches the search query
func (s *Service) matchesQuery(repo *Repository, queryLower string) bool {
	return strings.Contains(strings.ToLower(repo.Name), queryLower) ||
		strings.Contains(strings.ToLower(repo.FullName), queryLower) ||
		strings.Contains(strings.ToLower(repo.Description), queryLower)
}

// GetRepositoryStats gets statistics about user's repositories
func (s *Service) GetRepositoryStats(ctx context.Context, user *models.User) (*RepositoryStats, error) {
	repos, err := s.GetUserRepositories(ctx, user)
	if err != nil {
		return nil, err
	}

	stats := &RepositoryStats{
		TotalCount: len(repos),
		Languages:  make(map[string]int),
	}

	for _, repo := range repos {
		if repo.Private {
			stats.PrivateCount++
		} else {
			stats.PublicCount++
		}

		if repo.Language != "" {
			stats.Languages[repo.Language]++
		}
	}

	return stats, nil
}

// RepositoryStats contains statistics about repositories
type RepositoryStats struct {
	TotalCount   int            `json:"totalCount"`
	PrivateCount int            `json:"privateCount"`
	PublicCount  int            `json:"publicCount"`
	Languages    map[string]int `json:"languages"`
}

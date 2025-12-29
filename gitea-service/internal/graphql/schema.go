package graphql

import (
	"context"
	"fmt"
	"time"

	"github.com/devplatform/gitea-service/internal/config"
	"github.com/devplatform/gitea-service/internal/gitea"
	"github.com/devplatform/gitea-service/internal/ldap"
	"github.com/devplatform/gitea-service/internal/models"
	"github.com/golang-jwt/jwt/v5"
	"github.com/graphql-go/graphql"
	"github.com/sirupsen/logrus"
)

// Schema represents the GraphQL schema
type Schema struct {
	schema       graphql.Schema
	giteaService *gitea.Service
	ldapClient   *ldap.Client
	giteaClient  *gitea.Client
	config       *config.Config
	logger       *logrus.Logger
}

// JWT Claims (must match LDAP Manager for token validation)
type Claims struct {
	UID        string `json:"uid"`
	Mail       string `json:"mail"`
	Department string `json:"department"`
	jwt.RegisteredClaims
}

// NewSchema creates a new GraphQL schema
func NewSchema(giteaService *gitea.Service, ldapClient *ldap.Client, giteaClient *gitea.Client, cfg *config.Config, logger *logrus.Logger) *Schema {
	s := &Schema{
		giteaService: giteaService,
		ldapClient:   ldapClient,
		giteaClient:  giteaClient,
		config:       cfg,
		logger:       logger,
	}

	// Define types
	giteaRepoType := s.defineGiteaRepositoryType()
	repoStatsType := s.defineRepositoryStatsType()
	healthType := s.defineHealthType()
	paginatedReposType := s.definePaginatedRepositoriesType(giteaRepoType)

	// Define root query
	queryType := graphql.NewObject(graphql.ObjectConfig{
		Name: "Query",
		Fields: graphql.Fields{
			"listRepositories": &graphql.Field{
				Type: paginatedReposType,
				Args: graphql.FieldConfigArgument{
					"limit": &graphql.ArgumentConfig{
						Type:         graphql.Int,
						DefaultValue: 10,
						Description:  "Number of items per page (default: 10, max: 100)",
					},
					"offset": &graphql.ArgumentConfig{
						Type:         graphql.Int,
						DefaultValue: 0,
						Description:  "Number of items to skip",
					},
				},
				Resolve: s.resolveListRepositories,
			},
			"searchRepositories": &graphql.Field{
				Type: paginatedReposType,
				Args: graphql.FieldConfigArgument{
					"query": &graphql.ArgumentConfig{
						Type:        graphql.String,
						Description: "Search query string",
					},
					"limit": &graphql.ArgumentConfig{
						Type:         graphql.Int,
						DefaultValue: 10,
						Description:  "Number of items per page (default: 10, max: 100)",
					},
					"offset": &graphql.ArgumentConfig{
						Type:         graphql.Int,
						DefaultValue: 0,
						Description:  "Number of items to skip",
					},
				},
				Resolve: s.resolveSearchRepositories,
			},
			"getRepository": &graphql.Field{
				Type: giteaRepoType,
				Args: graphql.FieldConfigArgument{
					"owner": &graphql.ArgumentConfig{
						Type:        graphql.NewNonNull(graphql.String),
						Description: "Repository owner username",
					},
					"name": &graphql.ArgumentConfig{
						Type:        graphql.NewNonNull(graphql.String),
						Description: "Repository name",
					},
				},
				Resolve: s.resolveGetRepository,
			},
			"myRepositories": &graphql.Field{
				Type: paginatedReposType,
				Args: graphql.FieldConfigArgument{
					"limit": &graphql.ArgumentConfig{
						Type:         graphql.Int,
						DefaultValue: 10,
						Description:  "Number of items per page",
					},
					"offset": &graphql.ArgumentConfig{
						Type:         graphql.Int,
						DefaultValue: 0,
						Description:  "Number of items to skip",
					},
				},
				Resolve: s.resolveMyRepositories,
			},
			"repositoryStats": &graphql.Field{
				Type:    repoStatsType,
				Resolve: s.resolveRepositoryStats,
			},
			"health": &graphql.Field{
				Type:    healthType,
				Resolve: s.resolveHealth,
			},
		},
	})

	// Define root mutation
	mutationType := graphql.NewObject(graphql.ObjectConfig{
		Name: "Mutation",
		Fields: graphql.Fields{
			"deleteRepository": &graphql.Field{
				Type: graphql.Boolean,
				Args: graphql.FieldConfigArgument{
					"owner": &graphql.ArgumentConfig{
						Type:        graphql.NewNonNull(graphql.String),
						Description: "Repository owner username",
					},
					"name": &graphql.ArgumentConfig{
						Type:        graphql.NewNonNull(graphql.String),
						Description: "Repository name",
					},
				},
				Resolve: s.resolveDeleteRepository,
			},
			"updateRepository": &graphql.Field{
				Type: giteaRepoType,
				Args: graphql.FieldConfigArgument{
					"owner": &graphql.ArgumentConfig{
						Type:        graphql.NewNonNull(graphql.String),
						Description: "Repository owner username",
					},
					"name": &graphql.ArgumentConfig{
						Type:        graphql.NewNonNull(graphql.String),
						Description: "Repository name",
					},
					"description": &graphql.ArgumentConfig{
						Type:        graphql.String,
						Description: "Repository description",
					},
					"private": &graphql.ArgumentConfig{
						Type:        graphql.Boolean,
						Description: "Make repository private",
					},
					"defaultBranch": &graphql.ArgumentConfig{
						Type:        graphql.String,
						Description: "Default branch name",
					},
				},
				Resolve: s.resolveUpdateRepository,
			},
		},
	})

	// Create schema
	schemaConfig := graphql.SchemaConfig{
		Query:    queryType,
		Mutation: mutationType,
	}

	schema, err := graphql.NewSchema(schemaConfig)
	if err != nil {
		logger.WithError(err).Fatal("Failed to create schema")
	}

	s.schema = schema
	return s
}

// GetSchema returns the GraphQL schema
func (s *Schema) GetSchema() graphql.Schema {
	return s.schema
}

// Type Definitions

func (s *Schema) defineGiteaRepositoryType() *graphql.Object {
	ownerType := graphql.NewObject(graphql.ObjectConfig{
		Name: "RepositoryOwner",
		Fields: graphql.Fields{
			"id":        &graphql.Field{Type: graphql.Int},
			"login":     &graphql.Field{Type: graphql.String},
			"fullName":  &graphql.Field{Type: graphql.String},
			"email":     &graphql.Field{Type: graphql.String},
			"avatarUrl": &graphql.Field{Type: graphql.String},
		},
	})

	return graphql.NewObject(graphql.ObjectConfig{
		Name: "GiteaRepository",
		Fields: graphql.Fields{
			"id":            &graphql.Field{Type: graphql.Int},
			"name":          &graphql.Field{Type: graphql.String},
			"fullName":      &graphql.Field{Type: graphql.String},
			"description":   &graphql.Field{Type: graphql.String},
			"private":       &graphql.Field{Type: graphql.Boolean},
			"fork":          &graphql.Field{Type: graphql.Boolean},
			"htmlUrl":       &graphql.Field{Type: graphql.String},
			"sshUrl":        &graphql.Field{Type: graphql.String},
			"cloneUrl":      &graphql.Field{Type: graphql.String},
			"defaultBranch": &graphql.Field{Type: graphql.String},
			"language":      &graphql.Field{Type: graphql.String},
			"stars":         &graphql.Field{Type: graphql.Int},
			"forks":         &graphql.Field{Type: graphql.Int},
			"size":          &graphql.Field{Type: graphql.Int},
			"createdAt":     &graphql.Field{Type: graphql.String},
			"updatedAt":     &graphql.Field{Type: graphql.String},
			"owner":         &graphql.Field{Type: ownerType},
		},
	})
}

func (s *Schema) defineRepositoryStatsType() *graphql.Object {
	languageDistType := graphql.NewObject(graphql.ObjectConfig{
		Name: "LanguageDistribution",
		Fields: graphql.Fields{
			"language": &graphql.Field{Type: graphql.String},
			"count":    &graphql.Field{Type: graphql.Int},
		},
	})

	return graphql.NewObject(graphql.ObjectConfig{
		Name: "RepositoryStats",
		Fields: graphql.Fields{
			"totalCount":   &graphql.Field{Type: graphql.Int},
			"privateCount": &graphql.Field{Type: graphql.Int},
			"publicCount":  &graphql.Field{Type: graphql.Int},
			"languages":    &graphql.Field{Type: graphql.NewList(languageDistType)},
		},
	})
}

func (s *Schema) defineHealthType() *graphql.Object {
	return graphql.NewObject(graphql.ObjectConfig{
		Name: "HealthStatus",
		Fields: graphql.Fields{
			"status":      &graphql.Field{Type: graphql.String},
			"timestamp":   &graphql.Field{Type: graphql.Int},
			"gitea":       &graphql.Field{Type: graphql.Boolean},
			"ldapManager": &graphql.Field{Type: graphql.Boolean},
		},
	})
}

func (s *Schema) definePaginatedRepositoriesType(repoType *graphql.Object) *graphql.Object {
	return graphql.NewObject(graphql.ObjectConfig{
		Name: "PaginatedRepositories",
		Fields: graphql.Fields{
			"items": &graphql.Field{
				Type:        graphql.NewList(repoType),
				Description: "List of repositories",
			},
			"total": &graphql.Field{
				Type:        graphql.Int,
				Description: "Total number of repositories",
			},
			"limit": &graphql.Field{
				Type:        graphql.Int,
				Description: "Number of items per page",
			},
			"offset": &graphql.Field{
				Type:        graphql.Int,
				Description: "Number of items skipped",
			},
			"hasMore": &graphql.Field{
				Type:        graphql.Boolean,
				Description: "Whether there are more items",
			},
		},
	})
}

// Query Resolvers

func (s *Schema) resolveListRepositories(p graphql.ResolveParams) (interface{}, error) {
	limit := p.Args["limit"].(int)
	offset := p.Args["offset"].(int)

	if limit > 100 {
		limit = 100
	}
	if limit <= 0 {
		limit = 10
	}

	allRepos, err := s.giteaClient.ListRepositories()
	if err != nil {
		s.logger.WithError(err).Error("Failed to list repositories")
		return nil, fmt.Errorf("failed to list repositories: %w", err)
	}

	total := len(allRepos)
	start := offset
	end := offset + limit

	if start > total {
		start = total
	}
	if end > total {
		end = total
	}

	paginatedRepos := allRepos[start:end]

	return map[string]interface{}{
		"items":   s.convertGiteaReposToMap(paginatedRepos),
		"total":   total,
		"limit":   limit,
		"offset":  offset,
		"hasMore": end < total,
	}, nil
}

func (s *Schema) resolveSearchRepositories(p graphql.ResolveParams) (interface{}, error) {
	query := ""
	if q, ok := p.Args["query"].(string); ok {
		query = q
	}

	limit := p.Args["limit"].(int)
	offset := p.Args["offset"].(int)

	if limit > 100 {
		limit = 100
	}
	if limit <= 0 {
		limit = 10
	}

	// Fetch more to handle offset client-side
	fetchLimit := limit + offset + 50
	allRepos, err := s.giteaClient.SearchRepositories(query, fetchLimit)
	if err != nil {
		s.logger.WithError(err).Error("Failed to search repositories")
		return nil, fmt.Errorf("failed to search repositories: %w", err)
	}

	total := len(allRepos)
	start := offset
	end := offset + limit

	if start > total {
		start = total
	}
	if end > total {
		end = total
	}

	paginatedRepos := allRepos[start:end]

	return map[string]interface{}{
		"items":   s.convertGiteaReposToMap(paginatedRepos),
		"total":   total,
		"limit":   limit,
		"offset":  offset,
		"hasMore": end < total,
	}, nil
}

func (s *Schema) resolveGetRepository(p graphql.ResolveParams) (interface{}, error) {
	owner := p.Args["owner"].(string)
	name := p.Args["name"].(string)

	repo, err := s.giteaClient.GetRepository(owner, name)
	if err != nil {
		s.logger.WithError(err).Error("Failed to get repository")
		return nil, fmt.Errorf("failed to get repository: %w", err)
	}

	return s.convertGiteaRepoToMap(repo), nil
}

func (s *Schema) resolveMyRepositories(p graphql.ResolveParams) (interface{}, error) {
	user, token, err := s.getUserFromContext(p.Context)
	if err != nil {
		return nil, err
	}

	limit := p.Args["limit"].(int)
	offset := p.Args["offset"].(int)

	if limit > 100 {
		limit = 100
	}
	if limit <= 0 {
		limit = 10
	}

	repos, err := s.giteaService.GetUserRepositories(p.Context, user, token)
	if err != nil {
		s.logger.WithError(err).Error("Failed to get user repositories")
		return nil, fmt.Errorf("failed to get repositories: %w", err)
	}

	convertedRepos := s.convertReposToModels(repos)
	total := len(convertedRepos)
	start := offset
	end := offset + limit

	if start > total {
		start = total
	}
	if end > total {
		end = total
	}

	paginatedRepos := convertedRepos[start:end]

	return map[string]interface{}{
		"items":   paginatedRepos,
		"total":   total,
		"limit":   limit,
		"offset":  offset,
		"hasMore": end < total,
	}, nil
}

func (s *Schema) resolveDeleteRepository(p graphql.ResolveParams) (interface{}, error) {
	owner := p.Args["owner"].(string)
	name := p.Args["name"].(string)

	err := s.giteaClient.DeleteRepository(owner, name)
	if err != nil {
		s.logger.WithError(err).Error("Failed to delete repository")
		return false, fmt.Errorf("failed to delete repository: %w", err)
	}

	s.logger.WithFields(logrus.Fields{
		"owner": owner,
		"name":  name,
	}).Info("Repository deleted successfully")

	return true, nil
}

func (s *Schema) resolveUpdateRepository(p graphql.ResolveParams) (interface{}, error) {
	owner := p.Args["owner"].(string)
	name := p.Args["name"].(string)

	updates := make(map[string]interface{})
	if desc, ok := p.Args["description"].(string); ok {
		updates["description"] = desc
	}
	if private, ok := p.Args["private"].(bool); ok {
		updates["private"] = private
	}
	if branch, ok := p.Args["defaultBranch"].(string); ok {
		updates["default_branch"] = branch
	}

	repo, err := s.giteaClient.UpdateRepository(owner, name, updates)
	if err != nil {
		s.logger.WithError(err).Error("Failed to update repository")
		return nil, fmt.Errorf("failed to update repository: %w", err)
	}

	s.logger.WithFields(logrus.Fields{
		"owner": owner,
		"name":  name,
	}).Info("Repository updated successfully")

	return s.convertGiteaRepoToMap(repo), nil
}

func (s *Schema) resolveRepositoryStats(p graphql.ResolveParams) (interface{}, error) {
	user, token, err := s.getUserFromContext(p.Context)
	if err != nil {
		return nil, err
	}

	stats, err := s.giteaService.GetRepositoryStats(p.Context, user, token)
	if err != nil {
		s.logger.WithError(err).Error("Failed to get repository stats")
		return nil, fmt.Errorf("failed to get stats: %w", err)
	}

	return s.convertStatsToModel(stats), nil
}

func (s *Schema) resolveHealth(p graphql.ResolveParams) (interface{}, error) {
	giteaHealthy := s.giteaClient.HealthCheck() == nil
	ldapHealthy := s.ldapClient.HealthCheck(p.Context) == nil

	status := "healthy"
	if !giteaHealthy || !ldapHealthy {
		status = "unhealthy"
	}

	return &models.HealthStatus{
		Status:      status,
		Timestamp:   time.Now().Unix(),
		Gitea:       giteaHealthy,
		LDAPManager: ldapHealthy,
	}, nil
}

// Helper functions

func (s *Schema) getUserFromContext(ctx context.Context) (*models.User, string, error) {
	// Get user from context (set by auth middleware)
	user, ok := ctx.Value("user").(*models.User)
	if !ok || user == nil {
		return nil, "", fmt.Errorf("unauthorized")
	}

	// Get token from context
	token, ok := ctx.Value("token").(string)
	if !ok {
		token = ""
	}

	return user, token, nil
}

// ExtractUserFromToken validates JWT and fetches user from LDAP Manager
func (s *Schema) ExtractUserFromToken(ctx context.Context, tokenString string) (*models.User, error) {
	claims := &Claims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(s.config.JWTSecret), nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	// Fetch full user details from LDAP Manager
	ldapUser, err := s.ldapClient.GetUser(ctx, claims.UID, tokenString)
	if err != nil {
		return nil, fmt.Errorf("failed to get user from LDAP Manager: %w", err)
	}

	// Convert to our User model
	return &models.User{
		UID:          ldapUser.UID,
		CN:           ldapUser.CN,
		SN:           ldapUser.SN,
		GivenName:    ldapUser.GivenName,
		Mail:         ldapUser.Mail,
		Department:   ldapUser.Department,
		UIDNumber:    ldapUser.UIDNumber,
		GIDNumber:    ldapUser.GIDNumber,
		HomeDir:      ldapUser.HomeDir,
		Repositories: ldapUser.Repositories,
		DN:           ldapUser.DN,
	}, nil
}

// Converter functions

func (s *Schema) convertRepoToModel(repo *gitea.Repository) *models.GiteaRepository {
	if repo == nil {
		return nil
	}

	return &models.GiteaRepository{
		ID:            repo.ID,
		Name:          repo.Name,
		FullName:      repo.FullName,
		Description:   repo.Description,
		Private:       repo.Private,
		Fork:          repo.Fork,
		HTMLURL:       repo.HTMLURL,
		SSHURL:        repo.SSHURL,
		CloneURL:      repo.CloneURL,
		DefaultBranch: repo.DefaultBranch,
		Language:      repo.Language,
		Stars:         repo.Stars,
		Forks:         repo.Forks,
		Size:          repo.Size,
		CreatedAt:     repo.CreatedAt.Format(time.RFC3339),
		UpdatedAt:     repo.UpdatedAt.Format(time.RFC3339),
		Owner: models.RepositoryOwner{
			ID:        repo.Owner.ID,
			Login:     repo.Owner.Login,
			FullName:  repo.Owner.FullName,
			Email:     repo.Owner.Email,
			AvatarURL: repo.Owner.AvatarURL,
		},
	}
}

func (s *Schema) convertReposToModels(repos []*gitea.Repository) []*models.GiteaRepository {
	result := make([]*models.GiteaRepository, len(repos))
	for i, r := range repos {
		result[i] = s.convertRepoToModel(r)
	}
	return result
}

func (s *Schema) convertStatsToModel(stats *gitea.RepositoryStats) *models.RepositoryStats {
	if stats == nil {
		return nil
	}

	languages := make([]models.LanguageDistribution, 0, len(stats.Languages))
	for lang, count := range stats.Languages {
		languages = append(languages, models.LanguageDistribution{
			Language: lang,
			Count:    count,
		})
	}

	return &models.RepositoryStats{
		TotalCount:   stats.TotalCount,
		PrivateCount: stats.PrivateCount,
		PublicCount:  stats.PublicCount,
		Languages:    languages,
	}
}

// Helper function to convert Gitea repos to map format
func (s *Schema) convertGiteaReposToMap(repos []*gitea.Repository) []map[string]interface{} {
	result := make([]map[string]interface{}, len(repos))
	for i, repo := range repos {
		result[i] = s.convertGiteaRepoToMap(repo)
	}
	return result
}

func (s *Schema) convertGiteaRepoToMap(repo *gitea.Repository) map[string]interface{} {
	return map[string]interface{}{
		"id":            repo.ID,
		"name":          repo.Name,
		"fullName":      repo.FullName,
		"description":   repo.Description,
		"private":       repo.Private,
		"fork":          repo.Fork,
		"htmlUrl":       repo.HTMLURL,
		"sshUrl":        repo.SSHURL,
		"cloneUrl":      repo.CloneURL,
		"defaultBranch": repo.DefaultBranch,
		"language":      repo.Language,
		"stars":         repo.Stars,
		"forks":         repo.Forks,
		"size":          repo.Size,
		"createdAt":     repo.CreatedAt.Format(time.RFC3339),
		"updatedAt":     repo.UpdatedAt.Format(time.RFC3339),
		"owner": map[string]interface{}{
			"id":        repo.Owner.ID,
			"login":     repo.Owner.Login,
			"fullName":  repo.Owner.FullName,
			"email":     repo.Owner.Email,
			"avatarUrl": repo.Owner.AvatarURL,
		},
	}
}

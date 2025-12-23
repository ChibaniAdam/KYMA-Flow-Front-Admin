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

	// Define root query
	queryType := graphql.NewObject(graphql.ObjectConfig{
		Name: "Query",
		Fields: graphql.Fields{
			"myRepositories": &graphql.Field{
				Type:    graphql.NewList(giteaRepoType),
				Resolve: s.resolveMyRepositories,
			},
			"repository": &graphql.Field{
				Type: giteaRepoType,
				Args: graphql.FieldConfigArgument{
					"owner": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.String),
					},
					"name": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.String),
					},
				},
				Resolve: s.resolveRepository,
			},
			"searchRepositories": &graphql.Field{
				Type: graphql.NewList(giteaRepoType),
				Args: graphql.FieldConfigArgument{
					"query": &graphql.ArgumentConfig{
						Type: graphql.String,
					},
					"limit": &graphql.ArgumentConfig{
						Type: graphql.Int,
					},
				},
				Resolve: s.resolveSearchRepositories,
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

	// Create schema
	schemaConfig := graphql.SchemaConfig{
		Query: queryType,
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

// Query Resolvers

func (s *Schema) resolveMyRepositories(p graphql.ResolveParams) (interface{}, error) {
	user, token, err := s.getUserFromContext(p.Context)
	if err != nil {
		return nil, err
	}

	repos, err := s.giteaService.GetUserRepositories(p.Context, user, token)
	if err != nil {
		s.logger.WithError(err).Error("Failed to get user repositories")
		return nil, fmt.Errorf("failed to get repositories: %w", err)
	}

	return s.convertReposToModels(repos), nil
}

func (s *Schema) resolveRepository(p graphql.ResolveParams) (interface{}, error) {
	user, token, err := s.getUserFromContext(p.Context)
	if err != nil {
		return nil, err
	}

	owner := p.Args["owner"].(string)
	name := p.Args["name"].(string)

	repo, err := s.giteaService.GetRepository(p.Context, user, owner, name, token)
	if err != nil {
		s.logger.WithError(err).Error("Failed to get repository")
		return nil, fmt.Errorf("failed to get repository: %w", err)
	}

	return s.convertRepoToModel(repo), nil
}

func (s *Schema) resolveSearchRepositories(p graphql.ResolveParams) (interface{}, error) {
	user, token, err := s.getUserFromContext(p.Context)
	if err != nil {
		return nil, err
	}

	query := ""
	if q, ok := p.Args["query"].(string); ok {
		query = q
	}

	limit := 50
	if l, ok := p.Args["limit"].(int); ok {
		limit = l
	}

	repos, err := s.giteaService.SearchUserRepositories(p.Context, user, query, limit, token)
	if err != nil {
		s.logger.WithError(err).Error("Failed to search repositories")
		return nil, fmt.Errorf("failed to search repositories: %w", err)
	}

	return s.convertReposToModels(repos), nil
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

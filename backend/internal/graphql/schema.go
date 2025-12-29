package graphql

import (
	"context"
	"fmt"
	"time"

	"github.com/devplatform/ldap-manager/internal/config"
	"github.com/devplatform/ldap-manager/internal/ldap"
	"github.com/devplatform/ldap-manager/internal/models"
	"github.com/golang-jwt/jwt/v5"
	"github.com/graphql-go/graphql"
	"github.com/sirupsen/logrus"
)

// Schema represents the GraphQL schema
type Schema struct {
	schema  graphql.Schema
	ldapMgr *ldap.Manager
	config  *config.Config
	logger  *logrus.Logger
}

// JWT Claims
type Claims struct {
	UID        string `json:"uid"`
	Mail       string `json:"mail"`
	Department string `json:"department"`
	jwt.RegisteredClaims
}

// NewSchema creates a new GraphQL schema
func NewSchema(ldapMgr *ldap.Manager, cfg *config.Config, logger *logrus.Logger) *Schema {
	s := &Schema{
		ldapMgr: ldapMgr,
		config:  cfg,
		logger:  logger,
	}

	// Define types
	userType := s.defineUserType()
	departmentType := s.defineDepartmentType()
	groupType := s.defineGroupType()
	authPayloadType := s.defineAuthPayloadType(userType)
	statsType := s.defineStatsType()
	healthType := s.defineHealthType()

	// Define paginated types
	paginatedUsersType := s.definePaginatedUsersType(userType)
	paginatedDepartmentsType := s.definePaginatedDepartmentsType(departmentType)
	paginatedGroupsType := s.definePaginatedGroupsType(groupType)

	// Define input types
	createUserInputType := s.defineCreateUserInput()
	updateUserInputType := s.defineUpdateUserInput()
	createDepartmentInputType := s.defineCreateDepartmentInput()
	searchFilterInputType := s.defineSearchFilterInput()

	// Define root query
	queryType := graphql.NewObject(graphql.ObjectConfig{
		Name: "Query",
		Fields: graphql.Fields{
			"me": &graphql.Field{
				Type:    userType,
				Resolve: s.resolveMe,
			},
			"user": &graphql.Field{
				Type: userType,
				Args: graphql.FieldConfigArgument{
					"uid": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.String),
					},
				},
				Resolve: s.resolveUser,
			},
			"users": &graphql.Field{
				Type: paginatedUsersType,
				Args: graphql.FieldConfigArgument{
					"filter": &graphql.ArgumentConfig{
						Type:        searchFilterInputType,
						Description: "Search filter for users",
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
				Resolve: s.resolveUsers,
			},
			"department": &graphql.Field{
				Type: departmentType,
				Args: graphql.FieldConfigArgument{
					"ou": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.String),
					},
				},
				Resolve: s.resolveDepartment,
			},
			"departments": &graphql.Field{
				Type: paginatedDepartmentsType,
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
				Resolve: s.resolveDepartments,
			},
			"departmentUsers": &graphql.Field{
				Type: paginatedUsersType,
				Args: graphql.FieldConfigArgument{
					"department": &graphql.ArgumentConfig{
						Type:        graphql.NewNonNull(graphql.String),
						Description: "Department name to filter by",
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
				Resolve: s.resolveDepartmentUsers,
			},
			"group": &graphql.Field{
				Type: groupType,
				Args: graphql.FieldConfigArgument{
					"cn": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.String),
					},
				},
				Resolve: s.resolveGroup,
			},
			"groups": &graphql.Field{
				Type: paginatedGroupsType,
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
				Resolve: s.resolveGroups,
			},
			"health": &graphql.Field{
				Type:    healthType,
				Resolve: s.resolveHealth,
			},
			"stats": &graphql.Field{
				Type:    statsType,
				Resolve: s.resolveStats,
			},
		},
	})

	// Define root mutation
	mutationType := graphql.NewObject(graphql.ObjectConfig{
		Name: "Mutation",
		Fields: graphql.Fields{
			"login": &graphql.Field{
				Type: authPayloadType,
				Args: graphql.FieldConfigArgument{
					"uid": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.String),
					},
					"password": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.String),
					},
				},
				Resolve: s.resolveLogin,
			},
			"createUser": &graphql.Field{
				Type: userType,
				Args: graphql.FieldConfigArgument{
					"input": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(createUserInputType),
					},
				},
				Resolve: s.resolveCreateUser,
			},
			"updateUser": &graphql.Field{
				Type: userType,
				Args: graphql.FieldConfigArgument{
					"input": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(updateUserInputType),
					},
				},
				Resolve: s.resolveUpdateUser,
			},
			"deleteUser": &graphql.Field{
				Type: graphql.Boolean,
				Args: graphql.FieldConfigArgument{
					"uid": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.String),
					},
				},
				Resolve: s.resolveDeleteUser,
			},
			"createDepartment": &graphql.Field{
				Type: departmentType,
				Args: graphql.FieldConfigArgument{
					"input": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(createDepartmentInputType),
					},
				},
				Resolve: s.resolveCreateDepartment,
			},
			"deleteDepartment": &graphql.Field{
				Type: graphql.Boolean,
				Args: graphql.FieldConfigArgument{
					"ou": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.String),
					},
				},
				Resolve: s.resolveDeleteDepartment,
			},
			"assignRepoToDepartment": &graphql.Field{
				Type: departmentType,
				Args: graphql.FieldConfigArgument{
					"ou": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.String),
					},
					"repositories": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.NewList(graphql.NewNonNull(graphql.String))),
					},
				},
				Resolve: s.resolveAssignRepoToDepartment,
			},
			"assignRepoToUser": &graphql.Field{
				Type: userType,
				Args: graphql.FieldConfigArgument{
					"uid": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.String),
					},
					"repositories": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.NewList(graphql.NewNonNull(graphql.String))),
					},
				},
				Resolve: s.resolveAssignRepoToUser,
			},
			"createGroup": &graphql.Field{
				Type: groupType,
				Args: graphql.FieldConfigArgument{
					"cn": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.String),
					},
					"description": &graphql.ArgumentConfig{
						Type: graphql.String,
					},
				},
				Resolve: s.resolveCreateGroup,
			},
			"addUserToGroup": &graphql.Field{
				Type: graphql.Boolean,
				Args: graphql.FieldConfigArgument{
					"uid": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.String),
					},
					"groupCn": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.String),
					},
				},
				Resolve: s.resolveAddUserToGroup,
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

func (s *Schema) defineUserType() *graphql.Object {
	return graphql.NewObject(graphql.ObjectConfig{
		Name: "User",
		Fields: graphql.Fields{
			"uid":           &graphql.Field{Type: graphql.String},
			"cn":            &graphql.Field{Type: graphql.String},
			"sn":            &graphql.Field{Type: graphql.String},
			"givenName":     &graphql.Field{Type: graphql.String},
			"mail":          &graphql.Field{Type: graphql.String},
			"department":    &graphql.Field{Type: graphql.String},
			"uidNumber":     &graphql.Field{Type: graphql.Int},
			"gidNumber":     &graphql.Field{Type: graphql.Int},
			"homeDirectory": &graphql.Field{Type: graphql.String},
			"repositories":  &graphql.Field{Type: graphql.NewList(graphql.String)},
			"dn":            &graphql.Field{Type: graphql.String},
		},
	})
}

func (s *Schema) defineDepartmentType() *graphql.Object {
	return graphql.NewObject(graphql.ObjectConfig{
		Name: "Department",
		Fields: graphql.Fields{
			"ou":           &graphql.Field{Type: graphql.String},
			"description":  &graphql.Field{Type: graphql.String},
			"manager":      &graphql.Field{Type: graphql.String},
			"members":      &graphql.Field{Type: graphql.NewList(graphql.String)},
			"repositories": &graphql.Field{Type: graphql.NewList(graphql.String)},
			"dn":           &graphql.Field{Type: graphql.String},
		},
	})
}

func (s *Schema) defineGroupType() *graphql.Object {
	return graphql.NewObject(graphql.ObjectConfig{
		Name: "Group",
		Fields: graphql.Fields{
			"cn":        &graphql.Field{Type: graphql.String},
			"gidNumber": &graphql.Field{Type: graphql.Int},
			"members":   &graphql.Field{Type: graphql.NewList(graphql.String)},
			"dn":        &graphql.Field{Type: graphql.String},
		},
	})
}

func (s *Schema) defineAuthPayloadType(userType *graphql.Object) *graphql.Object {
	return graphql.NewObject(graphql.ObjectConfig{
		Name: "AuthPayload",
		Fields: graphql.Fields{
			"token": &graphql.Field{Type: graphql.String},
			"user":  &graphql.Field{Type: userType},
		},
	})
}

func (s *Schema) defineStatsType() *graphql.Object {
	return graphql.NewObject(graphql.ObjectConfig{
		Name: "Stats",
		Fields: graphql.Fields{
			"poolSize":      &graphql.Field{Type: graphql.Int},
			"available":     &graphql.Field{Type: graphql.Int},
			"inUse":         &graphql.Field{Type: graphql.Int},
			"totalRequests": &graphql.Field{Type: graphql.Int},
		},
	})
}

func (s *Schema) defineHealthType() *graphql.Object {
	return graphql.NewObject(graphql.ObjectConfig{
		Name: "HealthStatus",
		Fields: graphql.Fields{
			"status":    &graphql.Field{Type: graphql.String},
			"timestamp": &graphql.Field{Type: graphql.Int},
			"ldap":      &graphql.Field{Type: graphql.Boolean},
			"gitea":     &graphql.Field{Type: graphql.Boolean},
		},
	})
}

// Paginated Type Definitions

func (s *Schema) definePaginatedUsersType(userType *graphql.Object) *graphql.Object {
	return graphql.NewObject(graphql.ObjectConfig{
		Name: "PaginatedUsers",
		Fields: graphql.Fields{
			"items": &graphql.Field{
				Type:        graphql.NewList(userType),
				Description: "List of users",
			},
			"total": &graphql.Field{
				Type:        graphql.Int,
				Description: "Total number of users",
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

func (s *Schema) definePaginatedDepartmentsType(departmentType *graphql.Object) *graphql.Object {
	return graphql.NewObject(graphql.ObjectConfig{
		Name: "PaginatedDepartments",
		Fields: graphql.Fields{
			"items": &graphql.Field{
				Type:        graphql.NewList(departmentType),
				Description: "List of departments",
			},
			"total": &graphql.Field{
				Type:        graphql.Int,
				Description: "Total number of departments",
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

func (s *Schema) definePaginatedGroupsType(groupType *graphql.Object) *graphql.Object {
	return graphql.NewObject(graphql.ObjectConfig{
		Name: "PaginatedGroups",
		Fields: graphql.Fields{
			"items": &graphql.Field{
				Type:        graphql.NewList(groupType),
				Description: "List of groups",
			},
			"total": &graphql.Field{
				Type:        graphql.Int,
				Description: "Total number of groups",
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

// Input Type Definitions

func (s *Schema) defineCreateUserInput() *graphql.InputObject {
	return graphql.NewInputObject(graphql.InputObjectConfig{
		Name: "CreateUserInput",
		Fields: graphql.InputObjectConfigFieldMap{
			"uid":          &graphql.InputObjectFieldConfig{Type: graphql.NewNonNull(graphql.String)},
			"cn":           &graphql.InputObjectFieldConfig{Type: graphql.NewNonNull(graphql.String)},
			"sn":           &graphql.InputObjectFieldConfig{Type: graphql.NewNonNull(graphql.String)},
			"givenName":    &graphql.InputObjectFieldConfig{Type: graphql.NewNonNull(graphql.String)},
			"mail":         &graphql.InputObjectFieldConfig{Type: graphql.NewNonNull(graphql.String)},
			"department":   &graphql.InputObjectFieldConfig{Type: graphql.NewNonNull(graphql.String)},
			"password":     &graphql.InputObjectFieldConfig{Type: graphql.NewNonNull(graphql.String)},
			"repositories": &graphql.InputObjectFieldConfig{Type: graphql.NewList(graphql.String)},
		},
	})
}

func (s *Schema) defineUpdateUserInput() *graphql.InputObject {
	return graphql.NewInputObject(graphql.InputObjectConfig{
		Name: "UpdateUserInput",
		Fields: graphql.InputObjectConfigFieldMap{
			"uid":          &graphql.InputObjectFieldConfig{Type: graphql.NewNonNull(graphql.String)},
			"cn":           &graphql.InputObjectFieldConfig{Type: graphql.String},
			"sn":           &graphql.InputObjectFieldConfig{Type: graphql.String},
			"givenName":    &graphql.InputObjectFieldConfig{Type: graphql.String},
			"mail":         &graphql.InputObjectFieldConfig{Type: graphql.String},
			"department":   &graphql.InputObjectFieldConfig{Type: graphql.String},
			"password":     &graphql.InputObjectFieldConfig{Type: graphql.String},
			"repositories": &graphql.InputObjectFieldConfig{Type: graphql.NewList(graphql.String)},
		},
	})
}

func (s *Schema) defineCreateDepartmentInput() *graphql.InputObject {
	return graphql.NewInputObject(graphql.InputObjectConfig{
		Name: "CreateDepartmentInput",
		Fields: graphql.InputObjectConfigFieldMap{
			"ou":           &graphql.InputObjectFieldConfig{Type: graphql.NewNonNull(graphql.String)},
			"description":  &graphql.InputObjectFieldConfig{Type: graphql.String},
			"manager":      &graphql.InputObjectFieldConfig{Type: graphql.String},
			"repositories": &graphql.InputObjectFieldConfig{Type: graphql.NewList(graphql.String)},
		},
	})
}

func (s *Schema) defineSearchFilterInput() *graphql.InputObject {
	return graphql.NewInputObject(graphql.InputObjectConfig{
		Name: "SearchFilterInput",
		Fields: graphql.InputObjectConfigFieldMap{
			"uid":          &graphql.InputObjectFieldConfig{Type: graphql.String, Description: "Filter by username (UID)"},
			"cn":           &graphql.InputObjectFieldConfig{Type: graphql.String, Description: "Filter by common name (full name)"},
			"sn":           &graphql.InputObjectFieldConfig{Type: graphql.String, Description: "Filter by surname"},
			"givenName":    &graphql.InputObjectFieldConfig{Type: graphql.String, Description: "Filter by given name (first name)"},
			"mail":         &graphql.InputObjectFieldConfig{Type: graphql.String, Description: "Filter by email address"},
			"department":   &graphql.InputObjectFieldConfig{Type: graphql.String, Description: "Filter by department name"},
			"uidNumber":    &graphql.InputObjectFieldConfig{Type: graphql.Int, Description: "Filter by UID number"},
			"gidNumber":    &graphql.InputObjectFieldConfig{Type: graphql.Int, Description: "Filter by GID number"},
			"repository":   &graphql.InputObjectFieldConfig{Type: graphql.String, Description: "Filter by repository URL"},
		},
	})
}

// Query Resolvers

func (s *Schema) resolveMe(p graphql.ResolveParams) (interface{}, error) {
	user, ok := p.Context.Value("user").(*models.User)
	if !ok {
		return nil, fmt.Errorf("unauthorized")
	}
	return user, nil
}

func (s *Schema) resolveUser(p graphql.ResolveParams) (interface{}, error) {
	uid := p.Args["uid"].(string)
	return s.ldapMgr.GetUser(p.Context, uid)
}

func (s *Schema) resolveUsers(p graphql.ResolveParams) (interface{}, error) {
	// Get pagination parameters
	limit := p.Args["limit"].(int)
	offset := p.Args["offset"].(int)

	// Enforce limit constraints
	if limit > 100 {
		limit = 100
	}
	if limit <= 0 {
		limit = 10
	}

	// Parse filter
	var filter *models.SearchFilter
	if filterInput, ok := p.Args["filter"].(map[string]interface{}); ok {
		filter = &models.SearchFilter{}
		if uid, ok := filterInput["uid"].(string); ok {
			filter.UID = uid
		}
		if cn, ok := filterInput["cn"].(string); ok {
			filter.CN = cn
		}
		if sn, ok := filterInput["sn"].(string); ok {
			filter.SN = sn
		}
		if givenName, ok := filterInput["givenName"].(string); ok {
			filter.GivenName = givenName
		}
		if mail, ok := filterInput["mail"].(string); ok {
			filter.Mail = mail
		}
		if dept, ok := filterInput["department"].(string); ok {
			filter.Department = dept
		}
		if uidNumber, ok := filterInput["uidNumber"].(int); ok {
			filter.UIDNumber = uidNumber
		}
		if gidNumber, ok := filterInput["gidNumber"].(int); ok {
			filter.GIDNumber = gidNumber
		}
		if repo, ok := filterInput["repository"].(string); ok {
			filter.Repository = repo
		}
	}

	// Get all users matching filter
	allUsers, err := s.ldapMgr.ListUsers(p.Context, filter)
	if err != nil {
		s.logger.WithError(err).Error("Failed to list users")
		return nil, fmt.Errorf("failed to list users: %w", err)
	}

	// Apply pagination
	total := len(allUsers)
	start := offset
	end := offset + limit

	if start > total {
		start = total
	}
	if end > total {
		end = total
	}

	paginatedUsers := allUsers[start:end]

	return map[string]interface{}{
		"items":   paginatedUsers,
		"total":   total,
		"limit":   limit,
		"offset":  offset,
		"hasMore": end < total,
	}, nil
}

func (s *Schema) resolveDepartment(p graphql.ResolveParams) (interface{}, error) {
	ou := p.Args["ou"].(string)
	return s.ldapMgr.GetDepartment(p.Context, ou)
}

func (s *Schema) resolveDepartments(p graphql.ResolveParams) (interface{}, error) {
	// Get pagination parameters
	limit := p.Args["limit"].(int)
	offset := p.Args["offset"].(int)

	// Enforce limit constraints
	if limit > 100 {
		limit = 100
	}
	if limit <= 0 {
		limit = 10
	}

	// Get all departments
	allDepartments, err := s.ldapMgr.ListDepartments(p.Context)
	if err != nil {
		s.logger.WithError(err).Error("Failed to list departments")
		return nil, fmt.Errorf("failed to list departments: %w", err)
	}

	// Apply pagination
	total := len(allDepartments)
	start := offset
	end := offset + limit

	if start > total {
		start = total
	}
	if end > total {
		end = total
	}

	paginatedDepartments := allDepartments[start:end]

	return map[string]interface{}{
		"items":   paginatedDepartments,
		"total":   total,
		"limit":   limit,
		"offset":  offset,
		"hasMore": end < total,
	}, nil
}

func (s *Schema) resolveDepartmentUsers(p graphql.ResolveParams) (interface{}, error) {
	department := p.Args["department"].(string)

	// Get pagination parameters
	limit := p.Args["limit"].(int)
	offset := p.Args["offset"].(int)

	// Enforce limit constraints
	if limit > 100 {
		limit = 100
	}
	if limit <= 0 {
		limit = 10
	}

	// Get all users in department
	allUsers, err := s.ldapMgr.GetUsersByDepartment(p.Context, department)
	if err != nil {
		s.logger.WithError(err).Error("Failed to get department users")
		return nil, fmt.Errorf("failed to get department users: %w", err)
	}

	// Apply pagination
	total := len(allUsers)
	start := offset
	end := offset + limit

	if start > total {
		start = total
	}
	if end > total {
		end = total
	}

	paginatedUsers := allUsers[start:end]

	return map[string]interface{}{
		"items":   paginatedUsers,
		"total":   total,
		"limit":   limit,
		"offset":  offset,
		"hasMore": end < total,
	}, nil
}

func (s *Schema) resolveGroup(p graphql.ResolveParams) (interface{}, error) {
	cn := p.Args["cn"].(string)
	return s.ldapMgr.GetGroup(p.Context, cn)
}

func (s *Schema) resolveHealth(p graphql.ResolveParams) (interface{}, error) {
	ldapHealthy := s.ldapMgr.HealthCheck(p.Context) == nil
	status := "healthy"
	if !ldapHealthy {
		status = "unhealthy"
	}

	return &models.HealthStatus{
		Status:    status,
		Timestamp: time.Now().Unix(),
		LDAP:      ldapHealthy,
	}, nil
}

func (s *Schema) resolveStats(p graphql.ResolveParams) (interface{}, error) {
	return s.ldapMgr.GetStats(), nil
}

// Mutation Resolvers

func (s *Schema) resolveLogin(p graphql.ResolveParams) (interface{}, error) {
	uid := p.Args["uid"].(string)
	password := p.Args["password"].(string)

	user, err := s.ldapMgr.Authenticate(p.Context, uid, password)
	if err != nil {
		s.logger.WithError(err).Warn("Login failed")
		return nil, fmt.Errorf("authentication failed")
	}

	token, err := s.generateJWT(user)
	if err != nil {
		s.logger.WithError(err).Error("Failed to generate JWT")
		return nil, fmt.Errorf("failed to generate token")
	}

	return &models.AuthPayload{
		Token: token,
		User:  user,
	}, nil
}

func (s *Schema) resolveCreateUser(p graphql.ResolveParams) (interface{}, error) {
	inputMap := p.Args["input"].(map[string]interface{})

	input := &models.CreateUserInput{
		UID:        inputMap["uid"].(string),
		CN:         inputMap["cn"].(string),
		SN:         inputMap["sn"].(string),
		GivenName:  inputMap["givenName"].(string),
		Mail:       inputMap["mail"].(string),
		Department: inputMap["department"].(string),
		Password:   inputMap["password"].(string),
	}

	if repos, ok := inputMap["repositories"].([]interface{}); ok {
		input.Repositories = make([]string, len(repos))
		for i, r := range repos {
			input.Repositories[i] = r.(string)
		}
	}

	return s.ldapMgr.CreateUser(p.Context, input)
}

func (s *Schema) resolveUpdateUser(p graphql.ResolveParams) (interface{}, error) {
	inputMap := p.Args["input"].(map[string]interface{})

	input := &models.UpdateUserInput{
		UID: inputMap["uid"].(string),
	}

	if cn, ok := inputMap["cn"].(string); ok {
		input.CN = &cn
	}
	if sn, ok := inputMap["sn"].(string); ok {
		input.SN = &sn
	}
	if givenName, ok := inputMap["givenName"].(string); ok {
		input.GivenName = &givenName
	}
	if mail, ok := inputMap["mail"].(string); ok {
		input.Mail = &mail
	}
	if dept, ok := inputMap["department"].(string); ok {
		input.Department = &dept
	}
	if password, ok := inputMap["password"].(string); ok {
		input.Password = &password
	}
	if repos, ok := inputMap["repositories"].([]interface{}); ok {
		input.Repositories = make([]string, len(repos))
		for i, r := range repos {
			input.Repositories[i] = r.(string)
		}
	}

	return s.ldapMgr.UpdateUser(p.Context, input)
}

func (s *Schema) resolveDeleteUser(p graphql.ResolveParams) (interface{}, error) {
	uid := p.Args["uid"].(string)
	err := s.ldapMgr.DeleteUser(p.Context, uid)
	return err == nil, err
}

func (s *Schema) resolveCreateDepartment(p graphql.ResolveParams) (interface{}, error) {
	inputMap := p.Args["input"].(map[string]interface{})

	input := &models.CreateDepartmentInput{
		OU: inputMap["ou"].(string),
	}

	if desc, ok := inputMap["description"].(string); ok {
		input.Description = desc
	}
	if mgr, ok := inputMap["manager"].(string); ok {
		input.Manager = mgr
	}
	if repos, ok := inputMap["repositories"].([]interface{}); ok {
		input.Repositories = make([]string, len(repos))
		for i, r := range repos {
			input.Repositories[i] = r.(string)
		}
	}

	return s.ldapMgr.CreateDepartment(p.Context, input)
}

func (s *Schema) resolveDeleteDepartment(p graphql.ResolveParams) (interface{}, error) {
	ou := p.Args["ou"].(string)
	err := s.ldapMgr.DeleteDepartment(p.Context, ou)
	return err == nil, err
}

func (s *Schema) resolveAssignRepoToDepartment(p graphql.ResolveParams) (interface{}, error) {
	ou := p.Args["ou"].(string)
	repoInterfaces := p.Args["repositories"].([]interface{})

	repos := make([]string, len(repoInterfaces))
	for i, r := range repoInterfaces {
		repos[i] = r.(string)
	}

	if err := s.ldapMgr.AssignRepositoryToDepartment(p.Context, ou, repos); err != nil {
		return nil, err
	}

	return s.ldapMgr.GetDepartment(p.Context, ou)
}

func (s *Schema) resolveAssignRepoToUser(p graphql.ResolveParams) (interface{}, error) {
	uid := p.Args["uid"].(string)
	repoInterfaces := p.Args["repositories"].([]interface{})

	repos := make([]string, len(repoInterfaces))
	for i, r := range repoInterfaces {
		repos[i] = r.(string)
	}

	input := &models.UpdateUserInput{
		UID:          uid,
		Repositories: repos,
	}

	return s.ldapMgr.UpdateUser(p.Context, input)
}

func (s *Schema) resolveCreateGroup(p graphql.ResolveParams) (interface{}, error) {
	cn := p.Args["cn"].(string)
	description := ""
	if desc, ok := p.Args["description"].(string); ok {
		description = desc
	}

	return s.ldapMgr.CreateGroup(p.Context, cn, description)
}

func (s *Schema) resolveAddUserToGroup(p graphql.ResolveParams) (interface{}, error) {
	uid := p.Args["uid"].(string)
	groupCn := p.Args["groupCn"].(string)

	err := s.ldapMgr.AddUserToGroup(p.Context, uid, groupCn)
	return err == nil, err
}

// JWT Functions

func (s *Schema) generateJWT(user *models.User) (string, error) {
	expirationTime := time.Now().Add(24 * time.Hour)

	claims := &Claims{
		UID:        user.UID,
		Mail:       user.Mail,
		Department: user.Department,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "ldap-manager",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.config.JWTSecret))
}

// ExtractUserFromToken validates JWT and extracts user information
func (s *Schema) ExtractUserFromToken(tokenString string) (*models.User, error) {
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

	// Fetch full user details from LDAP
	ctx := context.Background()
	return s.ldapMgr.GetUser(ctx, claims.UID)
}

func (s *Schema) resolveGroups(p graphql.ResolveParams) (interface{}, error) {
	// Get pagination parameters
	limit := p.Args["limit"].(int)
	offset := p.Args["offset"].(int)

	// Enforce limit constraints
	if limit > 100 {
		limit = 100
	}
	if limit <= 0 {
		limit = 10
	}

	// Get all groups
	allGroups, err := s.ldapMgr.ListGroups(p.Context)
	if err != nil {
		s.logger.WithError(err).Error("Failed to list groups")
		return nil, fmt.Errorf("failed to list groups: %w", err)
	}

	// Apply pagination
	total := len(allGroups)
	start := offset
	end := offset + limit

	if start > total {
		start = total
	}
	if end > total {
		end = total
	}

	paginatedGroups := allGroups[start:end]

	return map[string]interface{}{
		"items":   paginatedGroups,
		"total":   total,
		"limit":   limit,
		"offset":  offset,
		"hasMore": end < total,
	}, nil
}

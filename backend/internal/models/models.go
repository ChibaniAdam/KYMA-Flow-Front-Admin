package models

// User represents an LDAP user with all attributes
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

// Department represents an organizational unit in LDAP
type Department struct {
	OU           string   `json:"ou"`
	Description  string   `json:"description"`
	Manager      string   `json:"manager,omitempty"`
	Members      []string `json:"members"`
	Repositories []string `json:"repositories"`
	DN           string   `json:"dn"`
}

// Group represents an LDAP group
type Group struct {
	CN        string   `json:"cn"`
	GIDNumber int      `json:"gidNumber"`
	Members   []string `json:"members"`
	DN        string   `json:"dn"`
}

// CreateUserInput contains fields for creating a new user
type CreateUserInput struct {
	UID          string   `json:"uid"`
	CN           string   `json:"cn"`
	SN           string   `json:"sn"`
	GivenName    string   `json:"givenName"`
	Mail         string   `json:"mail"`
	Department   string   `json:"department"`
	Password     string   `json:"password"`
	Repositories []string `json:"repositories"`
}

// UpdateUserInput contains fields for updating a user
type UpdateUserInput struct {
	UID          string   `json:"uid"`
	CN           *string  `json:"cn,omitempty"`
	SN           *string  `json:"sn,omitempty"`
	GivenName    *string  `json:"givenName,omitempty"`
	Mail         *string  `json:"mail,omitempty"`
	Department   *string  `json:"department,omitempty"`
	Password     *string  `json:"password,omitempty"`
	Repositories []string `json:"repositories,omitempty"`
}

// CreateDepartmentInput contains fields for creating a department
type CreateDepartmentInput struct {
	OU           string   `json:"ou"`
	Description  string   `json:"description"`
	Manager      string   `json:"manager,omitempty"`
	Repositories []string `json:"repositories,omitempty"`
}

// SearchFilter contains optional filters for user searches
type SearchFilter struct {
	Department string `json:"department,omitempty"`
	Mail       string `json:"mail,omitempty"`
	CN         string `json:"cn,omitempty"`
}

// AuthPayload is returned after successful authentication
type AuthPayload struct {
	Token string `json:"token"`
	User  *User  `json:"user"`
}

// Stats contains connection pool statistics
type Stats struct {
	PoolSize      int `json:"poolSize"`
	Available     int `json:"available"`
	InUse         int `json:"inUse"`
	TotalRequests int `json:"totalRequests"`
}

// HealthStatus represents the health status of the service
type HealthStatus struct {
	Status    string `json:"status"`
	Timestamp int64  `json:"timestamp"`
	LDAP      bool   `json:"ldap"`
	Gitea     bool   `json:"gitea,omitempty"`
}

// GiteaRepository represents a Gitea repository
type GiteaRepository struct {
	ID            int64            `json:"id"`
	Name          string           `json:"name"`
	FullName      string           `json:"fullName"`
	Description   string           `json:"description"`
	Private       bool             `json:"private"`
	Fork          bool             `json:"fork"`
	HTMLURL       string           `json:"htmlUrl"`
	SSHURL        string           `json:"sshUrl"`
	CloneURL      string           `json:"cloneUrl"`
	DefaultBranch string           `json:"defaultBranch"`
	Language      string           `json:"language"`
	Stars         int              `json:"stars"`
	Forks         int              `json:"forks"`
	Size          int              `json:"size"`
	CreatedAt     string           `json:"createdAt"`
	UpdatedAt     string           `json:"updatedAt"`
	Owner         RepositoryOwner  `json:"owner"`
}

// RepositoryOwner represents the owner of a repository
type RepositoryOwner struct {
	ID        int64  `json:"id"`
	Login     string `json:"login"`
	FullName  string `json:"fullName"`
	Email     string `json:"email"`
	AvatarURL string `json:"avatarUrl"`
}

// RepositoryStats contains statistics about repositories
type RepositoryStats struct {
	TotalCount   int                    `json:"totalCount"`
	PrivateCount int                    `json:"privateCount"`
	PublicCount  int                    `json:"publicCount"`
	Languages    []LanguageDistribution `json:"languages"`
}

// LanguageDistribution represents language distribution in repositories
type LanguageDistribution struct {
	Language string `json:"language"`
	Count    int    `json:"count"`
}

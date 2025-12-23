package models

// User represents a user (fetched from LDAP Manager service)
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

// GiteaRepository represents a repository from Gitea
type GiteaRepository struct {
	ID            int64           `json:"id"`
	Name          string          `json:"name"`
	FullName      string          `json:"fullName"`
	Description   string          `json:"description"`
	Private       bool            `json:"private"`
	Fork          bool            `json:"fork"`
	HTMLURL       string          `json:"htmlUrl"`
	SSHURL        string          `json:"sshUrl"`
	CloneURL      string          `json:"cloneUrl"`
	DefaultBranch string          `json:"defaultBranch"`
	Language      string          `json:"language"`
	Stars         int             `json:"stars"`
	Forks         int             `json:"forks"`
	Size          int             `json:"size"`
	CreatedAt     string          `json:"createdAt"`
	UpdatedAt     string          `json:"updatedAt"`
	Owner         RepositoryOwner `json:"owner"`
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

// LanguageDistribution represents language usage statistics
type LanguageDistribution struct {
	Language string `json:"language"`
	Count    int    `json:"count"`
}

// HealthStatus represents the health status of the service
type HealthStatus struct {
	Status      string `json:"status"`
	Timestamp   int64  `json:"timestamp"`
	Gitea       bool   `json:"gitea"`
	LDAPManager bool   `json:"ldapManager"`
}

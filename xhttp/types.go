package xhttp

type Credentials struct {
	ProjectDomain  string
	ProjectName    string
	AuthUrl        string
	Username       string
	Password       string
	UserDomainName string
}

// Config is a structure containing configuration
// options for connecting to a Barbican server.
type Config struct {
	// Endpoint is the Barbican instance endpoint.
	Endpoint string

	// Credentials used to login to OpenStack to retrieve the APIKey
	Login Credentials
}

// Auth request structures

type AuthRequest struct {
	Auth Auth `json:"auth"`
}

type Auth struct {
	Identity AuthIdentity `json:"identity"`
	Scope    Scope        `json:"scope"`
}

type Scope struct {
	Project Project `json:"project"`
}

type Project struct {
	Domain Name   `json:"domain"`
	Name   string `json:"name"`
}

type AuthIdentity struct {
	Methods  []string `json:"methods"`
	Password Password `json:"password"`
}

type Name struct {
	Name string `json:"name"`
}

type Password struct {
	User User `json:"user"`
}

type User struct {
	Domain   Name   `json:"domain"`
	Name     string `json:"name"`
	Password string `json:"password"`
}

// Auth response structures

type AuthResponse struct {
	Token Token
}

type Token struct {
	ExpiresAt string `json:"expires_at"`
}

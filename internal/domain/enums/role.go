package enums

// Role represents user roles in the system
type Role string

const (
	RoleUser  Role = "USER"
	RoleAdmin Role = "ADMIN"
)

// IsValid checks if the role is valid
func (r Role) IsValid() bool {
	switch r {
	case RoleUser, RoleAdmin:
		return true
	}
	return false
}

// String returns the string representation of the role
func (r Role) String() string {
	return string(r)
}

package models

type User struct {
	ID             string  `json:"id"`
	Email          string  `json:"email"`
	EmailValidAt   *string `json:"email_valid_at"`
	Name           string  `json:"name"`
	Password       string  `json:"password"`
	Status         string  `json:"status"`
	RegisterType   string  `json:"register_type"`
	RegisterDetail string  `json:"register_detail"`
	LastSeen       string  `json:"last_seen"`
	CreatedAt      string  `json:"created_at"`
	UpdatedAt      string  `json:"updated_at"`
	DeletedAt      *string `json:"deleted_at"`
}

// UserParamater ...
type UserParamater struct {
	ID           string `json:"id"`
	Email        string `json:"email"`
	Name         string `json:"name"`
	ShowPassword bool   `json:"show_password"`
	Status       string `json:"status"`
	Search       string `json:"search"`
}

var (
	UserStatusPending   = "pending"
	UserStatusActive    = "active"
	UserStatusInactive  = "inactive"
	UserStatusWhitelist = []string{UserStatusPending, UserStatusActive, UserStatusInactive}

	UserRegisterTypeEmail = "email"

	UserSocialMediaWhitelist = []string{
		UserRegisterTypeEmail,
	}

	UserSelectStatement = `
	SELECT 
	def.id, def.name, def.email, def.password, 
	def.status, def.register_type, def.register_detail, def.email_valid_at, def.last_seen,
	def.created_at, def.updated_at
	FROM users def
	`

	// UserOrderBy ...
	UserOrderBy = []string{"def.id", "def.name", "def.email", "def.created_at", "def.update_at"}

	// UserWhereStatement ...
	UserWhereStatement = `WHERE def.deleted_at IS NULL`

	// UserGroupByStatement ...
	UserGroupByStatement = `GROUP BY def.id`
)

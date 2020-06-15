package validatores

// Login Faz a validação do login
type Login struct {
	User     string `form:"user" json:"user" xml:"user"  binding:"required"`
	Password string `form:"password" json:"password" xml:"password" binding:"required"`
}

// Register Faz a validação do Register
type Register struct {
	Username string `validate:"required"`
	Email    string `validate:"required,email"`
	Password string `validate:"required"`
}

// UpdateUser Faz a validação do UpdateUser
type UpdateUser struct {
	Username string `validate:"required"`
}

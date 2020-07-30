package validatores

// Login Faz a validação do login
type Login struct {
	Email    string `form:"email" json:"email" xml:"email"  binding:"required" validate:"required,email"`
	Password string `form:"password" json:"password" xml:"password" binding:"required"`
}

// Register Faz a validação do Register
type Register struct {
	Username string `validate:"required"`
	Email    string `validate:"required,email"`
	Password string `validate:"required"`
}

func (data Register) CreateRegister() (valid bool, listError Errors) {
	return Validate(data)
}

// UpdateUser Faz a validação do UpdateUser
type UpdateUser struct {
	Username string `validate:"required"`
}

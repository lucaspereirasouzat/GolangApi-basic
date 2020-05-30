package user

type User struct {
	ID                        uint64
	Username, Password, Email string
}

//A sample use
var user = User{
	// ID:       1,
	Username: "username",
	Password: "password",
	Email:    "lucas@teste.com",
}

package user

import "time"

/* User é o modelo do usuario padrão */
type User struct {
	ID                        uint64
	Username, Password, Email string
	CreatedAt                 time.Time `db:"create_at"`
}

type Token struct {
	token string
	user  User
}

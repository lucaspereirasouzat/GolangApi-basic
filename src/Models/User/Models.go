package user

import (
	"database/sql"
	"time"
)

/* User é o modelo do usuario padrão */
type User struct {
	ID                        uint64
	FileId                    sql.NullInt64 `db:"file_id"`
	Username, Password, Email string
	Securelevel               string    `db:"securelevel"`
	CreatedAt                 time.Time `db:"created_at"`
}

type Token struct {
	token string
	user  User
}

func (u User) Relations() string {
	return ""
}

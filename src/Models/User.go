package models

import (
	"database/sql"
	"time"
)

/* User é o modelo do usuario padrão */
type User struct {
	ID                        uint64
	Pathfile                  sql.NullString `db:"pathfile"`
	Username, Password, Email string
	Securelevel               string    `db:"securelevel"`
	CreatedAt                 time.Time `db:"created_at"`
}

type Token struct {
	token     string
	isRevoked bool
	UserID    sql.NullInt64 `db:"user_id"`
	CreatedAt time.Time     `db:"created_at"`
}

func (u User) Relations() string {
	return ""
}

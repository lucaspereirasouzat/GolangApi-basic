package file

import (
	"database/sql"
	"time"
)

/* File é o padrão para se salvar os arquivos no sistema */
type File struct {
	ID        uint64
	UserId    sql.NullInt64 `db:"user_id"`
	Path      string
	CreatedAt time.Time `db:"created_at"`
}

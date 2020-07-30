package models

import (
	"database/sql"
	"time"
)

/* Notification é o modelo das notificações */
type Notification struct {
	TokenNotification string        `db:"tokennotification"`
	UserID            sql.NullInt64 `db:"user_id"`
	CreatedAt         time.Time     `db:"created_at"`
	User              User          `db:"User"`
}

/* DataNotification é uma notificação */
type DataNotification struct {
	ID          uint64
	UserID      sql.NullInt64 `db:"user_id"`
	Title, Body string
	CreatedAt   time.Time `db:"created_at"`
}

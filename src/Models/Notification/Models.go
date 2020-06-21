package notification

import (
	"database/sql"
	"time"
)

/* Notification é o modelo das notificações */
type Notification struct {
	TokenNotification string        `db:"tokennotification"`
	UserID            sql.NullInt64 `db:"user_id"`
	CreatedAt         time.Time     `db:"created_at"`
}

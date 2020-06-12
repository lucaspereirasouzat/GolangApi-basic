package seeder

import (
	"github.com/jmoiron/sqlx"
)

/*
	Schema para se fazer as migrations
*/
func Seed(db *sqlx.DB) {

	/**
	make the migration to users;
	DROP DATABASE token;
	*/
	tx := db.MustBegin()
	tx.MustExec("INSERT INTO users (username, email, password,secureLevel) VALUES ($1, $2, $3, $4)", "Jason", "jmoiron@jmoiron.net", "1234", "ADM")

	tx.Commit()
	defer db.Close()
}

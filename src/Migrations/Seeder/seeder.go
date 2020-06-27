package seeder

import (
	"docker.go/src/functions"
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
	Password := functions.GenerateMD5("1234")

	tx.MustExec("INSERT INTO users (username, email, password,secureLevel) VALUES ($1, $2, $3, $4)", "Jason", "jmoiron@jmoiron.net", Password, "adm")

	tx.Commit()
	defer db.Close()
}

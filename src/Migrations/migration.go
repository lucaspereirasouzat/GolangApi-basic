package migration

import (
	"database/sql"
)

// var schema = `
// CREATE TABLE person (
//     first_name text,
//     last_name text,
//     email text
// );

// CREATE TABLE place (
//     country text,
//     city text NULL,
//     telcode integer
// )`

/*
	Schema para se fazer as migrations
*/
func Schema() string {
	return `
	CREATE TABLE IF NOT EXISTS users (
		id SERIAL,
		username text,
		email text,
		password text
	) `

	// return `
	// CREATE TABLE user (
	// 	username text,
	// 	email text,
	// 	password text
	// );

	// CREATE TABLE place (
	// 	country text,
	// 	city text NULL,
	// 	telcode integer
	// )`
}

/*
	Person struct
*/
type Person struct {
	FirstName string `db:"first_name"`
	LastName  string `db:"last_name"`
	Email     string
}

/*
	Place struct
*/
type Place struct {
	Country string
	City    sql.NullString
	TelCode int
}

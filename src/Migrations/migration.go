package migration

/*
	Schema para se fazer as migrations
*/
func Schema() string {

	/**
	DROP DATABASE users;
	DROP DATABASE token;
	*/
	return `

	CREATE TABLE IF NOT EXISTS users (
		id SERIAL UNIQUE,
		username text,
		email text,
		password text,
		created_at DATE NOT NULL DEFAULT CURRENT_DATE
	);

	CREATE TABLE IF NOT EXISTS token (
		token text NOT NULL PRIMARY KEY,
		is_revoked bool DEFAULT FALSE,
		user_id INTEGER REFERENCES users(id),
		created_at DATE NOT NULL DEFAULT CURRENT_DATE
	);
	`

}

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
		email text UNIQUE,
		password text,
		secureLevel text DEFAULT 'user',
		pathfile text,
		created_at TIMESTAMP NOT NULL DEFAULT CURRENT_DATE
		updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_DATE
	);

	CREATE TABLE IF NOT EXISTS token (
		token text NOT NULL PRIMARY KEY,
		is_revoked bool DEFAULT FALSE,
		user_id INTEGER REFERENCES users(id) NOT NULL,
		created_at TIMESTAMP NOT NULL DEFAULT CURRENT_DATE
		updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_DATE
	);

	CREATE TABLE IF NOT EXISTS notification (
		tokenNotification text NOT NULL PRIMARY KEY,
		user_id INTEGER REFERENCES users(id) NOT NULL,
		created_at TIMESTAMP NOT NULL DEFAULT CURRENT_DATE
		updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_DATE
	);

	CREATE TABLE IF NOT EXISTS dataNotification (
		id SERIAL UNIQUE,
		user_id INTEGER REFERENCES users(id) NOT NULL,
		title text,
		body text,
		created_at TIMESTAMP NOT NULL DEFAULT CURRENT_DATE
		updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_DATE
	);

	`

}

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
		secureLevel text DEFAULT "user",
		file_id text,
		created_at DATE NOT NULL DEFAULT CURRENT_DATE
	);

	CREATE TABLE IF NOT EXISTS token (
		token text NOT NULL PRIMARY KEY,
		is_revoked bool DEFAULT FALSE,
		user_id INTEGER REFERENCES users(id) NOT NULL,
		created_at DATE NOT NULL DEFAULT CURRENT_DATE
	);

	CREATE TABLE IF NOT EXISTS file (
		id SERIAL UNIQUE PRIMARY KEY,
		path text,
		user_id INTEGER REFERENCES users(id) NOT NULL,
		created_at DATE NOT NULL DEFAULT CURRENT_DATE
	);

	CREATE TABLE IF NOT EXISTS notification (
		tokenNotification text NOT NULL PRIMARY KEY,
		user_id INTEGER REFERENCES users(id) NOT NULL,
		created_at DATE NOT NULL DEFAULT CURRENT_DATE
	);

	do $$
	begin
		IF EXISTS (SELECT file_id FROM users) THEN 
			ALTER TABLE users ADD FOREIGN KEY (file_id) REFERENCES file(id);
		END IF;
	end;
	$$
	`

}

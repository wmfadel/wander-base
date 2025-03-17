CREATE TABLE IF NOT EXISTS roles (
			id SERIAL PRIMARY KEY,
			name TEXT NOT NULL UNIQUE,
			description Text NOT NULL,
			default_role BOOLEAN NOT NULL
		);
CREATE TABLE IF NOT EXISTS user_roles (
  			user_id INTEGER NOT NULL,
    		role_id INTEGER NOT NULL,
    		FOREIGN KEY (user_id) REFERENCES users(id),
   			FOREIGN KEY (role_id) REFERENCES roles(id),
    		PRIMARY KEY (user_id, role_id)
		);
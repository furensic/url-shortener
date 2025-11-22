CREATE TABLE IF NOT EXISTS user_extension (
    id SERIAL PRIMARY KEY,
    user_id SERIAL NOT NULL,
    display_name TEXT,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

/* add ON DELETE CASCADE */
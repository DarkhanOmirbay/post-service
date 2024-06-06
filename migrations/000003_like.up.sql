CREATE TABLE IF NOT EXISTS likes (
                       id SERIAL PRIMARY KEY,
                       user_id INT NOT NULL,
                       post_id INT NOT NULL,
                       created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                       FOREIGN KEY (user_id) REFERENCES users(id),
                       FOREIGN KEY (post_id) REFERENCES posts(id),
                       UNIQUE (user_id, post_id)
);

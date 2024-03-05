-- Create the users table
CREATE TABLE users (
    id SERIAL PRIMARY KEY, username VARCHAR(100), password VARCHAR(100), api_key VARCHAR(255)
);
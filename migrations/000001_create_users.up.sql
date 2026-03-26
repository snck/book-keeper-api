CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_name VARCHAR(50) UNIQUE,
    password_hash CHAR(60),
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP
);

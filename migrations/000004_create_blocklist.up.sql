CREATE TABLE blocklists (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID REFERENCES users(id),
    token TEXT,
    expired_at TIMESTAMP
)
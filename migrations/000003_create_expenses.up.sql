CREATE TABLE expenses (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    amount INT NOT NULL,
    category_id UUID NOT NULL REFERENCES categories(id),
    note VARCHAR(255),
    user_id UUID NOT NULL REFERENCES users(id),
    expense_date TIMESTAMP NOT NULL,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP
);

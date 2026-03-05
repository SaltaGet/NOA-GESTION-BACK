CREATE TABLE categories (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    delete_at TIMESTAMP WITH TIME ZONE
);
CREATE UNIQUE INDEX idx_categories_name ON categories(name);
CREATE INDEX idx_categories_delete_at ON categories(delete_at);

CREATE TABLE plans (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(255),
    price_mounthly DECIMAL(10,2),
    price_yearly DECIMAL(10,2),
    description TEXT,
    features TEXT,
    amount_point_sale BIGINT NOT NULL,
    amount_member BIGINT NOT NULL,
    amount_product BIGINT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);
CREATE UNIQUE INDEX idx_plans_name ON plans(name);

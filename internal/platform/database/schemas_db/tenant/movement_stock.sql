CREATE TABLE movement_stocks (
    id BIGSERIAL PRIMARY KEY,
    member_id BIGINT NOT NULL,
    product_id BIGINT NOT NULL,
    amount NUMERIC(15, 2) NOT NULL,
    from_id BIGINT NOT NULL,
    from_type VARCHAR(20) NOT NULL,
    to_id BIGINT NOT NULL,
    to_type VARCHAR(20) NOT NULL,
    ignore_stock BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_movement_stocks_member FOREIGN KEY (member_id) REFERENCES members(id),
    CONSTRAINT fk_movement_stocks_product FOREIGN KEY (product_id) REFERENCES products(id)
);

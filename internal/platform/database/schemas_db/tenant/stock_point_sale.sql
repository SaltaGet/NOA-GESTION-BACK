CREATE TABLE stock_point_sales (
    id BIGSERIAL PRIMARY KEY,
    product_id BIGINT NOT NULL,
    point_sale_id BIGINT NOT NULL,
    stock NUMERIC(15, 2) NOT NULL DEFAULT 0,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_stock_point_sales_product FOREIGN KEY (product_id) REFERENCES products(id),
    CONSTRAINT fk_stock_point_sales_point_sale FOREIGN KEY (point_sale_id) REFERENCES point_sales(id)
);

CREATE TABLE income_sales (
    id BIGSERIAL PRIMARY KEY,
    point_sale_id BIGINT NOT NULL,
    member_id BIGINT NOT NULL,
    client_id BIGINT NOT NULL,
    cash_register_id BIGINT NOT NULL,
    subtotal NUMERIC(15, 2) NOT NULL,
    discount NUMERIC(15, 2) NOT NULL DEFAULT 0,
    type VARCHAR(20) NOT NULL DEFAULT 'percent',
    total NUMERIC(15, 2) NOT NULL,
    is_budget BOOLEAN NOT NULL DEFAULT FALSE,
    invoice_id BIGINT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_income_sales_point_sale FOREIGN KEY (point_sale_id) REFERENCES point_sales(id),
    CONSTRAINT fk_income_sales_member FOREIGN KEY (member_id) REFERENCES members(id),
    CONSTRAINT fk_income_sales_client FOREIGN KEY (client_id) REFERENCES clients(id),
    CONSTRAINT fk_income_sales_cash_register FOREIGN KEY (cash_register_id) REFERENCES cash_registers(id),
    CONSTRAINT fk_income_sales_invoice FOREIGN KEY (invoice_id) REFERENCES invoices(id)
);

CREATE TABLE income_sale_items (
    id BIGSERIAL PRIMARY KEY,
    income_sale_id BIGINT,
    product_id BIGINT,
    amount NUMERIC(15, 2) NOT NULL,
    price_cost NUMERIC(15, 2) NOT NULL,
    price NUMERIC(15, 2) NOT NULL,
    discount NUMERIC(15, 2) NOT NULL DEFAULT 0,
    type_discount VARCHAR(20) NOT NULL DEFAULT 'percent',
    subtotal NUMERIC(15, 2) NOT NULL,
    total NUMERIC(15, 2) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_income_sale_items_income_sale FOREIGN KEY (income_sale_id) REFERENCES income_sales(id),
    CONSTRAINT fk_income_sale_items_product FOREIGN KEY (product_id) REFERENCES products(id)
);
CREATE INDEX idx_income_sale_items_income_sale_id ON income_sale_items(income_sale_id);
CREATE INDEX idx_income_sale_items_product_id ON income_sale_items(product_id);

CREATE TABLE income_others (
    id BIGSERIAL PRIMARY KEY,
    point_sale_id BIGINT,
    member_id BIGINT,
    cash_register_id BIGINT,
    total NUMERIC(15, 2) NOT NULL,
    type_income_id BIGINT NOT NULL,
    details VARCHAR(255),
    method_income VARCHAR(30) NOT NULL DEFAULT 'cash',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_income_others_point_sale FOREIGN KEY (point_sale_id) REFERENCES point_sales(id),
    CONSTRAINT fk_income_others_member FOREIGN KEY (member_id) REFERENCES members(id),
    CONSTRAINT fk_income_others_cash_register FOREIGN KEY (cash_register_id) REFERENCES cash_registers(id),
    CONSTRAINT fk_income_others_type_income FOREIGN KEY (type_income_id) REFERENCES type_incomes(id)
);

CREATE TABLE income_ecommerces (
    id BIGSERIAL PRIMARY KEY,
    payment_id VARCHAR(255) NOT NULL,
    external_reference VARCHAR(255) NOT NULL,
    status VARCHAR(255) NOT NULL,
    total NUMERIC(15, 2) NOT NULL,
    delivery_status VARCHAR(255) NOT NULL,
    delivery_id VARCHAR(255) NOT NULL,
    date_created VARCHAR(255) NOT NULL,
    date_approved VARCHAR(255) NOT NULL,
    transaction_amount NUMERIC(15, 2),
    net_received_amount NUMERIC(15, 2) NOT NULL,
    payer_first_name VARCHAR(255) NOT NULL,
    payer_last_name VARCHAR(255) NOT NULL,
    payer_email VARCHAR(255) NOT NULL,
    pay_method VARCHAR(255) NOT NULL,
    operation_type VARCHAR(255) NOT NULL,
    message VARCHAR(255),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);
CREATE INDEX idx_income_ecommerces_payment_id ON income_ecommerces(payment_id);
CREATE UNIQUE INDEX idx_income_ecommerces_external_reference ON income_ecommerces(external_reference);

CREATE TABLE income_ecommerce_items (
    id BIGSERIAL PRIMARY KEY,
    income_ecommerce_id BIGINT,
    product_id BIGINT,
    amount NUMERIC(15, 2) NOT NULL,
    price_cost NUMERIC(15, 2) NOT NULL,
    price NUMERIC(15, 2) NOT NULL,
    discount NUMERIC(15, 2) NOT NULL DEFAULT 0,
    type_discount VARCHAR(20) NOT NULL DEFAULT 'percent',
    subtotal NUMERIC(15, 2) NOT NULL,
    total NUMERIC(15, 2) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_income_ecommerce_items_income_ecommerce FOREIGN KEY (income_ecommerce_id) REFERENCES income_ecommerces(id),
    CONSTRAINT fk_income_ecommerce_items_product FOREIGN KEY (product_id) REFERENCES products(id)
);
CREATE INDEX idx_income_ecommerce_items_income_ecommerce_id ON income_ecommerce_items(income_ecommerce_id);
CREATE INDEX idx_income_ecommerce_items_product_id ON income_ecommerce_items(product_id);

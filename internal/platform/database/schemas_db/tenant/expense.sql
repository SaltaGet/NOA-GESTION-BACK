CREATE TABLE expense_buys (
    id BIGSERIAL PRIMARY KEY,
    member_id BIGINT NOT NULL,
    supplier_id BIGINT NOT NULL,
    details VARCHAR(255),
    subtotal NUMERIC(15, 2),
    discount NUMERIC(15, 2) NOT NULL DEFAULT 0,
    type_discount VARCHAR(20) NOT NULL DEFAULT 'percent',
    total NUMERIC(15, 2),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_expense_buys_member FOREIGN KEY (member_id) REFERENCES members(id),
    CONSTRAINT fk_expense_buys_supplier FOREIGN KEY (supplier_id) REFERENCES suppliers(id)
);

CREATE TABLE expense_buy_items (
    id BIGSERIAL PRIMARY KEY,
    expense_buy_id BIGINT NOT NULL,
    product_id BIGINT NOT NULL,
    amount NUMERIC(15, 2) NOT NULL,
    price NUMERIC(15, 2) NOT NULL,
    discount NUMERIC(15, 2) NOT NULL DEFAULT 0,
    type_discount VARCHAR(20) NOT NULL DEFAULT 'percent',
    subtotal NUMERIC(15, 2) NOT NULL,
    total NUMERIC(15, 2) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_expense_buy_items_expense_buy FOREIGN KEY (expense_buy_id) REFERENCES expense_buys(id),
    CONSTRAINT fk_expense_buy_items_product FOREIGN KEY (product_id) REFERENCES products(id)
);

CREATE TABLE expense_others (
    id BIGSERIAL PRIMARY KEY,
    point_sale_id BIGINT,
    member_id BIGINT NOT NULL,
    cash_register_id BIGINT,
    details VARCHAR(255),
    type_expense_id BIGINT NOT NULL,
    total NUMERIC(15, 2) NOT NULL,
    pay_method VARCHAR(30) DEFAULT 'efectivo',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_expense_others_point_sale FOREIGN KEY (point_sale_id) REFERENCES point_sales(id),
    CONSTRAINT fk_expense_others_member FOREIGN KEY (member_id) REFERENCES members(id),
    CONSTRAINT fk_expense_others_cash_register FOREIGN KEY (cash_register_id) REFERENCES cash_registers(id),
    CONSTRAINT fk_expense_others_type_expense FOREIGN KEY (type_expense_id) REFERENCES type_expenses(id)
);

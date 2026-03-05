CREATE TABLE pay_incomes (
    id BIGSERIAL PRIMARY KEY,
    income_sale_id BIGINT,
    cash_register_id BIGINT,
    client_id BIGINT,
    total NUMERIC(15, 2) NOT NULL,
    method_pay VARCHAR(30) NOT NULL DEFAULT 'cash',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_pay_incomes_income_sale FOREIGN KEY (income_sale_id) REFERENCES income_sales(id),
    CONSTRAINT fk_pay_incomes_cash_register FOREIGN KEY (cash_register_id) REFERENCES cash_registers(id),
    CONSTRAINT fk_pay_incomes_client FOREIGN KEY (client_id) REFERENCES clients(id)
);
CREATE INDEX idx_pay_incomes_income_sale_id ON pay_incomes(income_sale_id);
CREATE INDEX idx_pay_incomes_cash_register_id ON pay_incomes(cash_register_id);
CREATE INDEX idx_pay_incomes_client_id ON pay_incomes(client_id);

CREATE TABLE pay_expense_buys (
    id BIGSERIAL PRIMARY KEY,
    expense_buy_id BIGINT,
    cash_register_id BIGINT,
    total NUMERIC(15, 2) NOT NULL,
    method_pay VARCHAR(30) NOT NULL DEFAULT 'cash',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_pay_expense_buys_expense_buy FOREIGN KEY (expense_buy_id) REFERENCES expense_buys(id),
    CONSTRAINT fk_pay_expense_buys_cash_register FOREIGN KEY (cash_register_id) REFERENCES cash_registers(id)
);
CREATE INDEX idx_pay_expense_buys_expense_buy_id ON pay_expense_buys(expense_buy_id);
CREATE INDEX idx_pay_expense_buys_cash_register_id ON pay_expense_buys(cash_register_id);

CREATE TABLE pay_expense_others (
    id BIGSERIAL PRIMARY KEY,
    expense_other_id BIGINT,
    cash_register_id BIGINT,
    total NUMERIC(15, 2) NOT NULL,
    method_pay VARCHAR(30) NOT NULL DEFAULT 'cash',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_pay_expense_others_expense_other FOREIGN KEY (expense_other_id) REFERENCES expense_others(id),
    CONSTRAINT fk_pay_expense_others_cash_register FOREIGN KEY (cash_register_id) REFERENCES cash_registers(id)
);
CREATE INDEX idx_pay_expense_others_expense_other_id ON pay_expense_others(expense_other_id);
CREATE INDEX idx_pay_expense_others_cash_register_id ON pay_expense_others(cash_register_id);

CREATE TABLE cash_registers (
    id BIGSERIAL PRIMARY KEY,
    point_sale_id BIGINT NOT NULL,
    member_open_id BIGINT NOT NULL,
    open_amount NUMERIC(15, 2) DEFAULT 0,
    hour_open TIMESTAMP WITH TIME ZONE,
    member_close_id BIGINT,
    close_amount NUMERIC(15, 2),
    hour_close TIMESTAMP WITH TIME ZONE,
    is_close BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_cash_registers_point_sale FOREIGN KEY (point_sale_id) REFERENCES point_sales(id),
    CONSTRAINT fk_cash_registers_member_open FOREIGN KEY (member_open_id) REFERENCES members(id),
    CONSTRAINT fk_cash_registers_member_close FOREIGN KEY (member_close_id) REFERENCES members(id)
);

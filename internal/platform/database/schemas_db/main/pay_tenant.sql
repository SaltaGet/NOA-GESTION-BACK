CREATE TABLE pay_tenants (
    id BIGSERIAL PRIMARY KEY,
    tenant_id BIGINT NOT NULL,
    admin_id BIGINT NOT NULL,
    amount_month BIGINT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_pay_tenants_tenant FOREIGN KEY (tenant_id) REFERENCES tenants(id) ON DELETE CASCADE,
    CONSTRAINT fk_pay_tenants_admin FOREIGN KEY (admin_id) REFERENCES admins(id)
);

CREATE TABLE pay_details (
    id BIGSERIAL PRIMARY KEY,
    pay_tenant_id BIGINT NOT NULL,
    pay_id VARCHAR(255),
    amount NUMERIC(15, 2) NOT NULL,
    method_pay VARCHAR(30) NOT NULL DEFAULT 'cash',
    state_pay VARCHAR(30) NOT NULL DEFAULT 'pending',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_pay_details_pay_tenant FOREIGN KEY (pay_tenant_id) REFERENCES pay_tenants(id)
);

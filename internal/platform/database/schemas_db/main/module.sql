CREATE TABLE modules (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(100) UNIQUE NOT NULL,
    price_monthly DECIMAL(10,2),
    price_yearly DECIMAL(10,2),
    description TEXT,
    features TEXT,
    amount_images_per_product INT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP WITH TIME ZONE
);
CREATE INDEX idx_modules_deleted_at ON modules(deleted_at);

CREATE TABLE tenant_modules (
    id BIGSERIAL PRIMARY KEY,
    module_id BIGINT NOT NULL,
    tenant_id BIGINT NOT NULL,
    expiration TIMESTAMP WITH TIME ZONE,
    accepted_terms BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP WITH TIME ZONE,
    CONSTRAINT fk_tenant_modules_module FOREIGN KEY (module_id) REFERENCES modules(id) ON UPDATE CASCADE ON DELETE CASCADE,
    CONSTRAINT fk_tenant_modules_tenant FOREIGN KEY (tenant_id) REFERENCES tenants(id) ON UPDATE CASCADE ON DELETE CASCADE
);
CREATE UNIQUE INDEX idx_tenant_module ON tenant_modules(module_id, tenant_id);
CREATE INDEX idx_tenant_modules_deleted_at ON tenant_modules(deleted_at);

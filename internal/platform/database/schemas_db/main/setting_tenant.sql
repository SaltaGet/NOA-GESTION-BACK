CREATE TABLE setting_tenants (
    id BIGSERIAL PRIMARY KEY,
    tenant_id BIGINT UNIQUE NOT NULL,
    logo VARCHAR(255),
    front_page VARCHAR(255),
    title VARCHAR(255),
    slogan TEXT,
    primary_color VARCHAR(255),
    secondary_color VARCHAR(255),
    phone VARCHAR(255),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_setting_tenants_tenant FOREIGN KEY (tenant_id) REFERENCES tenants(id) ON DELETE CASCADE
);

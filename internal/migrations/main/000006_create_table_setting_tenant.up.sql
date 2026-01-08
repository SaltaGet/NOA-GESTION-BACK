CREATE TABLE IF NOT EXISTS setting_tenants (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    tenant_id BIGINT NOT NULL,
    logo VARCHAR(255),
    front_page VARCHAR(255),
    title VARCHAR(255),
    slogan TEXT,
    primary_color VARCHAR(255),
    secondary_color VARCHAR(255),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    UNIQUE INDEX idx_setting_tenant_id (tenant_id), -- Garantiza el 1 a 1
    CONSTRAINT fk_tenant_setting FOREIGN KEY (tenant_id) 
        REFERENCES tenants(id) ON DELETE CASCADE ON UPDATE CASCADE
);
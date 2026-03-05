CREATE TABLE credentials (
    id BIGSERIAL PRIMARY KEY,
    tenant_id BIGINT NOT NULL,
    access_token_mp VARCHAR(255),
    access_token_test_mp VARCHAR(255),
    social_reason VARCHAR(255),
    business_name VARCHAR(255),
    address VARCHAR(255),
    responsibility_front_iva VARCHAR(255),
    gross_income VARCHAR(255),
    start_activities DATE,
    cuit VARCHAR(255) UNIQUE,
    concept VARCHAR(255),
    arca_certificate TEXT,
    arca_key TEXT,
    token_arca VARCHAR(255),
    sign_arca VARCHAR(255),
    expire_token_arca TIMESTAMP WITH TIME ZONE,
    token_email VARCHAR(255),
    CONSTRAINT fk_credentials_tenant FOREIGN KEY (tenant_id) REFERENCES tenants(id)
);
CREATE UNIQUE INDEX idx_credentials_tenant_id ON credentials(tenant_id);

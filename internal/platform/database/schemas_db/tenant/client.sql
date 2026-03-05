CREATE TABLE clients (
    id BIGSERIAL PRIMARY KEY,
    first_name VARCHAR(30) NOT NULL,
    last_name VARCHAR(30) NOT NULL,
    company_name VARCHAR(255),
    identifier VARCHAR(20) UNIQUE,
    email VARCHAR(100) UNIQUE,
    phone VARCHAR(20),
    address VARCHAR(255),
    responsability_front_iva VARCHAR(255),
    member_create_id BIGINT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP WITH TIME ZONE,
    CONSTRAINT fk_clients_member_create FOREIGN KEY (member_create_id) REFERENCES members(id)
);

CREATE INDEX idx_clients_deleted_at ON clients(deleted_at);

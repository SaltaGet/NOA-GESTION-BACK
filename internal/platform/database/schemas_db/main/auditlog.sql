CREATE TABLE audit_log_admins (
    id BIGSERIAL PRIMARY KEY,
    transaction_id BIGINT,
    admin_id BIGINT NOT NULL,
    method VARCHAR(10) NOT NULL,
    path VARCHAR(255) NOT NULL,
    old_value JSONB,
    new_value JSONB,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_audit_log_admins_admin FOREIGN KEY (admin_id) REFERENCES admins(id)
);
CREATE INDEX idx_audit_log_admins_transaction_id ON audit_log_admins(transaction_id);

CREATE TABLE audit_logs (
    id BIGSERIAL PRIMARY KEY,
    transaction_id BIGINT,
    member_id BIGINT NOT NULL,
    method VARCHAR(10) NOT NULL,
    path VARCHAR(255) NOT NULL,
    old_value JSONB,
    new_value JSONB,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_audit_logs_member FOREIGN KEY (member_id) REFERENCES members(id)
);
CREATE INDEX idx_audit_logs_transaction_id ON audit_logs(transaction_id);

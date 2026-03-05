CREATE TABLE members (
    id BIGSERIAL PRIMARY KEY,
    first_name VARCHAR(30) NOT NULL,
    last_name VARCHAR(30) NOT NULL,
    username VARCHAR(30) UNIQUE NOT NULL,
    email VARCHAR(100) UNIQUE NOT NULL,
    password VARCHAR(255) NOT NULL,
    address VARCHAR(255) DEFAULT NULL,
    phone VARCHAR(20) DEFAULT NULL,
    is_admin BOOLEAN NOT NULL DEFAULT FALSE,
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    role_id BIGINT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP WITH TIME ZONE,
    CONSTRAINT fk_members_role FOREIGN KEY (role_id) REFERENCES roles(id)
);

CREATE INDEX idx_members_deleted_at ON members(deleted_at);

-- Generate junction table for member_point_sales since it's many-to-many
CREATE TABLE member_point_sales (
    member_id BIGINT NOT NULL,
    point_sale_id BIGINT NOT NULL,
    PRIMARY KEY (member_id, point_sale_id),
    CONSTRAINT fk_member_point_sales_member FOREIGN KEY (member_id) REFERENCES members(id) ON UPDATE CASCADE ON DELETE CASCADE,
    CONSTRAINT fk_member_point_sales_point_sale FOREIGN KEY (point_sale_id) REFERENCES point_sales(id) ON UPDATE CASCADE ON DELETE CASCADE
);

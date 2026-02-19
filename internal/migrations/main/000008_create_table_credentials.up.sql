CREATE TABLE IF NOT EXISTS `credentials` (
  `id`                          BIGINT AUTO_INCREMENT,
  `tenant_id`                   BIGINT NOT NULL,
  `access_token_mp`             VARCHAR(255),
  `access_token_test_mp`        VARCHAR(255),
  `social_reason`                VARCHAR(255),
  `responsibility_front_iva`  VARCHAR(255),
  `cuit`                        VARCHAR(255),
  `arca_certificate`            TEXT,
  `arca_key`                    TEXT,
  `arca_certificate_test`       TEXT,
  `arca_key_test`               TEXT,
  
  PRIMARY KEY (`id`),
  UNIQUE INDEX `idx_credentials_tenant_id` (`tenant_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

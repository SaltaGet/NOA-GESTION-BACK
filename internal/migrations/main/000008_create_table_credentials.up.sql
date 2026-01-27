CREATE TABLE IF NOT EXISTS `credentials` (
  `id`                          BIGINT AUTO_INCREMENT,
  `tenant_id`                   BIGINT NOT NULL,
  `access_token_mp`             VARCHAR(255),
  `access_token_test_mp`        VARCHAR(255),
  `social_reason`                VARCHAR(255),
  `responsibility_front_iva`  VARCHAR(255),
  `cuit`                        VARCHAR(255),
  `arca_certificate`            LONGTEXT,
  `arca_key`                    LONGTEXT,
  `arca_certificate_test`       LONGTEXT,
  `arca_key_test`               LONGTEXT,
  
  PRIMARY KEY (`id`),
  UNIQUE INDEX `idx_credentials_tenant_id` (`tenant_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

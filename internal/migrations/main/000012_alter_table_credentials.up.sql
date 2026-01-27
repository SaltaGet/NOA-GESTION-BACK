ALTER TABLE `credentials` DROP COLUMN IF EXISTS `arca_certificate_test`;
ALTER TABLE `credentials` DROP COLUMN IF EXISTS `arca_key_test`;

ALTER TABLE `credentials` 
    ADD COLUMN IF NOT EXISTS `business_name`     VARCHAR(255) AFTER `social_reason`,
    ADD COLUMN IF NOT EXISTS `address`           VARCHAR(255) AFTER `business_name`,
    ADD COLUMN IF NOT EXISTS `gross_income`      VARCHAR(255) AFTER `responsibility_front_iva`,
    ADD COLUMN IF NOT EXISTS `start_activities`  DATE         AFTER `gross_income`,
    ADD COLUMN IF NOT EXISTS `token_arca`        TEXT         AFTER `cuit`,
    ADD COLUMN IF NOT EXISTS `sign_arca`         VARCHAR(255) AFTER `token_arca`,
    ADD COLUMN IF NOT EXISTS `expire_token_arca` DATETIME     AFTER `sign_arca`,
    ADD COLUMN IF NOT EXISTS `token_email`       VARCHAR(255) AFTER `arca_key`;

ALTER TABLE `credentials` 
    ADD UNIQUE INDEX `idx_credentials_cuit` (`cuit`);
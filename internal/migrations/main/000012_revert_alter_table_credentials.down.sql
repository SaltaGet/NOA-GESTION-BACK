ALTER TABLE `credentials` 
    DROP COLUMN IF EXISTS `business_name`,
    DROP COLUMN IF EXISTS `address`,
    DROP COLUMN IF EXISTS `gross_income`,
    DROP COLUMN IF EXISTS `start_activities`,
    
    ADD COLUMN `arca_certificate_test` LONGTEXT AFTER `arca_key`,
    ADD COLUMN `arca_key_test`         LONGTEXT AFTER `arca_certificate_test`;
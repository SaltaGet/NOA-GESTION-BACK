ALTER TABLE `credentials` 
    DROP COLUMN `business_name`,
    DROP COLUMN `address`,
    DROP COLUMN `gross_income`,
    DROP COLUMN `start_activities`,
    
    ADD COLUMN `arca_certificate_test` LONGTEXT AFTER `arca_key`,
    ADD COLUMN `arca_key_test`         LONGTEXT AFTER `arca_certificate_test`;
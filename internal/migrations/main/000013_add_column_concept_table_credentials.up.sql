ALTER TABLE `credentials` 
    ADD COLUMN IF NOT EXISTS `concept` VARCHAR(255) AFTER `expire_token_arca`;
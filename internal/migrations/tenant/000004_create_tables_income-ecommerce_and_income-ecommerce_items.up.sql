CREATE TABLE IF NOT EXISTS `income_ecommerces` (
  `id` BIGINT NOT NULL AUTO_INCREMENT,
  `payment_id` VARCHAR(255) NOT NULL,
  `external_reference` VARCHAR(255) NOT NULL,
  `status` VARCHAR(100) NOT NULL,
  `total` DECIMAL(19,4) NOT NULL,
  `delivery_status` VARCHAR(100) NOT NULL,
  `delivery_id` VARCHAR(255) NULL,
  `date_created` VARCHAR(100) NOT NULL,
  `date_approved` VARCHAR(100) NOT NULL,
  `transaction_amount` DECIMAL(19,4) NULL,
  `net_received_amount` DECIMAL(19,4) NOT NULL,
  `payer_first_name` VARCHAR(255) NOT NULL,
  `payer_last_name` VARCHAR(255) NOT NULL,
  `payer_email` VARCHAR(255) NOT NULL,
  `pay_method` VARCHAR(100) NOT NULL,
  `operation_type` VARCHAR(100) NOT NULL,
  `message` TEXT NULL,
  `created_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
  `updated_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3)
    ON UPDATE CURRENT_TIMESTAMP(3),
  PRIMARY KEY (`id`),
  UNIQUE INDEX `idx_income_ecommerces_external_reference` (`external_reference`),
  INDEX `idx_income_ecommerces_payment_id` (`payment_id`),
  INDEX `idx_income_ecommerces_status` (`status`),
  INDEX `idx_income_ecommerces_created_at` (`created_at`)
) ENGINE=InnoDB
  DEFAULT CHARSET=utf8mb4
  COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS `income_ecommerce_items` (
  `id` BIGINT NOT NULL AUTO_INCREMENT,
  `income_ecommerce_id` BIGINT NOT NULL,
  `product_id` BIGINT NOT NULL,
  `amount` DECIMAL(19,4) NOT NULL,
  `price_cost` DECIMAL(19,4) NOT NULL,
  `price` DECIMAL(19,4) NOT NULL,
  `discount` DECIMAL(19,4) NOT NULL DEFAULT 0,
  `type_discount` VARCHAR(50) NOT NULL DEFAULT 'percent',
  `subtotal` DECIMAL(19,4) NOT NULL,
  `total` DECIMAL(19,4) NOT NULL,
  `created_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
  PRIMARY KEY (`id`),
  INDEX `idx_income_ecommerce_items_income_ecommerce_id` (`income_ecommerce_id`),
  INDEX `idx_income_ecommerce_items_product_id` (`product_id`),
  CONSTRAINT `fk_income_ecommerce` 
    FOREIGN KEY (`income_ecommerce_id`) 
    REFERENCES `income_ecommerces` (`id`) 
    ON DELETE CASCADE
    ON UPDATE NO ACTION,
  CONSTRAINT `fk_product` 
    FOREIGN KEY (`product_id`) 
    REFERENCES `products` (`id`)
    ON DELETE RESTRICT
    ON UPDATE NO ACTION
) ENGINE=InnoDB
  DEFAULT CHARSET=utf8mb4
  COLLATE=utf8mb4_unicode_ci;
/*M!999999\- enable the sandbox mode */ 
-- MariaDB dump 10.19  Distrib 10.11.13-MariaDB, for debian-linux-gnu (x86_64)
--
-- Host: 127.0.0.1    Database: string_tenant2
-- ------------------------------------------------------
-- Server version	10.11.13-MariaDB-0ubuntu0.24.04.1-log

/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/*!40101 SET @OLD_CHARACTER_SET_RESULTS=@@CHARACTER_SET_RESULTS */;
/*!40101 SET @OLD_COLLATION_CONNECTION=@@COLLATION_CONNECTION */;
/*!40101 SET NAMES utf8mb4 */;
/*!40103 SET @OLD_TIME_ZONE=@@TIME_ZONE */;
/*!40103 SET TIME_ZONE='+00:00' */;
/*!40014 SET @OLD_UNIQUE_CHECKS=@@UNIQUE_CHECKS, UNIQUE_CHECKS=0 */;
/*!40014 SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0 */;
/*!40101 SET @OLD_SQL_MODE=@@SQL_MODE, SQL_MODE='NO_AUTO_VALUE_ON_ZERO' */;
/*!40111 SET @OLD_SQL_NOTES=@@SQL_NOTES, SQL_NOTES=0 */;

--
-- Current Database: `string_tenant2`
--

CREATE DATABASE /*!32312 IF NOT EXISTS*/ `string_tenant2` /*!40100 DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci */;

USE `string_tenant2`;

--
-- Table structure for table `admins`
--

DROP TABLE IF EXISTS `admins`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8mb4 */;
CREATE TABLE `admins` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `first_name` varchar(30) NOT NULL,
  `last_name` varchar(30) NOT NULL,
  `username` varchar(30) NOT NULL,
  `email` varchar(191) NOT NULL,
  `password` longtext NOT NULL,
  `is_super_admin` tinyint(1) NOT NULL DEFAULT 0,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `deleted_at` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_admins_username` (`username`),
  UNIQUE KEY `uni_admins_email` (`email`),
  KEY `idx_admins_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `admins`
--

LOCK TABLES `admins` WRITE;
/*!40000 ALTER TABLE `admins` DISABLE KEYS */;
/*!40000 ALTER TABLE `admins` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `audit_log_admins`
--

DROP TABLE IF EXISTS `audit_log_admins`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8mb4 */;
CREATE TABLE `audit_log_admins` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `admin_id` bigint(20) NOT NULL,
  `method` longtext NOT NULL,
  `path` longtext NOT NULL,
  `old_value` longtext DEFAULT NULL,
  `new_value` longtext DEFAULT NULL,
  `created_at` longtext DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `audit_log_admins`
--

LOCK TABLES `audit_log_admins` WRITE;
/*!40000 ALTER TABLE `audit_log_admins` DISABLE KEYS */;
/*!40000 ALTER TABLE `audit_log_admins` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `audit_logs`
--

DROP TABLE IF EXISTS `audit_logs`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8mb4 */;
CREATE TABLE `audit_logs` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `member_id` bigint(20) NOT NULL,
  `tenant_id` bigint(20) NOT NULL,
  `method` longtext NOT NULL,
  `path` longtext NOT NULL,
  `old_value` longtext DEFAULT NULL,
  `new_value` longtext DEFAULT NULL,
  `created_at` longtext DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `audit_logs`
--

LOCK TABLES `audit_logs` WRITE;
/*!40000 ALTER TABLE `audit_logs` DISABLE KEYS */;
/*!40000 ALTER TABLE `audit_logs` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `cash_registers`
--

DROP TABLE IF EXISTS `cash_registers`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8mb4 */;
CREATE TABLE `cash_registers` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `point_sale_id` bigint(20) NOT NULL,
  `member_open_id` bigint(20) NOT NULL,
  `open_amount` double DEFAULT NULL,
  `hour_open` datetime(3) DEFAULT NULL,
  `member_close_id` bigint(20) DEFAULT NULL,
  `close_amount` double DEFAULT NULL,
  `hour_close` datetime(3) DEFAULT NULL,
  `is_close` tinyint(1) DEFAULT NULL,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `fk_cash_registers_point_sale` (`point_sale_id`),
  KEY `fk_cash_registers_member_open` (`member_open_id`),
  KEY `fk_cash_registers_member_close` (`member_close_id`),
  CONSTRAINT `fk_cash_registers_member_close` FOREIGN KEY (`member_close_id`) REFERENCES `members` (`id`),
  CONSTRAINT `fk_cash_registers_member_open` FOREIGN KEY (`member_open_id`) REFERENCES `members` (`id`),
  CONSTRAINT `fk_cash_registers_point_sale` FOREIGN KEY (`point_sale_id`) REFERENCES `point_sales` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `cash_registers`
--

LOCK TABLES `cash_registers` WRITE;
/*!40000 ALTER TABLE `cash_registers` DISABLE KEYS */;
/*!40000 ALTER TABLE `cash_registers` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `categories`
--

DROP TABLE IF EXISTS `categories`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8mb4 */;
CREATE TABLE `categories` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `name` varchar(100) NOT NULL,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `delete_at` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_categories_name` (`name`),
  KEY `idx_categories_delete_at` (`delete_at`)
) ENGINE=InnoDB AUTO_INCREMENT=2 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `categories`
--

LOCK TABLES `categories` WRITE;
/*!40000 ALTER TABLE `categories` DISABLE KEYS */;
INSERT INTO `categories` VALUES
(1,'sin categoría','2025-12-05 02:04:49.710','2025-12-05 02:04:49.710',NULL);
/*!40000 ALTER TABLE `categories` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `clients`
--

DROP TABLE IF EXISTS `clients`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8mb4 */;
CREATE TABLE `clients` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `first_name` varchar(30) NOT NULL,
  `last_name` varchar(30) NOT NULL,
  `company_name` longtext DEFAULT NULL,
  `identifier` varchar(20) DEFAULT NULL,
  `email` varchar(191) DEFAULT NULL,
  `phone` longtext DEFAULT NULL,
  `address` longtext DEFAULT NULL,
  `member_create_id` bigint(20) NOT NULL,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `deleted_at` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uni_clients_identifier` (`identifier`),
  UNIQUE KEY `uni_clients_email` (`email`),
  KEY `idx_clients_deleted_at` (`deleted_at`),
  KEY `fk_clients_member_create` (`member_create_id`),
  CONSTRAINT `fk_clients_member_create` FOREIGN KEY (`member_create_id`) REFERENCES `members` (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=2 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `clients`
--

LOCK TABLES `clients` WRITE;
/*!40000 ALTER TABLE `clients` DISABLE KEYS */;
INSERT INTO `clients` VALUES
(1,'Consumidor','Final',NULL,NULL,NULL,NULL,NULL,1,'2025-12-05 02:04:49.708','2025-12-05 02:04:49.708',NULL);
/*!40000 ALTER TABLE `clients` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `deposits`
--

DROP TABLE IF EXISTS `deposits`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8mb4 */;
CREATE TABLE `deposits` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `product_id` bigint(20) NOT NULL,
  `stock` double NOT NULL DEFAULT 0,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `fk_products_stock_deposit` (`product_id`),
  CONSTRAINT `fk_products_stock_deposit` FOREIGN KEY (`product_id`) REFERENCES `products` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `deposits`
--

LOCK TABLES `deposits` WRITE;
/*!40000 ALTER TABLE `deposits` DISABLE KEYS */;
/*!40000 ALTER TABLE `deposits` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `expense_buy_items`
--

DROP TABLE IF EXISTS `expense_buy_items`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8mb4 */;
CREATE TABLE `expense_buy_items` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `expense_buy_id` bigint(20) NOT NULL,
  `product_id` bigint(20) NOT NULL,
  `amount` double NOT NULL,
  `price` double NOT NULL,
  `discount` double NOT NULL DEFAULT 0,
  `type_discount` varchar(191) NOT NULL DEFAULT 'percent',
  `subtotal` double NOT NULL,
  `total` double NOT NULL,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `fk_expense_buy_items_product` (`product_id`),
  KEY `fk_expense_buys_expense_buy_item` (`expense_buy_id`),
  CONSTRAINT `fk_expense_buy_items_product` FOREIGN KEY (`product_id`) REFERENCES `products` (`id`),
  CONSTRAINT `fk_expense_buys_expense_buy_item` FOREIGN KEY (`expense_buy_id`) REFERENCES `expense_buys` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `expense_buy_items`
--

LOCK TABLES `expense_buy_items` WRITE;
/*!40000 ALTER TABLE `expense_buy_items` DISABLE KEYS */;
/*!40000 ALTER TABLE `expense_buy_items` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `expense_buys`
--

DROP TABLE IF EXISTS `expense_buys`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8mb4 */;
CREATE TABLE `expense_buys` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `member_id` bigint(20) NOT NULL,
  `supplier_id` bigint(20) NOT NULL,
  `details` varchar(255) DEFAULT NULL,
  `subtotal` double DEFAULT NULL,
  `discount` double NOT NULL DEFAULT 0,
  `type_discount` varchar(191) NOT NULL DEFAULT 'percent',
  `total` double DEFAULT NULL,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `fk_expense_buys_supplier` (`supplier_id`),
  KEY `fk_expense_buys_member` (`member_id`),
  CONSTRAINT `fk_expense_buys_member` FOREIGN KEY (`member_id`) REFERENCES `members` (`id`),
  CONSTRAINT `fk_expense_buys_supplier` FOREIGN KEY (`supplier_id`) REFERENCES `suppliers` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `expense_buys`
--

LOCK TABLES `expense_buys` WRITE;
/*!40000 ALTER TABLE `expense_buys` DISABLE KEYS */;
/*!40000 ALTER TABLE `expense_buys` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `expense_others`
--

DROP TABLE IF EXISTS `expense_others`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8mb4 */;
CREATE TABLE `expense_others` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `point_sale_id` bigint(20) DEFAULT NULL,
  `member_id` bigint(20) NOT NULL,
  `cash_register_id` bigint(20) DEFAULT NULL,
  `details` varchar(255) DEFAULT NULL,
  `type_expense_id` bigint(20) NOT NULL,
  `total` double NOT NULL,
  `pay_method` varchar(30) DEFAULT 'efectivo',
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `fk_expense_others_type_expense` (`type_expense_id`),
  KEY `fk_expense_others_point_sale` (`point_sale_id`),
  KEY `fk_expense_others_member` (`member_id`),
  KEY `fk_expense_others_cash_register` (`cash_register_id`),
  CONSTRAINT `fk_expense_others_cash_register` FOREIGN KEY (`cash_register_id`) REFERENCES `cash_registers` (`id`),
  CONSTRAINT `fk_expense_others_member` FOREIGN KEY (`member_id`) REFERENCES `members` (`id`),
  CONSTRAINT `fk_expense_others_point_sale` FOREIGN KEY (`point_sale_id`) REFERENCES `point_sales` (`id`),
  CONSTRAINT `fk_expense_others_type_expense` FOREIGN KEY (`type_expense_id`) REFERENCES `type_expenses` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `expense_others`
--

LOCK TABLES `expense_others` WRITE;
/*!40000 ALTER TABLE `expense_others` DISABLE KEYS */;
/*!40000 ALTER TABLE `expense_others` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `income_others`
--

DROP TABLE IF EXISTS `income_others`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8mb4 */;
CREATE TABLE `income_others` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `point_sale_id` bigint(20) DEFAULT NULL,
  `member_id` bigint(20) DEFAULT NULL,
  `cash_register_id` bigint(20) DEFAULT NULL,
  `total` double NOT NULL,
  `type_income_id` bigint(20) NOT NULL,
  `details` longtext DEFAULT NULL,
  `method_income` varchar(191) NOT NULL DEFAULT 'cash',
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `fk_income_others_point_sale` (`point_sale_id`),
  KEY `fk_income_others_member` (`member_id`),
  KEY `fk_income_others_cash_register` (`cash_register_id`),
  KEY `fk_income_others_type_income` (`type_income_id`),
  CONSTRAINT `fk_income_others_cash_register` FOREIGN KEY (`cash_register_id`) REFERENCES `cash_registers` (`id`),
  CONSTRAINT `fk_income_others_member` FOREIGN KEY (`member_id`) REFERENCES `members` (`id`),
  CONSTRAINT `fk_income_others_point_sale` FOREIGN KEY (`point_sale_id`) REFERENCES `point_sales` (`id`),
  CONSTRAINT `fk_income_others_type_income` FOREIGN KEY (`type_income_id`) REFERENCES `type_incomes` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `income_others`
--

LOCK TABLES `income_others` WRITE;
/*!40000 ALTER TABLE `income_others` DISABLE KEYS */;
/*!40000 ALTER TABLE `income_others` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `income_sale_items`
--

DROP TABLE IF EXISTS `income_sale_items`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8mb4 */;
CREATE TABLE `income_sale_items` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `income_sale_id` bigint(20) DEFAULT NULL,
  `product_id` bigint(20) DEFAULT NULL,
  `amount` double NOT NULL,
  `price_cost` double NOT NULL,
  `price` double NOT NULL,
  `discount` double NOT NULL DEFAULT 0,
  `type_discount` varchar(191) NOT NULL DEFAULT 'percent',
  `subtotal` double NOT NULL,
  `total` double NOT NULL,
  `created_at` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_income_sale_items_income_sale_id` (`income_sale_id`),
  KEY `idx_income_sale_items_product_id` (`product_id`),
  CONSTRAINT `fk_income_sale_items_product` FOREIGN KEY (`product_id`) REFERENCES `products` (`id`),
  CONSTRAINT `fk_income_sales_items` FOREIGN KEY (`income_sale_id`) REFERENCES `income_sales` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `income_sale_items`
--

LOCK TABLES `income_sale_items` WRITE;
/*!40000 ALTER TABLE `income_sale_items` DISABLE KEYS */;
/*!40000 ALTER TABLE `income_sale_items` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `income_sales`
--

DROP TABLE IF EXISTS `income_sales`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8mb4 */;
CREATE TABLE `income_sales` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `point_sale_id` bigint(20) NOT NULL,
  `member_id` bigint(20) NOT NULL,
  `client_id` bigint(20) NOT NULL,
  `cash_register_id` bigint(20) NOT NULL,
  `subtotal` double NOT NULL,
  `discount` double NOT NULL DEFAULT 0,
  `type` varchar(191) NOT NULL DEFAULT 'percent',
  `total` double NOT NULL,
  `is_budget` tinyint(1) NOT NULL DEFAULT 0,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `fk_income_sales_client` (`client_id`),
  KEY `fk_income_sales_cash_register` (`cash_register_id`),
  KEY `fk_income_sales_point_sale` (`point_sale_id`),
  KEY `fk_income_sales_member` (`member_id`),
  CONSTRAINT `fk_income_sales_cash_register` FOREIGN KEY (`cash_register_id`) REFERENCES `cash_registers` (`id`),
  CONSTRAINT `fk_income_sales_client` FOREIGN KEY (`client_id`) REFERENCES `clients` (`id`),
  CONSTRAINT `fk_income_sales_member` FOREIGN KEY (`member_id`) REFERENCES `members` (`id`),
  CONSTRAINT `fk_income_sales_point_sale` FOREIGN KEY (`point_sale_id`) REFERENCES `point_sales` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `income_sales`
--

LOCK TABLES `income_sales` WRITE;
/*!40000 ALTER TABLE `income_sales` DISABLE KEYS */;
/*!40000 ALTER TABLE `income_sales` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `member_point_sales`
--

DROP TABLE IF EXISTS `member_point_sales`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8mb4 */;
CREATE TABLE `member_point_sales` (
  `point_sale_id` bigint(20) NOT NULL,
  `member_id` bigint(20) NOT NULL,
  PRIMARY KEY (`point_sale_id`,`member_id`),
  KEY `fk_member_point_sales_member` (`member_id`),
  CONSTRAINT `fk_member_point_sales_member` FOREIGN KEY (`member_id`) REFERENCES `members` (`id`) ON DELETE CASCADE ON UPDATE CASCADE,
  CONSTRAINT `fk_member_point_sales_point_sale` FOREIGN KEY (`point_sale_id`) REFERENCES `point_sales` (`id`) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `member_point_sales`
--

LOCK TABLES `member_point_sales` WRITE;
/*!40000 ALTER TABLE `member_point_sales` DISABLE KEYS */;
/*!40000 ALTER TABLE `member_point_sales` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `members`
--

DROP TABLE IF EXISTS `members`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8mb4 */;
CREATE TABLE `members` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `first_name` varchar(30) NOT NULL,
  `last_name` varchar(30) NOT NULL,
  `username` varchar(30) NOT NULL,
  `email` varchar(191) NOT NULL,
  `password` longtext NOT NULL,
  `address` varchar(255) DEFAULT NULL,
  `phone` varchar(20) DEFAULT NULL,
  `is_admin` tinyint(1) NOT NULL DEFAULT 0,
  `is_active` tinyint(1) NOT NULL DEFAULT 1,
  `role_id` bigint(20) NOT NULL,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `deleted_at` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uni_members_username` (`username`),
  UNIQUE KEY `uni_members_email` (`email`),
  KEY `idx_members_deleted_at` (`deleted_at`),
  KEY `fk_members_role` (`role_id`),
  CONSTRAINT `fk_members_role` FOREIGN KEY (`role_id`) REFERENCES `roles` (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=2 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `members`
--

LOCK TABLES `members` WRITE;
/*!40000 ALTER TABLE `members` DISABLE KEYS */;
INSERT INTO `members` VALUES
(1,'string','string','user2','user2@gmail.com','$argon2id$v=19$m=65536,t=1,p=16$P8Mfs/2EV1iFc0ACV0m92A$QIvcorDIq3GXPmD6Bzhj9JMYYBM0LGzb0woUcEe19Uw','Sin dirección',NULL,1,1,1,'2025-12-05 02:04:49.699','2025-12-05 02:04:49.699',NULL);
/*!40000 ALTER TABLE `members` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `movement_stocks`
--

DROP TABLE IF EXISTS `movement_stocks`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8mb4 */;
CREATE TABLE `movement_stocks` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `member_id` bigint(20) NOT NULL,
  `product_id` bigint(20) NOT NULL,
  `amount` double NOT NULL,
  `from_id` bigint(20) NOT NULL,
  `from_type` longtext NOT NULL,
  `to_id` bigint(20) NOT NULL,
  `to_type` longtext NOT NULL,
  `ignore_stock` tinyint(1) NOT NULL DEFAULT 0,
  `created_at` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `fk_movement_stocks_product` (`product_id`),
  KEY `fk_movement_stocks_member` (`member_id`),
  CONSTRAINT `fk_movement_stocks_member` FOREIGN KEY (`member_id`) REFERENCES `members` (`id`),
  CONSTRAINT `fk_movement_stocks_product` FOREIGN KEY (`product_id`) REFERENCES `products` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `movement_stocks`
--

LOCK TABLES `movement_stocks` WRITE;
/*!40000 ALTER TABLE `movement_stocks` DISABLE KEYS */;
/*!40000 ALTER TABLE `movement_stocks` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `pay_expense_buys`
--

DROP TABLE IF EXISTS `pay_expense_buys`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8mb4 */;
CREATE TABLE `pay_expense_buys` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `expense_buy_id` bigint(20) DEFAULT NULL,
  `cash_register_id` bigint(20) DEFAULT NULL,
  `total` double NOT NULL,
  `method_pay` varchar(191) NOT NULL DEFAULT 'cash',
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_pay_expense_buys_expense_buy_id` (`expense_buy_id`),
  KEY `idx_pay_expense_buys_cash_register_id` (`cash_register_id`),
  CONSTRAINT `fk_expense_buys_pay_expense_buy` FOREIGN KEY (`expense_buy_id`) REFERENCES `expense_buys` (`id`),
  CONSTRAINT `fk_pay_expense_buys_cash_register` FOREIGN KEY (`cash_register_id`) REFERENCES `cash_registers` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `pay_expense_buys`
--

LOCK TABLES `pay_expense_buys` WRITE;
/*!40000 ALTER TABLE `pay_expense_buys` DISABLE KEYS */;
/*!40000 ALTER TABLE `pay_expense_buys` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `pay_expense_others`
--

DROP TABLE IF EXISTS `pay_expense_others`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8mb4 */;
CREATE TABLE `pay_expense_others` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `expense_other_id` bigint(20) DEFAULT NULL,
  `cash_register_id` bigint(20) DEFAULT NULL,
  `total` double NOT NULL,
  `method_pay` varchar(191) NOT NULL DEFAULT 'cash',
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_pay_expense_others_expense_other_id` (`expense_other_id`),
  KEY `idx_pay_expense_others_cash_register_id` (`cash_register_id`),
  CONSTRAINT `fk_pay_expense_others_cash_register` FOREIGN KEY (`cash_register_id`) REFERENCES `cash_registers` (`id`),
  CONSTRAINT `fk_pay_expense_others_expense_other` FOREIGN KEY (`expense_other_id`) REFERENCES `expense_others` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `pay_expense_others`
--

LOCK TABLES `pay_expense_others` WRITE;
/*!40000 ALTER TABLE `pay_expense_others` DISABLE KEYS */;
/*!40000 ALTER TABLE `pay_expense_others` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `pay_incomes`
--

DROP TABLE IF EXISTS `pay_incomes`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8mb4 */;
CREATE TABLE `pay_incomes` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `income_sale_id` bigint(20) DEFAULT NULL,
  `cash_register_id` bigint(20) DEFAULT NULL,
  `client_id` bigint(20) DEFAULT NULL,
  `total` double NOT NULL,
  `method_pay` varchar(191) NOT NULL DEFAULT 'cash',
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_pay_incomes_income_sale_id` (`income_sale_id`),
  KEY `idx_pay_incomes_cash_register_id` (`cash_register_id`),
  KEY `idx_pay_incomes_client_id` (`client_id`),
  CONSTRAINT `fk_clients_pay` FOREIGN KEY (`client_id`) REFERENCES `clients` (`id`),
  CONSTRAINT `fk_income_sales_pay` FOREIGN KEY (`income_sale_id`) REFERENCES `income_sales` (`id`),
  CONSTRAINT `fk_pay_incomes_cash_register` FOREIGN KEY (`cash_register_id`) REFERENCES `cash_registers` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `pay_incomes`
--

LOCK TABLES `pay_incomes` WRITE;
/*!40000 ALTER TABLE `pay_incomes` DISABLE KEYS */;
/*!40000 ALTER TABLE `pay_incomes` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `permissions`
--

DROP TABLE IF EXISTS `permissions`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8mb4 */;
CREATE TABLE `permissions` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `code` longtext DEFAULT NULL,
  `name` longtext NOT NULL,
  `details` longtext NOT NULL,
  `group` longtext NOT NULL,
  `environment` longtext NOT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=92 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `permissions`
--

LOCK TABLES `permissions` WRITE;
/*!40000 ALTER TABLE `permissions` DISABLE KEYS */;
INSERT INTO `permissions` VALUES
(1,'CR01','verificar apertura','Verifica si existe alguna apertura de caja del punto de venta','caja','point_sale'),
(2,'CR02','apertura','Crea una nueva apertura de caja','caja','point_sale'),
(3,'CR03','informes','Obtener infomes de caja','caja','point_sale'),
(4,'CR04','cierre','Cerrar una caja abierta','caja','point_sale'),
(5,'CR05','informe específico','Obtener informe de una caja específica','caja','point_sale'),
(6,'CAT01','crear','Crear una nueva categoría','categoría','point_sale'),
(7,'CAT02','actualizar','Actualizar una categoria existente','categoría','point_sale'),
(8,'CAT03','eliminar','Eliminar una cagotería','categoría','point_sale'),
(9,'CAT04','obtner uno en específico','Obtener una categoría en específico','categoría','point_sale'),
(10,'CAT05','obtner todos','Obtener todas las categorías','categoría','point_sale'),
(11,'CL01','crear','Crear un nuevo cliente','cliente','point_sale'),
(12,'CL02','actualizar','Actualizar un cliente existente','cliente','point_sale'),
(13,'CL03','eliminar','Eliminar un cliente','cliente','point_sale'),
(14,'CL04','obtener uno en específico','Obtener un cliente en específico','cliente','point_sale'),
(15,'CL05','obtener por filtro','Obtener clientes por filtros','cliente','point_sale'),
(16,'CL06','obtener todos','Obtener todos lo clientes','cliente','point_sale'),
(17,'DEP01','obtener uno en específico','Obtener información de un productos en específico del depósito','deposito','point_sale'),
(18,'DEP02','obtener por nombre de producto','Obtener información de un producto en especifico del depósito por nombre','deposito','point_sale'),
(19,'DEP03','obtener por código de producto','Obtener información de un producto en especifico del depósito por codigo','deposito','point_sale'),
(20,'DEP04','obtener todos','Obtener todos los productos del depósito','deposito','point_sale'),
(21,'DEP05','actualizar stock de producto','Actualizar el stock de productos del depósito','deposito','point_sale'),
(22,'EB01','crear','Crear un nuevo gasto de compra','gastos de compra','dashboard'),
(23,'EB02','actualizar','Actualizar un gasto de compra existente','gastos de compra','dashboard'),
(24,'EB03','eliminar','Eliminar un gasto de compra','gastos de compra','dashboard'),
(25,'EB04','obtener uno en específico','Obtener un gasto de compra en específico','gastos de compra','dashboard'),
(26,'EB05','obtener por fecha','Obtener gastos de compra por rango de fechas','gastos de compra','dashboard'),
(27,'EO01','crear','Crear un nuevo otros gastos','otros gastos','dashboard'),
(28,'EO02','actualizar','Actualizar un otros gastos existente','otros gastos','dashboard'),
(29,'EO03','eliminar','Eliminar un otros gastos','otros gastos','dashboard'),
(30,'EO04','obtener uno en específico','Obtener un otros gastos en específico','otros gastos','dashboard'),
(31,'EO05','obtener por fecha','Obtener otros gastos por rango de fechas','otros gastos','dashboard'),
(32,'EOPS01','crear','Crear otro gasto para un punto de venta','otros grastos - punto de venta','point_sale'),
(33,'EOPS02','actualizar','Actualizar gasto de un punto de venta','otros grastos - punto de venta','point_sale'),
(34,'EOPS03','eliminar','Eliminar un gasto de un punto de venta','otros grastos - punto de venta','point_sale'),
(35,'EOPS04','obtener uno en específico','Obtener un gasto en específico de un punto de venta','otros grastos - punto de venta','point_sale'),
(36,'EOPS05','obtener por fecha','Obtener gastos de compra por rango de fechas de un punto de venta','otros grastos - punto de venta','point_sale'),
(37,'INO01','crear','Crear otro ingreso','otros ingresos','dashboard'),
(38,'INO02','actualizar','Actualizar un ingreso existente','otros ingresos','dashboard'),
(39,'INO03','eliminar','Eliminar un ingreso','otros ingresos','dashboard'),
(40,'INO04','obtener uno en específico','Obtener un ingreso en específico','otros ingresos','dashboard'),
(41,'INO05','obtener por fecha','Obtener ingresos por rango de fechas','otros ingresos','dashboard'),
(42,'INOPS01','crear','Crear otro ingreso para un punto de venta','otros ingresos - punto de venta','point_sale'),
(43,'INOPS02','actualizar','Actualizar un ingreso de un punto de venta','otros ingresos - punto de venta','point_sale'),
(44,'INOPS03','eliminar','Eliminar un ingreso de un punto de venta','otros ingresos - punto de venta','point_sale'),
(45,'INOPS04','obtener uno en específico','Obtener un ingreso otros en específico','otros ingresos - punto de venta','point_sale'),
(46,'INOPS05','obtener por fecha','Obtener ingresos otros por rango de fechas','otros ingresos - punto de venta','point_sale'),
(47,'INS01','crear','Crear un nuevo ingreso de venta','ingreso de venta','dashboard'),
(48,'INS02','actualizar','Actualizar un ingreso de venta existente','ingreso de venta','dashboard'),
(49,'INS03','eliminar','Eliminar un ingreso de venta','ingreso de venta','dashboard'),
(50,'INS04','obtener uno en específico','Obtener un ingreso de venta en específico','ingreso de venta','dashboard'),
(51,'INS05','obtener por fecha','Obtener ingresos de venta por rango de fechas','ingreso de venta','dashboard'),
(52,'MB01','crear','Crear un nuevo miembro','miembro','dashboard'),
(53,'MB02','actualizar','Actualizar un miembro existente','miembro','dashboard'),
(54,'MB04','obtener uno en específico','Obtener un miembro en específico','miembro','dashboard'),
(55,'MB05','obtener todos','Obtener todos los miembros','miembro','dashboard'),
(56,'MS01','mover','Mover stock de productos','movimiento de stock','dashboard'),
(57,'MS02','obtener uno en específico','Obtener un movimiento en específico','movimiento de stock','dashboard'),
(58,'MS03','obtener por fecha','Obtener movimientos por rango de fechas','movimiento de stock','dashboard'),
(59,'PER01','obtener todos','Obtener todos los permisos','permisos','dashboard'),
(60,'PER02','obtener propios','Obtener los permisos propios del usuario','permisos','dashboard'),
(61,'PS01','crear','Crear un nuevo punto de venta','punto de venta','dashboard'),
(62,'PS02','actualizar','Actualizar un punto de venta existente','punto de venta','dashboard'),
(63,'PS03','obtener miembros','Obtener los miembros de un punto de venta','punto de venta','dashboard'),
(64,'PS04','obtener todos','Obtener todos los puntos de venta','punto de venta','dashboard'),
(65,'PR01','crear','Crear un nuevo producto','producto','dashboard'),
(66,'PR02','actualizar','Actualizar un producto existente','producto','dashboard'),
(67,'PR03','actualizar lista de precios','Actualizar la lista de precios de productos','producto','dashboard'),
(68,'PR04','eliminar','Eliminar un producto','producto','dashboard'),
(69,'PR05','obtener uno en específico','Obtener un producto en específico','producto','dashboard'),
(70,'PR06','obtener por nombre','Obtener productos por nombre','producto','dashboard'),
(71,'PR07','obtener por código','Obtener producto por codigo','producto','dashboard'),
(72,'PR08','obtener por categoría','Obtener productos por categoría','producto','dashboard'),
(73,'PR09','obtener todos','Obtener todos los productos','producto','dashboard'),
(74,'RP01','obtener excel','Obtener reporte en excel','report','dashboard'),
(75,'RP02','obtener rentabilidad de productos','Obtener rentabilidad de productos','report','dashboard'),
(76,'RP03','obtener por fecha','Obtener reporte por rango de fechas','report','dashboard'),
(77,'RL01','crear','Crear un nuevo rol','rol','dashboard'),
(78,'RL02','obtener todos','Obtener todos los roles','rol','dashboard'),
(79,'ST01','crear','Crear nuevo stock de producto','stock','point_sale'),
(80,'ST02','obtener uno en específico','Obtener stock de un producto en específico','stock','point_sale'),
(81,'ST03','obtener por código','Obtener stock por codigo de producto','stock','point_sale'),
(82,'ST04','obtener por categoría','Obtener stock por categoría de producto','stock','point_sale'),
(83,'ST05','obtener todos','Obtener todos los stocks de productos','stock','point_sale'),
(84,'SP01','crear','Crear un nuevo proveedor','proveedor','dashboard'),
(85,'SP02','actualizar','Actualizar un proveedor existente','proveedor','dashboard'),
(86,'SP03','eliminar','Eliminar un proveedor','proveedor','dashboard'),
(87,'SP04','obtener uno en específico','Obtener un proveedor en específico','proveedor','dashboard'),
(88,'SP05','obtener todos','Obtener todos los proveedores','proveedor','dashboard'),
(89,'TM01','crear','Crear un nuevo tipo de movimiento','tipo de movimiento','dashboard'),
(90,'TM02','actualizar','Actualizar un tipo de movimiento existente','tipo de movimiento','dashboard'),
(91,'TM03','obtener todos','Obtener todos los tipos de movimiento','tipo de movimiento','dashboard');
/*!40000 ALTER TABLE `permissions` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `point_sales`
--

DROP TABLE IF EXISTS `point_sales`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8mb4 */;
CREATE TABLE `point_sales` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `name` varchar(100) NOT NULL,
  `description` varchar(200) DEFAULT NULL,
  `is_deposit` tinyint(1) NOT NULL DEFAULT 0,
  `is_main` tinyint(1) NOT NULL DEFAULT 0,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `delete_at` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_point_sales_name` (`name`),
  KEY `idx_point_sales_delete_at` (`delete_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `point_sales`
--

LOCK TABLES `point_sales` WRITE;
/*!40000 ALTER TABLE `point_sales` DISABLE KEYS */;
/*!40000 ALTER TABLE `point_sales` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `products`
--

DROP TABLE IF EXISTS `products`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8mb4 */;
CREATE TABLE `products` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `code` varchar(50) NOT NULL,
  `name` varchar(100) NOT NULL,
  `description` varchar(200) DEFAULT NULL,
  `price` double DEFAULT NULL,
  `category_id` bigint(20) NOT NULL,
  `notifier` tinyint(1) NOT NULL DEFAULT 0,
  `min_amount` double NOT NULL DEFAULT 0,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_products_code` (`code`),
  KEY `fk_products_category` (`category_id`),
  CONSTRAINT `fk_products_category` FOREIGN KEY (`category_id`) REFERENCES `categories` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `products`
--

LOCK TABLES `products` WRITE;
/*!40000 ALTER TABLE `products` DISABLE KEYS */;
/*!40000 ALTER TABLE `products` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `role_permissions`
--

DROP TABLE IF EXISTS `role_permissions`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8mb4 */;
CREATE TABLE `role_permissions` (
  `role_id` bigint(20) NOT NULL,
  `permission_id` bigint(20) NOT NULL,
  PRIMARY KEY (`role_id`,`permission_id`),
  KEY `fk_role_permissions_permission` (`permission_id`),
  CONSTRAINT `fk_role_permissions_permission` FOREIGN KEY (`permission_id`) REFERENCES `permissions` (`id`) ON DELETE CASCADE ON UPDATE CASCADE,
  CONSTRAINT `fk_role_permissions_role` FOREIGN KEY (`role_id`) REFERENCES `roles` (`id`) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `role_permissions`
--

LOCK TABLES `role_permissions` WRITE;
/*!40000 ALTER TABLE `role_permissions` DISABLE KEYS */;
/*!40000 ALTER TABLE `role_permissions` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `roles`
--

DROP TABLE IF EXISTS `roles`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8mb4 */;
CREATE TABLE `roles` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `name` varchar(191) NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uni_roles_name` (`name`)
) ENGINE=InnoDB AUTO_INCREMENT=2 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `roles`
--

LOCK TABLES `roles` WRITE;
/*!40000 ALTER TABLE `roles` DISABLE KEYS */;
INSERT INTO `roles` VALUES
(1,'admin');
/*!40000 ALTER TABLE `roles` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `schema_migrations`
--

DROP TABLE IF EXISTS `schema_migrations`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8mb4 */;
CREATE TABLE `schema_migrations` (
  `version` bigint(20) NOT NULL,
  `dirty` tinyint(1) NOT NULL,
  PRIMARY KEY (`version`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `schema_migrations`
--

LOCK TABLES `schema_migrations` WRITE;
/*!40000 ALTER TABLE `schema_migrations` DISABLE KEYS */;
INSERT INTO `schema_migrations` VALUES
(1,0);
/*!40000 ALTER TABLE `schema_migrations` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `stock_point_sales`
--

DROP TABLE IF EXISTS `stock_point_sales`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8mb4 */;
CREATE TABLE `stock_point_sales` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `product_id` bigint(20) NOT NULL,
  `point_sale_id` bigint(20) NOT NULL,
  `stock` double NOT NULL DEFAULT 0,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `fk_stock_point_sales_point_sale` (`point_sale_id`),
  KEY `fk_products_stock_point_sales` (`product_id`),
  CONSTRAINT `fk_products_stock_point_sales` FOREIGN KEY (`product_id`) REFERENCES `products` (`id`),
  CONSTRAINT `fk_stock_point_sales_point_sale` FOREIGN KEY (`point_sale_id`) REFERENCES `point_sales` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `stock_point_sales`
--

LOCK TABLES `stock_point_sales` WRITE;
/*!40000 ALTER TABLE `stock_point_sales` DISABLE KEYS */;
/*!40000 ALTER TABLE `stock_point_sales` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `suppliers`
--

DROP TABLE IF EXISTS `suppliers`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8mb4 */;
CREATE TABLE `suppliers` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `name` longtext NOT NULL,
  `company_name` varchar(100) NOT NULL,
  `identifier` varchar(20) DEFAULT NULL,
  `address` varchar(150) DEFAULT NULL,
  `debt_limit` double DEFAULT NULL,
  `email` varchar(191) DEFAULT NULL,
  `phone` varchar(20) DEFAULT NULL,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `deleted_at` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uni_suppliers_identifier` (`identifier`),
  UNIQUE KEY `uni_suppliers_email` (`email`),
  KEY `idx_suppliers_deleted_at` (`deleted_at`)
) ENGINE=InnoDB AUTO_INCREMENT=2 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `suppliers`
--

LOCK TABLES `suppliers` WRITE;
/*!40000 ALTER TABLE `suppliers` DISABLE KEYS */;
INSERT INTO `suppliers` VALUES
(1,'Sin proveedor','Sin nombre',NULL,NULL,NULL,NULL,NULL,'2025-12-05 02:04:49.712','2025-12-05 02:04:49.712',NULL);
/*!40000 ALTER TABLE `suppliers` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `type_expenses`
--

DROP TABLE IF EXISTS `type_expenses`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8mb4 */;
CREATE TABLE `type_expenses` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `name` varchar(191) NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uni_type_expenses_name` (`name`)
) ENGINE=InnoDB AUTO_INCREMENT=2 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `type_expenses`
--

LOCK TABLES `type_expenses` WRITE;
/*!40000 ALTER TABLE `type_expenses` DISABLE KEYS */;
INSERT INTO `type_expenses` VALUES
(1,'Otros');
/*!40000 ALTER TABLE `type_expenses` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `type_incomes`
--

DROP TABLE IF EXISTS `type_incomes`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8mb4 */;
CREATE TABLE `type_incomes` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `name` varchar(191) NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uni_type_incomes_name` (`name`)
) ENGINE=InnoDB AUTO_INCREMENT=2 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `type_incomes`
--

LOCK TABLES `type_incomes` WRITE;
/*!40000 ALTER TABLE `type_incomes` DISABLE KEYS */;
INSERT INTO `type_incomes` VALUES
(1,'Otros');
/*!40000 ALTER TABLE `type_incomes` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Dumping events for database 'string_tenant2'
--

--
-- Dumping routines for database 'string_tenant2'
--
/*!40103 SET TIME_ZONE=@OLD_TIME_ZONE */;

/*!40101 SET SQL_MODE=@OLD_SQL_MODE */;
/*!40014 SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS */;
/*!40014 SET UNIQUE_CHECKS=@OLD_UNIQUE_CHECKS */;
/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
/*!40111 SET SQL_NOTES=@OLD_SQL_NOTES */;

-- Dump completed on 2025-12-05  3:46:34

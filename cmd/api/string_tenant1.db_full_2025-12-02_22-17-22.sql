/*M!999999\- enable the sandbox mode */ 
-- MariaDB dump 10.19  Distrib 10.11.13-MariaDB, for debian-linux-gnu (x86_64)
--
-- Host: 127.0.0.1    Database: string_tenant1.db
-- ------------------------------------------------------
-- Server version	10.11.13-MariaDB-0ubuntu0.24.04.1

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
-- Current Database: `string_tenant1.db`
--

CREATE DATABASE /*!32312 IF NOT EXISTS*/ `string_tenant1.db` /*!40100 DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci */;

USE `string_tenant1.db`;

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
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_admins_username` (`username`),
  UNIQUE KEY `uni_admins_email` (`email`)
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
) ENGINE=InnoDB AUTO_INCREMENT=2 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `cash_registers`
--

LOCK TABLES `cash_registers` WRITE;
/*!40000 ALTER TABLE `cash_registers` DISABLE KEYS */;
INSERT INTO `cash_registers` VALUES
(1,1,1,100,'2025-11-28 17:42:14.944',NULL,NULL,NULL,0,'2025-11-28 17:42:14.945','2025-11-28 17:42:14.945');
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
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_categories_name` (`name`)
) ENGINE=InnoDB AUTO_INCREMENT=3 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `categories`
--

LOCK TABLES `categories` WRITE;
/*!40000 ALTER TABLE `categories` DISABLE KEYS */;
INSERT INTO `categories` VALUES
(1,'categoria1','2025-11-15 23:58:02.765','2025-11-15 23:58:02.765');
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
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `clients`
--

LOCK TABLES `clients` WRITE;
/*!40000 ALTER TABLE `clients` DISABLE KEYS */;
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
) ENGINE=InnoDB AUTO_INCREMENT=4 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `deposits`
--

LOCK TABLES `deposits` WRITE;
/*!40000 ALTER TABLE `deposits` DISABLE KEYS */;
INSERT INTO `deposits` VALUES
(3,2,10,'2025-11-18 01:17:54.133','2025-11-28 17:40:38.154');
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
) ENGINE=InnoDB AUTO_INCREMENT=2 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `expense_buy_items`
--

LOCK TABLES `expense_buy_items` WRITE;
/*!40000 ALTER TABLE `expense_buy_items` DISABLE KEYS */;
INSERT INTO `expense_buy_items` VALUES
(1,1,2,10,200,0,'amount',2000,2000,'2025-11-28 17:40:38.157','2025-11-28 17:40:38.157');
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
  KEY `fk_expense_buys_member` (`member_id`),
  KEY `fk_expense_buys_supplier` (`supplier_id`),
  CONSTRAINT `fk_expense_buys_member` FOREIGN KEY (`member_id`) REFERENCES `members` (`id`),
  CONSTRAINT `fk_expense_buys_supplier` FOREIGN KEY (`supplier_id`) REFERENCES `suppliers` (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=2 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `expense_buys`
--

LOCK TABLES `expense_buys` WRITE;
/*!40000 ALTER TABLE `expense_buys` DISABLE KEYS */;
INSERT INTO `expense_buys` VALUES
(1,1,1,'string',2000,0,'amount',2000,'2025-11-28 17:40:38.156','2025-11-28 17:40:38.156');
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
  KEY `fk_expense_others_member` (`member_id`),
  KEY `fk_expense_others_cash_register` (`cash_register_id`),
  KEY `fk_expense_others_type_expense` (`type_expense_id`),
  KEY `fk_expense_others_point_sale` (`point_sale_id`),
  CONSTRAINT `fk_expense_others_cash_register` FOREIGN KEY (`cash_register_id`) REFERENCES `cash_registers` (`id`),
  CONSTRAINT `fk_expense_others_member` FOREIGN KEY (`member_id`) REFERENCES `members` (`id`),
  CONSTRAINT `fk_expense_others_point_sale` FOREIGN KEY (`point_sale_id`) REFERENCES `point_sales` (`id`),
  CONSTRAINT `fk_expense_others_type_expense` FOREIGN KEY (`type_expense_id`) REFERENCES `type_expenses` (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=3 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `expense_others`
--

LOCK TABLES `expense_others` WRITE;
/*!40000 ALTER TABLE `expense_others` DISABLE KEYS */;
INSERT INTO `expense_others` VALUES
(1,NULL,1,NULL,'la bamos a pasar bomba',1,1000,'cash','2025-11-28 16:03:25.256','2025-11-28 16:03:25.256'),
(2,NULL,1,NULL,'string',1,1000,'cash','2025-12-02 21:54:32.593','2025-12-02 21:54:32.593');
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
) ENGINE=InnoDB AUTO_INCREMENT=3 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `income_others`
--

LOCK TABLES `income_others` WRITE;
/*!40000 ALTER TABLE `income_others` DISABLE KEYS */;
INSERT INTO `income_others` VALUES
(1,NULL,1,NULL,2000,1,'string','cash','2025-11-28 17:44:27.787','2025-11-28 17:44:27.787'),
(2,1,1,1,5000,1,'string','cash','2025-11-28 17:44:58.595','2025-11-28 17:44:58.595');
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
  `created_at` datetime NOT NULL DEFAULT current_timestamp(),
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
  KEY `fk_income_sales_point_sale` (`point_sale_id`),
  KEY `fk_income_sales_member` (`member_id`),
  KEY `fk_income_sales_client` (`client_id`),
  KEY `fk_income_sales_cash_register` (`cash_register_id`),
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
INSERT INTO `member_point_sales` VALUES
(1,1),
(1,2),
(1,3),
(1,4),
(1,5),
(1,6),
(1,7),
(1,8);
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
) ENGINE=InnoDB AUTO_INCREMENT=9 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `members`
--

LOCK TABLES `members` WRITE;
/*!40000 ALTER TABLE `members` DISABLE KEYS */;
INSERT INTO `members` VALUES
(1,'string','string','user1','danielmchachagua@gmail.com','$argon2id$v=19$m=65536,t=1,p=16$/bO+5VMDpT/qYjsnrD0BSw$NVM7R76mBrhDUgNySGUvTD9KEd3Zl8MNJj999EtLvBI','Sin dirección',NULL,1,1,1,'2025-11-15 07:57:46.285','2025-11-15 07:57:46.285',NULL),
(2,'john','doe','user2','user2@b.com','$argon2id$v=19$m=65536,t=1,p=16$OnEBB5t+5CJQpzs4Mc5VCw$X3q7Ms+kzHHJdUHTW0OvVORLuch/OfxX3GkuuYN40d8','Sin dirección','123123123',0,1,2,'2025-11-15 18:42:57.513','2025-11-15 18:42:57.516',NULL),
(3,'john','doe','user3','user3@b.com','$argon2id$v=19$m=65536,t=1,p=16$wbWk9R0Kgz2yMzgtq0tg4w$X/QTYEGdKAKUTh/CETe2KvxeBDhQzDqLei83/K7jyTo','Sin dirección','123123123',0,1,2,'2025-11-15 18:48:15.826','2025-11-15 18:48:15.829',NULL),
(4,'john','doe','user4','user4@b.com','$argon2id$v=19$m=65536,t=1,p=16$DQTTjbTjFyjKmUFkKaUCrg$y2+kdJZHKKTIMo+UZekWYC3nmSvi7UdLd/FZy6NF30Q','Sin dirección','123123123',0,1,2,'2025-11-15 18:49:13.427','2025-11-15 18:49:13.429',NULL),
(5,'john','doe','user5','user5@gmail.com','$argon2id$v=19$m=65536,t=1,p=16$aVvpPIWStXyd+MnEpFeokQ$woBEjWfPZmxJnoXcBFChq4U8ws96HChcHtSFpTEWA7c','Sin dirección','123123123',0,1,2,'2025-11-15 18:55:51.242','2025-11-15 18:55:51.244',NULL),
(6,'john','doe','user6','user6@gmail.com','$argon2id$v=19$m=65536,t=1,p=16$XXCP4qBznEsz5JpLOJ39nA$QSI6ozlNqnpKRcs+Cu48ZdLUr+QQJQge1HQ11zPBaiw','Sin dirección','123123123',0,1,2,'2025-11-15 18:57:41.723','2025-11-15 18:57:41.725',NULL),
(7,'john','doe','user7','user7@gmail.com','$argon2id$v=19$m=65536,t=1,p=16$EfwJgrtmK0/i1LX3RyGPSA$IUZM+THlNVZ7LKXgMaIlF4tgdC9C2ZTV6CzSBln+7mE','Sin dirección','123123123',0,1,2,'2025-11-15 19:05:07.331','2025-11-15 19:05:07.334',NULL),
(8,'john','doe','user8','user8@gmail.com','$argon2id$v=19$m=65536,t=1,p=16$+0Z3eQg1X6gj8EoBdUy0Uw$4QN+hKQkdWIhYcrVPNarfFw0WtAA/7qihNbARopyDBY','Sin dirección','123123123',0,1,2,'2025-11-15 19:08:48.145','2025-11-15 19:08:48.147',NULL);
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
  KEY `fk_movement_stocks_member` (`member_id`),
  KEY `fk_movement_stocks_product` (`product_id`),
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
) ENGINE=InnoDB AUTO_INCREMENT=2 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `pay_expense_buys`
--

LOCK TABLES `pay_expense_buys` WRITE;
/*!40000 ALTER TABLE `pay_expense_buys` DISABLE KEYS */;
INSERT INTO `pay_expense_buys` VALUES
(1,1,NULL,2000,'cash','2025-11-28 17:40:38.158','2025-11-28 17:40:38.158');
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
  `details` longtext NOT NULL,
  `group` longtext NOT NULL,
  `environment` longtext NOT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=6 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `permissions`
--

LOCK TABLES `permissions` WRITE;
/*!40000 ALTER TABLE `permissions` DISABLE KEYS */;
INSERT INTO `permissions` VALUES
(1,'create_client','Crear clientes','clients','dashboard'),
(2,'update_client','Actualizar clientes','clients','dashboard'),
(3,'delete_client','Eliminar clientes','clients','dashboard'),
(4,'create_expense','Crear gastos','expenses','point_sale'),
(5,'update_expense','Actualizar gastos','expenses','point_sale');
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
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_point_sales_name` (`name`)
) ENGINE=InnoDB AUTO_INCREMENT=2 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `point_sales`
--

LOCK TABLES `point_sales` WRITE;
/*!40000 ALTER TABLE `point_sales` DISABLE KEYS */;
INSERT INTO `point_sales` VALUES
(1,'string','string',1,'2025-11-15 18:42:45.898','2025-11-15 18:42:45.900');
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
) ENGINE=InnoDB AUTO_INCREMENT=5 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `products`
--

LOCK TABLES `products` WRITE;
/*!40000 ALTER TABLE `products` DISABLE KEYS */;
INSERT INTO `products` VALUES
(2,'ABC123','Producto1','pepitos',100,1,0,10,'2025-11-15 23:58:08.067','2025-11-16 00:02:59.218'),
(4,'ABC124','Producto1','description',100,1,0,10,'2025-11-29 16:30:22.838','2025-11-29 16:30:22.838');
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
INSERT INTO `role_permissions` VALUES
(2,1),
(2,2),
(2,3),
(2,4);
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
) ENGINE=InnoDB AUTO_INCREMENT=3 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `roles`
--

LOCK TABLES `roles` WRITE;
/*!40000 ALTER TABLE `roles` DISABLE KEYS */;
INSERT INTO `roles` VALUES
(1,'admin'),
(2,'string');
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
(2,0);
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
(1,'string','string','supplier1','string',1000000,'supplier1@gmail.com','string','2025-11-28 17:40:29.389','2025-11-28 17:40:29.389',NULL);
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
(1,'muñeca inflables');
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
(1,'alquilarlo al gaston');
/*!40000 ALTER TABLE `type_incomes` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Dumping events for database 'string_tenant1.db'
--

--
-- Dumping routines for database 'string_tenant1.db'
--
/*!40103 SET TIME_ZONE=@OLD_TIME_ZONE */;

/*!40101 SET SQL_MODE=@OLD_SQL_MODE */;
/*!40014 SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS */;
/*!40014 SET UNIQUE_CHECKS=@OLD_UNIQUE_CHECKS */;
/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
/*!40111 SET SQL_NOTES=@OLD_SQL_NOTES */;

-- Dump completed on 2025-12-02 22:17:22

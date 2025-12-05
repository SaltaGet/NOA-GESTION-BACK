/*M!999999\- enable the sandbox mode */ 
-- MariaDB dump 10.19  Distrib 10.11.13-MariaDB, for debian-linux-gnu (x86_64)
--
-- Host: 127.0.0.1    Database: noa_gestion
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
-- Current Database: `noa_gestion`
--

CREATE DATABASE /*!32312 IF NOT EXISTS*/ `noa_gestion` /*!40100 DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci */;

USE `noa_gestion`;

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
) ENGINE=InnoDB AUTO_INCREMENT=2 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `admins`
--

LOCK TABLES `admins` WRITE;
/*!40000 ALTER TABLE `admins` DISABLE KEYS */;
INSERT INTO `admins` VALUES
(1,'saltaget','saltaget','SaltaGet','saltaget@gmail.com','$argon2id$v=19$m=65536,t=1,p=16$/N7xVuB83ErUEkEcKuQgdA$+VMtFWmrUqx0oqoPbIUaKwTAyY7NIIHpFc5F+vghODo',1,'2025-12-04 03:23:07.197','2025-12-04 03:23:07.197',NULL);
/*!40000 ALTER TABLE `admins` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `plans`
--

DROP TABLE IF EXISTS `plans`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8mb4 */;
CREATE TABLE `plans` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `name` varchar(255) DEFAULT NULL,
  `price_mounthly` decimal(10,2) DEFAULT NULL,
  `price_yearly` decimal(10,2) DEFAULT NULL,
  `description` text DEFAULT NULL,
  `features` text DEFAULT NULL,
  `amount_point_sale` bigint(20) NOT NULL,
  `amount_member` bigint(20) NOT NULL,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_plans_name` (`name`)
) ENGINE=InnoDB AUTO_INCREMENT=34 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `plans`
--

LOCK TABLES `plans` WRITE;
/*!40000 ALTER TABLE `plans` DISABLE KEYS */;
INSERT INTO `plans` VALUES
(1,'Básico',25.00,250.00,'Plan básico','emmmm, nada es básico, asi que no esperes mucho',1,5,'2025-12-04 03:23:07.142','2025-12-04 03:23:07.142');
/*!40000 ALTER TABLE `plans` ENABLE KEYS */;
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
-- Table structure for table `tenants`
--

DROP TABLE IF EXISTS `tenants`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8mb4 */;
CREATE TABLE `tenants` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `name` longtext NOT NULL,
  `identifier` varchar(191) NOT NULL,
  `address` longtext NOT NULL,
  `phone` longtext NOT NULL,
  `email` varchar(191) NOT NULL,
  `cuit_pdv` varchar(50) NOT NULL,
  `is_active` tinyint(1) NOT NULL DEFAULT 1,
  `plan_id` bigint(20) NOT NULL,
  `connection` longtext NOT NULL,
  `expiration` datetime(3) DEFAULT NULL,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `deleted_at` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_tenants_cuit_pdv` (`cuit_pdv`),
  UNIQUE KEY `uni_tenants_identifier` (`identifier`),
  KEY `idx_tenants_email` (`email`),
  KEY `idx_tenants_deleted_at` (`deleted_at`),
  KEY `fk_plans_tenants` (`plan_id`),
  CONSTRAINT `fk_plans_tenants` FOREIGN KEY (`plan_id`) REFERENCES `plans` (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=4 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `tenants`
--

LOCK TABLES `tenants` WRITE;
/*!40000 ALTER TABLE `tenants` DISABLE KEYS */;
INSERT INTO `tenants` VALUES
(1,'string','tenant1','string','string','t@gmail.com','string',1,1,'KolX02QjsTly/WXMiyaZUBs/B4TyVQjg5E+xp66kxjJVXqwAoxCqJsqehZKIFk3wKIpnarAI+myA8AwlWR9pKZ6caOydAPfsrWbXSD8lnNWeTcKaov+irO/WWeO1B4HyHvBGpgVgHO6xBW0KS13zw3P3u80HD6uZfg==','2025-12-12 00:00:00.000','2025-12-04 03:24:47.113','2025-12-04 03:31:38.166',NULL),
(3,'string','tenant2','string','string','tenant2@gmail.com','string2',1,1,'Kt3TUQDMNcWO/45w9r7UHqlUJJSbgOtQFCJb2JRxpVDjm8Sef8nVoVDatAngCMOH6GtDcX3OhTBIXz6hhqxjBWDbjrCR2b1Q+9ALXZNsG/jRQ0pbJnl2dcNJQx/GxGMHpRw0or4y5e7/V22bG7wK1PRG9FwJ8Q==',NULL,'2025-12-05 02:04:49.296','2025-12-05 02:04:49.296',NULL);
/*!40000 ALTER TABLE `tenants` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `user_tenants`
--

DROP TABLE IF EXISTS `user_tenants`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8mb4 */;
CREATE TABLE `user_tenants` (
  `user_id` bigint(20) NOT NULL,
  `tenant_id` bigint(20) NOT NULL,
  KEY `fk_user_tenants_user` (`user_id`),
  KEY `fk_tenants_user_tenants` (`tenant_id`),
  CONSTRAINT `fk_tenants_user_tenants` FOREIGN KEY (`tenant_id`) REFERENCES `tenants` (`id`),
  CONSTRAINT `fk_user_tenants_user` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `user_tenants`
--

LOCK TABLES `user_tenants` WRITE;
/*!40000 ALTER TABLE `user_tenants` DISABLE KEYS */;
INSERT INTO `user_tenants` VALUES
(1,1),
(2,3);
/*!40000 ALTER TABLE `user_tenants` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `users`
--

DROP TABLE IF EXISTS `users`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8mb4 */;
CREATE TABLE `users` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `first_name` varchar(30) NOT NULL,
  `last_name` varchar(30) NOT NULL,
  `email` varchar(191) NOT NULL,
  `username` varchar(191) NOT NULL,
  `address` varchar(255) DEFAULT NULL,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uni_users_email` (`email`),
  UNIQUE KEY `uni_users_username` (`username`)
) ENGINE=InnoDB AUTO_INCREMENT=3 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `users`
--

LOCK TABLES `users` WRITE;
/*!40000 ALTER TABLE `users` DISABLE KEYS */;
INSERT INTO `users` VALUES
(1,'string','string','user1@gmail.com','user1','string','2025-12-04 03:24:47.114','2025-12-04 03:24:47.114'),
(2,'string','string','user2@gmail.com','user2','string','2025-12-05 02:04:49.297','2025-12-05 02:04:49.297');
/*!40000 ALTER TABLE `users` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Dumping events for database 'noa_gestion'
--

--
-- Dumping routines for database 'noa_gestion'
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

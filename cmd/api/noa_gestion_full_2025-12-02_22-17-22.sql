/*M!999999\- enable the sandbox mode */ 
-- MariaDB dump 10.19  Distrib 10.11.13-MariaDB, for debian-linux-gnu (x86_64)
--
-- Host: 127.0.0.1    Database: noa_gestion
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
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_admins_username` (`username`),
  UNIQUE KEY `uni_admins_email` (`email`)
) ENGINE=InnoDB AUTO_INCREMENT=2 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `admins`
--

LOCK TABLES `admins` WRITE;
/*!40000 ALTER TABLE `admins` DISABLE KEYS */;
INSERT INTO `admins` VALUES
(1,'saltaget','saltaget','SaltaGet','saltaget@gmail.com','$argon2id$v=19$m=65536,t=1,p=16$acZovIbYDqmY9S0gtZuFxg$L2UJn0jefZjiEMVFYi9gdjoVefNrmbpudJLn91V9JDs',1,'2025-11-15 05:21:51.811','2025-11-15 05:21:51.811');
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
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=301 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `plans`
--

LOCK TABLES `plans` WRITE;
/*!40000 ALTER TABLE `plans` DISABLE KEYS */;
INSERT INTO `plans` VALUES
(1,'Free'),
(2,'Basic'),
(3,'Premium'),
(4,'Free'),
(5,'Basic'),
(6,'Premium'),
(7,'Free'),
(8,'Basic'),
(9,'Premium'),
(10,'Free'),
(11,'Basic'),
(12,'Premium'),
(13,'Free'),
(14,'Basic'),
(15,'Premium'),
(16,'Free'),
(17,'Basic'),
(18,'Premium'),
(19,'Free'),
(20,'Basic'),
(21,'Premium'),
(22,'Free'),
(23,'Basic'),
(24,'Premium'),
(25,'Free'),
(26,'Basic'),
(27,'Premium'),
(28,'Free'),
(29,'Basic'),
(30,'Premium'),
(31,'Free'),
(32,'Basic'),
(33,'Premium'),
(34,'Free'),
(35,'Basic'),
(36,'Premium'),
(37,'Free'),
(38,'Basic'),
(39,'Premium'),
(40,'Free'),
(41,'Basic'),
(42,'Premium'),
(43,'Free'),
(44,'Basic'),
(45,'Premium'),
(46,'Free'),
(47,'Basic'),
(48,'Premium'),
(49,'Free'),
(50,'Basic'),
(51,'Premium'),
(52,'Free'),
(53,'Basic'),
(54,'Premium'),
(55,'Free'),
(56,'Basic'),
(57,'Premium'),
(58,'Free'),
(59,'Basic'),
(60,'Premium'),
(61,'Free'),
(62,'Basic'),
(63,'Premium'),
(64,'Free'),
(65,'Basic'),
(66,'Premium'),
(67,'Free'),
(68,'Basic'),
(69,'Premium'),
(70,'Free'),
(71,'Basic'),
(72,'Premium'),
(73,'Free'),
(74,'Basic'),
(75,'Premium'),
(76,'Free'),
(77,'Basic'),
(78,'Premium'),
(79,'Free'),
(80,'Basic'),
(81,'Premium'),
(82,'Free'),
(83,'Basic'),
(84,'Premium'),
(85,'Free'),
(86,'Basic'),
(87,'Premium'),
(88,'Free'),
(89,'Basic'),
(90,'Premium'),
(91,'Free'),
(92,'Basic'),
(93,'Premium'),
(94,'Free'),
(95,'Basic'),
(96,'Premium'),
(97,'Free'),
(98,'Basic'),
(99,'Premium'),
(100,'Free'),
(101,'Basic'),
(102,'Premium'),
(103,'Free'),
(104,'Basic'),
(105,'Premium'),
(106,'Free'),
(107,'Basic'),
(108,'Premium'),
(109,'Free'),
(110,'Basic'),
(111,'Premium'),
(112,'Free'),
(113,'Basic'),
(114,'Premium'),
(115,'Free'),
(116,'Basic'),
(117,'Premium'),
(118,'Free'),
(119,'Basic'),
(120,'Premium'),
(121,'Free'),
(122,'Basic'),
(123,'Premium'),
(124,'Free'),
(125,'Basic'),
(126,'Premium'),
(127,'Free'),
(128,'Basic'),
(129,'Premium'),
(130,'Free'),
(131,'Basic'),
(132,'Premium'),
(133,'Free'),
(134,'Basic'),
(135,'Premium'),
(136,'Free'),
(137,'Basic'),
(138,'Premium'),
(139,'Free'),
(140,'Basic'),
(141,'Premium'),
(142,'Free'),
(143,'Basic'),
(144,'Premium'),
(145,'Free'),
(146,'Basic'),
(147,'Premium'),
(148,'Free'),
(149,'Basic'),
(150,'Premium'),
(151,'Free'),
(152,'Basic'),
(153,'Premium'),
(154,'Free'),
(155,'Basic'),
(156,'Premium'),
(157,'Free'),
(158,'Basic'),
(159,'Premium'),
(160,'Free'),
(161,'Basic'),
(162,'Premium'),
(163,'Free'),
(164,'Basic'),
(165,'Premium'),
(166,'Free'),
(167,'Basic'),
(168,'Premium'),
(169,'Free'),
(170,'Basic'),
(171,'Premium'),
(172,'Free'),
(173,'Basic'),
(174,'Premium'),
(175,'Free'),
(176,'Basic'),
(177,'Premium'),
(178,'Free'),
(179,'Basic'),
(180,'Premium'),
(181,'Free'),
(182,'Basic'),
(183,'Premium'),
(184,'Free'),
(185,'Basic'),
(186,'Premium'),
(187,'Free'),
(188,'Basic'),
(189,'Premium'),
(190,'Free'),
(191,'Basic'),
(192,'Premium'),
(193,'Free'),
(194,'Basic'),
(195,'Premium'),
(196,'Free'),
(197,'Basic'),
(198,'Premium'),
(199,'Free'),
(200,'Basic'),
(201,'Premium'),
(202,'Free'),
(203,'Basic'),
(204,'Premium'),
(205,'Free'),
(206,'Basic'),
(207,'Premium'),
(208,'Free'),
(209,'Basic'),
(210,'Premium'),
(211,'Free'),
(212,'Basic'),
(213,'Premium'),
(214,'Free'),
(215,'Basic'),
(216,'Premium'),
(217,'Free'),
(218,'Basic'),
(219,'Premium'),
(220,'Free'),
(221,'Basic'),
(222,'Premium'),
(223,'Free'),
(224,'Basic'),
(225,'Premium'),
(226,'Free'),
(227,'Basic'),
(228,'Premium'),
(229,'Free'),
(230,'Basic'),
(231,'Premium'),
(232,'Free'),
(233,'Basic'),
(234,'Premium'),
(235,'Free'),
(236,'Basic'),
(237,'Premium'),
(238,'Free'),
(239,'Basic'),
(240,'Premium'),
(241,'Free'),
(242,'Basic'),
(243,'Premium'),
(244,'Free'),
(245,'Basic'),
(246,'Premium'),
(247,'Free'),
(248,'Basic'),
(249,'Premium'),
(250,'Free'),
(251,'Basic'),
(252,'Premium'),
(253,'Free'),
(254,'Basic'),
(255,'Premium'),
(256,'Free'),
(257,'Basic'),
(258,'Premium'),
(259,'Free'),
(260,'Basic'),
(261,'Premium'),
(262,'Free'),
(263,'Basic'),
(264,'Premium'),
(265,'Free'),
(266,'Basic'),
(267,'Premium'),
(268,'Free'),
(269,'Basic'),
(270,'Premium'),
(271,'Free'),
(272,'Basic'),
(273,'Premium'),
(274,'Free'),
(275,'Basic'),
(276,'Premium'),
(277,'Free'),
(278,'Basic'),
(279,'Premium'),
(280,'Free'),
(281,'Basic'),
(282,'Premium'),
(283,'Free'),
(284,'Basic'),
(285,'Premium'),
(286,'Free'),
(287,'Basic'),
(288,'Premium'),
(289,'Free'),
(290,'Basic'),
(291,'Premium'),
(292,'Free'),
(293,'Basic'),
(294,'Premium'),
(295,'Free'),
(296,'Basic'),
(297,'Premium'),
(298,'Free'),
(299,'Basic'),
(300,'Premium');
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
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_tenants_cuit_pdv` (`cuit_pdv`),
  UNIQUE KEY `uni_tenants_identifier` (`identifier`),
  KEY `idx_tenants_email` (`email`),
  KEY `fk_plans_tenants` (`plan_id`),
  CONSTRAINT `fk_plans_tenants` FOREIGN KEY (`plan_id`) REFERENCES `plans` (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=2 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `tenants`
--

LOCK TABLES `tenants` WRITE;
/*!40000 ALTER TABLE `tenants` DISABLE KEYS */;
INSERT INTO `tenants` VALUES
(1,'string','tenant1','string','string','tenant1@gmail.com','string',1,1,'NmkGXhBePfjl6dJ7GQMcYMHf/gt+ko9rsaIaZjBrHS9yI0Iw5G+rvScvOk8lWIA0alwYSyq6/y9OpZA2Cz7cgc0i1/jIeRgS0+tTHYVrbOc24fqwt4CUiOAZBYhQjL8uX44d0d8J7tacW+IwQdfDkuNKdT8yvvPERg==','2025-11-15 07:57:45.852','2025-11-15 07:57:45.852');
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
(1,1);
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
) ENGINE=InnoDB AUTO_INCREMENT=2 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `users`
--

LOCK TABLES `users` WRITE;
/*!40000 ALTER TABLE `users` DISABLE KEYS */;
INSERT INTO `users` VALUES
(1,'string','string','user1@gmail.com','user1','string','2025-11-15 07:57:45.853','2025-11-15 07:57:45.853');
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

-- Dump completed on 2025-12-02 22:17:22

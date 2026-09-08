/*M!999999\- enable the sandbox mode */ 
-- MariaDB dump 10.19-11.8.6-MariaDB, for debian-linux-gnu (x86_64)
--
-- Host: localhost    Database: go_inspector
-- ------------------------------------------------------
-- Server version	11.8.6-MariaDB-0+deb13u1 from Debian

/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/*!40101 SET @OLD_CHARACTER_SET_RESULTS=@@CHARACTER_SET_RESULTS */;
/*!40101 SET @OLD_COLLATION_CONNECTION=@@COLLATION_CONNECTION */;
/*!40101 SET NAMES utf8mb4 */;
/*!40103 SET @OLD_TIME_ZONE=@@TIME_ZONE */;
/*!40103 SET TIME_ZONE='+00:00' */;
/*!40014 SET @OLD_UNIQUE_CHECKS=@@UNIQUE_CHECKS, UNIQUE_CHECKS=0 */;
/*!40014 SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0 */;
/*!40101 SET @OLD_SQL_MODE=@@SQL_MODE, SQL_MODE='NO_AUTO_VALUE_ON_ZERO' */;
/*M!100616 SET @OLD_NOTE_VERBOSITY=@@NOTE_VERBOSITY, NOTE_VERBOSITY=0 */;

--
-- Table structure for table `check_lastest_time`
--

DROP TABLE IF EXISTS `check_lastest_time`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8mb4 */;
CREATE TABLE `check_lastest_time` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `check_time_id` datetime DEFAULT NULL COMMENT '最新的检查时间',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `check_lastest_time`
--

SET @OLD_AUTOCOMMIT=@@AUTOCOMMIT, @@AUTOCOMMIT=0;
LOCK TABLES `check_lastest_time` WRITE;
/*!40000 ALTER TABLE `check_lastest_time` DISABLE KEYS */;
/*!40000 ALTER TABLE `check_lastest_time` ENABLE KEYS */;
UNLOCK TABLES;
COMMIT;
SET AUTOCOMMIT=@OLD_AUTOCOMMIT;

--
-- Table structure for table `host_check_items`
--

DROP TABLE IF EXISTS `host_check_items`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8mb4 */;
CREATE TABLE `host_check_items` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `host_id` varchar(16) DEFAULT NULL COMMENT '主机ip',
  `check_time_id` datetime DEFAULT NULL COMMENT '检查时间',
  `item_name` varchar(64) DEFAULT NULL COMMENT '指标名',
  `item_value` varchar(128) DEFAULT NULL,
  `status` varchar(16) DEFAULT NULL COMMENT '指标状态',
  `detail` text DEFAULT NULL COMMENT '详情或报错日志',
  PRIMARY KEY (`id`),
  KEY `idx_host_time` (`host_id`,`check_time_id`),
  KEY `idx_item_status` (`item_name`,`status`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `host_check_items`
--

SET @OLD_AUTOCOMMIT=@@AUTOCOMMIT, @@AUTOCOMMIT=0;
LOCK TABLES `host_check_items` WRITE;
/*!40000 ALTER TABLE `host_check_items` DISABLE KEYS */;
/*!40000 ALTER TABLE `host_check_items` ENABLE KEYS */;
UNLOCK TABLES;
COMMIT;
SET AUTOCOMMIT=@OLD_AUTOCOMMIT;

--
-- Table structure for table `host_check_tasks`
--

DROP TABLE IF EXISTS `host_check_tasks`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8mb4 */;
CREATE TABLE `host_check_tasks` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `host_id` varchar(16) DEFAULT NULL COMMENT '主机ip',
  `check_time_id` datetime DEFAULT NULL COMMENT '检查时间',
  `status` varchar(16) DEFAULT NULL COMMENT '整体检查状态',
  `error` text DEFAULT NULL COMMENT '整机巡检全局异常信息',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_host_time` (`host_id`,`check_time_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `host_check_tasks`
--

SET @OLD_AUTOCOMMIT=@@AUTOCOMMIT, @@AUTOCOMMIT=0;
LOCK TABLES `host_check_tasks` WRITE;
/*!40000 ALTER TABLE `host_check_tasks` DISABLE KEYS */;
/*!40000 ALTER TABLE `host_check_tasks` ENABLE KEYS */;
UNLOCK TABLES;
COMMIT;
SET AUTOCOMMIT=@OLD_AUTOCOMMIT;

--
-- Table structure for table `host_config`
--

DROP TABLE IF EXISTS `host_config`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8mb4 */;
CREATE TABLE `host_config` (
  `Host` varchar(16) DEFAULT NULL COMMENT '主机ip',
  `User` varchar(20) DEFAULT NULL COMMENT '用户名',
  `Port` smallint(6) DEFAULT NULL COMMENT '端口',
  `Password` varchar(20) DEFAULT NULL COMMENT '密码',
  `KeyFile` varchar(30) DEFAULT NULL COMMENT '私钥文件'
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `host_config`
--

SET @OLD_AUTOCOMMIT=@@AUTOCOMMIT, @@AUTOCOMMIT=0;
LOCK TABLES `host_config` WRITE;
/*!40000 ALTER TABLE `host_config` DISABLE KEYS */;
/*!40000 ALTER TABLE `host_config` ENABLE KEYS */;
UNLOCK TABLES;
COMMIT;
SET AUTOCOMMIT=@OLD_AUTOCOMMIT;
/*!40103 SET TIME_ZONE=@OLD_TIME_ZONE */;

/*!40101 SET SQL_MODE=@OLD_SQL_MODE */;
/*!40014 SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS */;
/*!40014 SET UNIQUE_CHECKS=@OLD_UNIQUE_CHECKS */;
/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
/*M!100616 SET NOTE_VERBOSITY=@OLD_NOTE_VERBOSITY */;

-- Dump completed on 2026-09-01 21:41:12

/*M!999999\- enable the sandbox mode */ 
-- MariaDB dump 10.19-11.7.2-MariaDB, for Linux (x86_64)
--
-- Host: localhost    Database: coalgdb
-- ------------------------------------------------------
-- Server version	11.7.2-MariaDB

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
-- Table structure for table `capitol`
--

DROP TABLE IF EXISTS `capitol`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8mb4 */;
CREATE TABLE `capitol` (
  `id_cap` int(11) NOT NULL AUTO_INCREMENT,
  `titlu` varchar(64) NOT NULL,
  `id_unit` int(11) NOT NULL,
  PRIMARY KEY (`id_cap`),
  KEY `fk_cap_unit` (`id_unit`),
  CONSTRAINT `fk_cap_unit` FOREIGN KEY (`id_unit`) REFERENCES `unitate` (`id_unit`)
) ENGINE=InnoDB AUTO_INCREMENT=7 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `capitol`
--

LOCK TABLES `capitol` WRITE;
/*!40000 ALTER TABLE `capitol` DISABLE KEYS */;
INSERT INTO `capitol` VALUES
(0,'null_cap',12),
(1,'Divide et impera',10),
(2,'Introducere in C++',9);
/*!40000 ALTER TABLE `capitol` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `pagina`
--

DROP TABLE IF EXISTS `pagina`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8mb4 */;
CREATE TABLE `pagina` (
  `id_pag` int(11) NOT NULL AUTO_INCREMENT,
  `titlu` varchar(64) NOT NULL,
  `descriere` text DEFAULT NULL,
  `autor` varchar(64) DEFAULT NULL,
  `aprobat` tinyint(1) NOT NULL DEFAULT 0,
  `id_cap` int(11) NOT NULL,
  PRIMARY KEY (`id_pag`),
  KEY `fk_pag_cap` (`id_cap`),
  CONSTRAINT `fk_pag_cap` FOREIGN KEY (`id_cap`) REFERENCES `capitol` (`id_cap`) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB AUTO_INCREMENT=3 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `pagina`
--

LOCK TABLES `pagina` WRITE;
/*!40000 ALTER TABLE `pagina` DISABLE KEYS */;
INSERT INTO `pagina` VALUES
(0,'null_pag',NULL,NULL,0,0),
/*!40000 ALTER TABLE `pagina` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `problema`
--

DROP TABLE IF EXISTS `problema`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8mb4 */;
CREATE TABLE `problema` (
  `id_prob` int(11) NOT NULL AUTO_INCREMENT,
  `titlu` varchar(64) DEFAULT '?',
  `link` varchar(256) DEFAULT NULL,
  `id_pag` int(11) NOT NULL,
  PRIMARY KEY (`id_prob`),
  KEY `fk_prob_pag` (`id_pag`),
  CONSTRAINT `fk_prob_pag` FOREIGN KEY (`id_pag`) REFERENCES `pagina` (`id_pag`) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB AUTO_INCREMENT=3 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `problema`
--

LOCK TABLES `problema` WRITE;
/*!40000 ALTER TABLE `problema` DISABLE KEYS */;
INSERT INTO `problema` VALUES
(0,'?',NULL,0),
/*!40000 ALTER TABLE `problema` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `resursa`
--

DROP TABLE IF EXISTS `resursa`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8mb4 */;
CREATE TABLE `resursa` (
  `id_res` int(11) NOT NULL AUTO_INCREMENT,
  `titlu` varchar(64) DEFAULT '?',
  `link` varchar(256) DEFAULT NULL,
  `id_pag` int(11) NOT NULL,
  PRIMARY KEY (`id_res`),
  KEY `fk_res_pag` (`id_pag`),
  CONSTRAINT `fk_res_pag` FOREIGN KEY (`id_pag`) REFERENCES `pagina` (`id_pag`) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB AUTO_INCREMENT=3 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `resursa`
--

LOCK TABLES `resursa` WRITE;
/*!40000 ALTER TABLE `resursa` DISABLE KEYS */;
INSERT INTO `resursa` VALUES
(0,'?',NULL,0),
/*!40000 ALTER TABLE `resursa` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `unitate`
--

DROP TABLE IF EXISTS `unitate`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8mb4 */;
CREATE TABLE `unitate` (
  `id_unit` int(11) NOT NULL,
  PRIMARY KEY (`id_unit`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `unitate`
--

LOCK TABLES `unitate` WRITE;
/*!40000 ALTER TABLE `unitate` DISABLE KEYS */;
INSERT INTO `unitate` VALUES
(9),
(10),
(11),
(12);
/*!40000 ALTER TABLE `unitate` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `user`
--

DROP TABLE IF EXISTS `user`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8mb4 */;
CREATE TABLE `user` (
  `id_user` int(11) NOT NULL AUTO_INCREMENT,
  `name` varchar(256) DEFAULT NULL,
  `parola` varchar(256) DEFAULT NULL,
  `admin` tinyint(4) NOT NULL DEFAULT 0,
  `sessiontoken` varchar(255) DEFAULT NULL,
  PRIMARY KEY (`id_user`)
) ENGINE=InnoDB AUTO_INCREMENT=3 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `user`
--

LOCK TABLES `user` WRITE;
/*!40000 ALTER TABLE `user` DISABLE KEYS */;
INSERT INTO `user` VALUES
(1,'admin','$2a$10$.Iod2zh0Zz7prPgdoj4XwuzxQOeitVJ3Su5nAX5p4tMaiDcXcKlKK',1,'K9QfKdOlAol7wyshTgrCSZ8fNgQ0ADa-4AJaaIlPPAA='),
(2,'matei','$2a$10$jGt7f1ywzZxWQ63k8qzVyeOyonirAyDFkxDhFJ75ToqpfkpsYh1fe',0,'');
/*!40000 ALTER TABLE `user` ENABLE KEYS */;
UNLOCK TABLES;
/*!40103 SET TIME_ZONE=@OLD_TIME_ZONE */;

/*!40101 SET SQL_MODE=@OLD_SQL_MODE */;
/*!40014 SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS */;
/*!40014 SET UNIQUE_CHECKS=@OLD_UNIQUE_CHECKS */;
/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
/*M!100616 SET NOTE_VERBOSITY=@OLD_NOTE_VERBOSITY */;

-- Dump completed on 2025-02-14 19:01:02

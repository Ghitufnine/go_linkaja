/*
 Navicat Premium Data Transfer

 Source Server         : local
 Source Server Type    : MySQL
 Source Server Version : 100140
 Source Host           : localhost:3306
 Source Schema         : go_linkaja

 Target Server Type    : MySQL
 Target Server Version : 100140
 File Encoding         : 65001

 Date: 30/08/2023 22:02:19
*/

SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

-- ----------------------------
-- Table structure for account
-- ----------------------------
DROP TABLE IF EXISTS `account`;
CREATE TABLE `account`  (
  `account_number` int(10) NOT NULL,
  `customer_number` int(10) NULL DEFAULT NULL,
  `balance` int(20) NULL DEFAULT NULL,
  PRIMARY KEY (`account_number`) USING BTREE,
  INDEX `customer_number_account_customer`(`customer_number`) USING BTREE,
  CONSTRAINT `customer_number_account_customer` FOREIGN KEY (`customer_number`) REFERENCES `customer` (`customer_number`) ON DELETE RESTRICT ON UPDATE RESTRICT
) ENGINE = InnoDB CHARACTER SET = latin1 COLLATE = latin1_swedish_ci ROW_FORMAT = Compact;

-- ----------------------------
-- Records of account
-- ----------------------------
INSERT INTO `account` VALUES (555001, 1001, 10000);
INSERT INTO `account` VALUES (555002, 1002, 15000);

-- ----------------------------
-- Table structure for customer
-- ----------------------------
DROP TABLE IF EXISTS `customer`;
CREATE TABLE `customer`  (
  `customer_number` int(10) NULL DEFAULT NULL,
  `name` varchar(255) CHARACTER SET latin1 COLLATE latin1_swedish_ci NULL DEFAULT NULL,
  INDEX `customer_number`(`customer_number`) USING BTREE
) ENGINE = InnoDB CHARACTER SET = latin1 COLLATE = latin1_swedish_ci ROW_FORMAT = Compact;

-- ----------------------------
-- Records of customer
-- ----------------------------
INSERT INTO `customer` VALUES (1001, 'Bob Martin');
INSERT INTO `customer` VALUES (1002, 'Linus Torvalds');

SET FOREIGN_KEY_CHECKS = 1;

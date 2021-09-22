/*
 Navicat Premium Data Transfer

 Source Server         : localhost
 Source Server Type    : MySQL
 Source Server Version : 50734
 Source Host           : 127.0.0.1:3306
 Source Schema         : mxshop

 Target Server Type    : MySQL
 Target Server Version : 50734
 File Encoding         : 65001

 Date: 21/09/2021 10:52:41
*/

SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

-- ----------------------------
-- Table structure for coupons
-- ----------------------------
DROP TABLE IF EXISTS `coupons`;
CREATE TABLE `coupons`  (
  `id` int(10) UNSIGNED NOT NULL AUTO_INCREMENT,
  `num` int(11) UNSIGNED NOT NULL,
  PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 2 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of coupons
-- ----------------------------
INSERT INTO `coupons` VALUES (1, 17);

-- ----------------------------
-- Table structure for coupons_history
-- ----------------------------
DROP TABLE IF EXISTS `coupons_history`;
CREATE TABLE `coupons_history`  (
  `id` int(10) UNSIGNED NOT NULL AUTO_INCREMENT,
  `mobile` varchar(11) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL,
  `num` int(11) NOT NULL,
  `status` tinyint(1) NOT NULL COMMENT '(1:未扣减,2:已归还)',
  PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 30 CHARACTER SET = utf8 COLLATE = utf8_general_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of coupons_history
-- ----------------------------
INSERT INTO `coupons_history` VALUES (28, '13743253322', 13, 1);
INSERT INTO `coupons_history` VALUES (29, '13743253321', 13, 2);

-- ----------------------------
-- Table structure for coupons_user
-- ----------------------------
DROP TABLE IF EXISTS `coupons_user`;
CREATE TABLE `coupons_user`  (
  `id` int(10) UNSIGNED NOT NULL AUTO_INCREMENT,
  `mobile` varchar(11) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL,
  `num` int(11) UNSIGNED NOT NULL,
  PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 34 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of coupons_user
-- ----------------------------
INSERT INTO `coupons_user` VALUES (33, '13743253322', 13);

-- ----------------------------
-- Table structure for inventory
-- ----------------------------
DROP TABLE IF EXISTS `inventory`;
CREATE TABLE `inventory`  (
  `id` int(10) UNSIGNED NOT NULL AUTO_INCREMENT,
  `goods` int(11) NULL DEFAULT NULL COMMENT '商品id',
  `stocks` int(11) NULL DEFAULT NULL COMMENT '库存数量',
  `version` int(11) NULL DEFAULT NULL COMMENT '版本号(分布式锁的乐观锁)',
  `add_time` int(10) NOT NULL,
  `is_deteted` tinyint(1) NOT NULL,
  `update_time` int(10) NOT NULL,
  PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 4 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of inventory
-- ----------------------------
INSERT INTO `inventory` VALUES (1, 1, 20, 0, 0, 0, 0);
INSERT INTO `inventory` VALUES (2, 2, 20, 0, 0, 0, 0);
INSERT INTO `inventory` VALUES (3, 3, 20, 0, 0, 0, 0);

-- ----------------------------
-- Table structure for users
-- ----------------------------
DROP TABLE IF EXISTS `users`;
CREATE TABLE `users`  (
  `id` int(10) UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '主键',
  `mobile` varchar(11) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL COMMENT '手机号码',
  `password` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL COMMENT '密码',
  `nickname` varchar(20) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL COMMENT '昵称',
  `head_url` varchar(200) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL COMMENT '头像',
  `birthday` int(10) NULL DEFAULT NULL COMMENT '生日',
  `address` varchar(200) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL COMMENT '地址',
  `desc` varchar(200) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL COMMENT '描述',
  `gender` tinyint(4) NULL DEFAULT NULL COMMENT '性别（1：男；2：女）',
  `role` tinyint(4) NULL DEFAULT NULL COMMENT '角色（1：普通用户；2：管理员）',
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE INDEX `idx_users_mobile`(`mobile`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 65 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of users
-- ----------------------------
INSERT INTO `users` VALUES (63, '13743253322', '6a2919851df02ca28d6b8a3604bcde13', 'test', '', 0, '', '', 0, 0);
INSERT INTO `users` VALUES (64, '13743253321', '6a2919851df02ca28d6b8a3604bcde13', 'test', '', 0, '', '', 0, 0);

SET FOREIGN_KEY_CHECKS = 1;

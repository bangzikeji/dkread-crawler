/*
Navicat MySQL Data Transfer

Source Server         : 127.0.0.1
Source Server Version : 50553
Source Host           : localhost:3306
Source Database       : dkread

Target Server Type    : MYSQL
Target Server Version : 50553
File Encoding         : 65001

Date: 2020-09-08 23:27:11
*/

SET FOREIGN_KEY_CHECKS=0;

-- ----------------------------
-- Table structure for books
-- ----------------------------
DROP TABLE IF EXISTS `books`;
CREATE TABLE `books` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT COMMENT 'ID',
  `title` varchar(100) NOT NULL COMMENT '标题',
  `describe` varchar(255) NOT NULL COMMENT '描述简介',
  `cover` varchar(100) NOT NULL COMMENT '封面',
  `content` varchar(255) NOT NULL COMMENT '详细描述简介',
  `about` varchar(255) NOT NULL COMMENT '关于作者',
  `category` varchar(50) NOT NULL COMMENT '分类',
  `author` varchar(520) NOT NULL COMMENT '作者',
  `label` varchar(100) NOT NULL COMMENT '标签',
  `score` decimal(10,2) NOT NULL DEFAULT '0.00' COMMENT '评分',
  `download` varchar(100) DEFAULT NULL COMMENT '下载链接',
  `created` datetime NOT NULL COMMENT '创建时间',
  `updated` datetime NOT NULL COMMENT '更新时间',
  `status` tinyint(4) NOT NULL DEFAULT '1' COMMENT '0删除，1可用',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='书籍表';

-- ----------------------------
-- Records of books
-- ----------------------------

package database

// import(
//   "github.com/kataras/iris/context"
// )

const (
  IDP_DONGPIN = `CREATE DATABASE IF NOT EXISTS idongpin DEFAULT CHARSET utf8mb4 COLLATE utf8mb4_unicode_ci`

  IDP_AUTH = `CREATE TABLE idp_auth (
    id INT(11) NOT NULL AUTO_INCREMENT COMMENT 'id',
    content TEXT(10000) NOT NULL COMMENT 'json格式配置文件',
    updated_at INT(11) NULL DEFAULT NULL COMMENT '修改时间',
    created_at INT(11) NULL DEFAULT NULL COMMENT '创建时间',
    PRIMARY KEY (id)
  )`

  IDP_ADMIN = `CREATE TABLE idp_admins (
    id INT(11) NOT NULL AUTO_INCREMENT COMMENT 'id',
    phone VARCHAR(64) NULL DEFAULT '' COMMENT '手机号码（用来登陆）',
    password VARCHAR(64) NULL DEFAULT '' COMMENT '登录密码',
    realname VARCHAR(64) NULL DEFAULT '' COMMENT '姓名',
    nickname VARCHAR(64) NULL DEFAULT '' COMMENT '昵称',
    avatar VARCHAR(255) NULL DEFAULT '' COMMENT '头像',
    sex int(2) DEFAULT 0 COMMENT '性别：0未知、1男、2女',
    identity VARCHAR(255) NULL DEFAULT '' COMMENT '身份证号码',
    groups INT(2) NULL DEFAULT 1 COMMENT '用户组',
    roles INT(2) NULL DEFAULT 1 COMMENT '用户组里面的角色',
    state INT(2) NULL DEFAULT 1 COMMENT '账号状态：1启动、2禁用',
    login_count INT(11) NULL DEFAULT 0 COMMENT '登陆次数',
    login_time INT(11) NULL DEFAULT NULL COMMENT '登陆时间',
    last_time INT(11) NULL DEFAULT NULL COMMENT '上次登陆时间',
    login_ip VARCHAR(255) NULL DEFAULT NULL COMMENT '登陆Ip',
    last_ip VARCHAR(255) NULL DEFAULT NULL COMMENT '上次登陆Ip',
    deleted_at INT(11) NULL DEFAULT NULL COMMENT '删除时间',
    updated_at INT(11) NULL DEFAULT NULL COMMENT '修改时间',
    created_at INT(11) NULL DEFAULT NULL COMMENT '创建时间',
    PRIMARY KEY (id)
  )`

  IDP_ADMIN_GROUP = `CREATE TABLE idp_admins_group (
    id INT(11) NOT NULL AUTO_INCREMENT COMMENT 'id',
    name VARCHAR(64) NULL DEFAULT '' COMMENT '组名称',
    state INT(2) NULL DEFAULT 1 COMMENT '状态：1启动、2禁用',
    deleted_at INT(11) NULL DEFAULT NULL COMMENT '删除时间',
    updated_at INT(11) NULL DEFAULT NULL COMMENT '修改时间',
    created_at INT(11) NULL DEFAULT NULL COMMENT '创建时间',
    PRIMARY KEY (id)
  )`

  IDP_ADMIN_ROLE = `CREATE TABLE idp_admins_role (
    id INT(11) NOT NULL AUTO_INCREMENT COMMENT 'id',
    name VARCHAR(64) NULL DEFAULT '' COMMENT '角色名称',
    gid INT(11) NULL DEFAULT NULL COMMENT 'GROUP表关联Id',
    state INT(2) NULL DEFAULT 1 COMMENT '状态：1启动、2禁用',
    deleted_at INT(11) NULL DEFAULT NULL COMMENT '删除时间',
    updated_at INT(11) NULL DEFAULT NULL COMMENT '修改时间',
    created_at INT(11) NULL DEFAULT NULL COMMENT '创建时间',
    PRIMARY KEY (id)
  )`
)

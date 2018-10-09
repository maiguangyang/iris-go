package database

// import(
//   "github.com/kataras/iris/context"
// )

const (
  IDP_ADMIN = `CREATE TABLE idp_admins (
    id INT(11) NOT NULL AUTO_INCREMENT COMMENT 'id',
    username VARCHAR(64) NULL DEFAULT '' COMMENT '账号',
    password VARCHAR(64) NULL DEFAULT '' COMMENT '登录密码',
    nickname VARCHAR(64) NULL DEFAULT '' COMMENT '昵称',
    avatar VARCHAR(255) NULL DEFAULT '' COMMENT '头像',
    groups INT(2) NULL DEFAULT 0 COMMENT '用户组',
    roles INT(2) NULL DEFAULT 0 COMMENT '用户组里面的角色',
    state INT(2) NULL DEFAULT 0 COMMENT '账号状态：0启动、1禁用',
    login_count INT(11) NULL DEFAULT 0 COMMENT '登陆次数',
    login_time INT(11) NULL DEFAULT NULL COMMENT '登陆时间',
    last_time INT(11) NULL DEFAULT NULL COMMENT '上次登陆时间',
    deleted_at INT(11) NULL DEFAULT NULL COMMENT '删除时间',
    updated_at INT(11) NULL DEFAULT NULL COMMENT '修改时间',
    created_at INT(11) NULL DEFAULT NULL COMMENT '创建时间',
    PRIMARY KEY (id)
  )`

  IDP_ADMIN_GROUP = `CREATE TABLE idp_admins_group (
    id INT(11) NOT NULL AUTO_INCREMENT COMMENT 'id',
    name VARCHAR(64) NULL DEFAULT '' COMMENT '组名称',
    value INT(11) NULL DEFAULT 0 COMMENT '用户组',
    deleted_at INT(11) NULL DEFAULT NULL COMMENT '删除时间',
    updated_at INT(11) NULL DEFAULT NULL COMMENT '修改时间',
    created_at INT(11) NULL DEFAULT NULL COMMENT '创建时间',
    PRIMARY KEY (id)
  )`

  IDP_ADMIN_ROLE = `CREATE TABLE idp_admins_role (
    id INT(11) NOT NULL AUTO_INCREMENT COMMENT 'id',
    name VARCHAR(64) NULL DEFAULT '' COMMENT '角色名称',
    value INT(11) NULL DEFAULT 0 COMMENT '用户角色',
    deleted_at INT(11) NULL DEFAULT NULL COMMENT '删除时间',
    updated_at INT(11) NULL DEFAULT NULL COMMENT '修改时间',
    created_at INT(11) NULL DEFAULT NULL COMMENT '创建时间',
    PRIMARY KEY (id)
  )`
)

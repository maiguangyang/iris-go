package database

// 创建数据表
const (
  // 登录记录
  C_T_LOGIN_LOG = `CREATE TABLE login_log(
    id int(11) NOT NULL AUTO_INCREMENT COMMENT 'id',
    uid int(11) NOT NULL COMMENT '用户id',
    type int(2) NOT NULL COMMENT '登录类型：0系统配置登录，1后台管理登录，2用户登录',
    ip varchar(20) NOT NULL COMMENT 'ip',
    token varchar(1000) NOT NULL COMMENT 'token',
    content varchar(1000) NOT NULL COMMENT '登录提交的数据',
    PRIMARY KEY (id)
  )`

  // 登录错误记录
  C_T_LOGIN_ERROR_LOG = `CREATE TABLE login_error_log(
    id int(11) NOT NULL AUTO_INCREMENT COMMENT 'id',
    type int(2) NOT NULL COMMENT '登录类型：0系统配置登录，1后台管理登录，2用户登录',
    ip varchar(20) NOT NULL COMMENT 'ip',
    created_time int(11) NOT NULL COMMENT '创建时间',
    PRIMARY KEY (id)
  )`

  // 错误记录类型表（开发环境要同步正式环境的数据）
  C_T_ERROR_LOG_TYPE = `CREATE TABLE error_log_type(
    id int(11) NOT NULL AUTO_INCREMENT COMMENT 'id',
    name varchar(255) NOT NULL COMMENT '类型名',
    del int(2) DEFAULT '0' COMMENT '0正常、1删除',
    PRIMARY KEY (id)
  )`

  // 错误记录表
  C_T_ERROR_LOG = `CREATE TABLE error_log(
    id int(11) NOT NULL AUTO_INCREMENT COMMENT 'id',
    tid int(2) NOT NULL COMMENT '错误类型id',
    token varchar(1000) NOT NULL COMMENT 'token',
    data varchar(10000) NOT NULL COMMENT '错误内容',
    remark varchar(1000) DEFAULT '' COMMENT '备注',
    del int(2) DEFAULT '0' COMMENT '0正常、1删除',
    created_time int(11) NOT NULL COMMENT '创建时间',
    PRIMARY KEY (id)
  )`

  // 管理员表
  C_T_ADMIN = `CREATE TABLE admin(
    id int(11) NOT NULL AUTO_INCREMENT COMMENT 'id',
    phone varchar(12) NOT NULL COMMENT '手机号',
    password varchar(33) NOT NULL COMMENT '密码',
    name varchar(100) NOT NULL DEFAULT '' COMMENT '昵称',
    avarat varchar(255) NOT NULL DEFAULT '' COMMENT '头像',
    role int(2) NOT NULL DEFAULT '0' COMMENT '角色：0普通管理员、1超级管理员',
    state int(2) NOT NULL DEFAULT '0' COMMENT '状态：0正常，1禁用',
    created_time int(11) NOT NULL COMMENT '创建时间',
    PRIMARY KEY (id)
  )`
)
package database

// import(
//   "github.com/kataras/iris/context"
// )

const (
  IDP_DONGPIN = `CREATE DATABASE IF NOT EXISTS idongpin DEFAULT CHARSET utf8mb4 COLLATE utf8mb4_unicode_ci`

  // 管理员权限表
  IDP_ADMIN_AUTH = `CREATE TABLE idp_admin_auth (
    id INT(11) NOT NULL AUTO_INCREMENT COMMENT 'id',
    rid INT(11) NULL DEFAULT NUll COMMENT '角色id',
    sid VARCHAR(1000) NULL DEFAULT null COMMENT '权限表对应id',
    content TEXT(60000) NOT NULL COMMENT 'json格式配置文件',
    auth INT(2) NULL DEFAULT 2 COMMENT '跨部门查看：1/是，2/否',
    deleted_at datetime NULL DEFAULT NULL COMMENT '删除时间',
    updated_at datetime NULL DEFAULT NULL COMMENT '修改时间',
    created_at datetime NULL DEFAULT NULL COMMENT '创建时间',
    PRIMARY KEY (id)
  )`

  // 权限配置表
  IDP_AUTH_SET = `CREATE TABLE idp_auth_set (
    id INT(11) NOT NULL AUTO_INCREMENT COMMENT 'id',
    name VARCHAR(255) NULL DEFAULT '' COMMENT '权限名称',
    table_name VARCHAR(1000) NULL DEFAULT '' COMMENT '数据表',
    routes TEXT(60000) NOT NULL COMMENT '前端路由权限',
    sub_id VARCHAR(255) NULL DEFAULT '' COMMENT '附属权限arrar',
    deleted_at datetime NULL DEFAULT NULL COMMENT '删除时间',
    updated_at datetime NULL DEFAULT NULL COMMENT '修改时间',
    created_at datetime NULL DEFAULT NULL COMMENT '创建时间',
    PRIMARY KEY (id)
  )`

  // 管理员表
  IDP_ADMIN = `CREATE TABLE idp_admins (
    id INT(11) NOT NULL AUTO_INCREMENT COMMENT 'id',
    phone VARCHAR(64) NULL DEFAULT '' COMMENT '手机号码（用来登陆）',
    password VARCHAR(64) NULL DEFAULT '' COMMENT '登录密码',
    username VARCHAR(64) NULL DEFAULT '' COMMENT '姓名',
    sex int(2) DEFAULT 0 COMMENT '性别：0未知、1男、2女',
    super int(2) DEFAULT 1 COMMENT '超级账户：1/否，2/是',
    gid VARCHAR(255) NULL DEFAULT NULL COMMENT '部门：idp_admins_group表id: 1,2,3,4,5',
    rid VARCHAR(255) NULL DEFAULT NULL COMMENT '部门职位：idp_admins_role表id: 1,2,3,4,5',
    money int(11) DEFAULT 0 COMMENT '月薪',
    job_state int(2) DEFAULT NULL COMMENT '职位状态：1试用期、2转正、3离职',
    entry_time datetime NULL DEFAULT NULL COMMENT '入职时间',
    trial_time datetime NULL DEFAULT NULL COMMENT '试用期时间',
    contract_time datetime NULL DEFAULT NULL COMMENT '合同到期时间',
    quit_time datetime NULL DEFAULT NULL COMMENT '离职时间',
    state INT(2) NULL DEFAULT 1 COMMENT '账号状态：1启动、2禁用',
    aid INT(11) NULL DEFAULT 1 COMMENT '操作员：1系统添加、其他对应该表的id字段',
    login_count INT(11) NULL DEFAULT 0 COMMENT '登陆次数',
    login_time datetime NULL DEFAULT NULL COMMENT '登陆时间',
    last_time datetime NULL DEFAULT NULL COMMENT '上次登陆时间',
    login_ip VARCHAR(255) NULL DEFAULT NULL COMMENT '登陆Ip',
    last_ip VARCHAR(255) NULL DEFAULT NULL COMMENT '上次登陆Ip',
    deleted_at datetime NULL DEFAULT NULL COMMENT '删除时间',
    updated_at datetime NULL DEFAULT NULL COMMENT '修改时间',
    created_at datetime NULL DEFAULT NULL COMMENT '创建时间',
    PRIMARY KEY (id)
  )`

  // 人员档案表
  IDP_ADMIN_ARCHIVE = `CREATE TABLE idp_admin_archive (
    id INT(11) NOT NULL AUTO_INCREMENT COMMENT 'id',
    aid INT(11) NULL DEFAULT NUll COMMENT '用户id',
    avatar VARCHAR(255) NULL DEFAULT '' COMMENT '头像',
    school VARCHAR(255) NULL DEFAULT NULL COMMENT '毕业学校',
    major VARCHAR(255) NULL DEFAULT NULL COMMENT '专业',
    education int(2) DEFAULT 0 COMMENT '学历：0未填',
    nation int(2) DEFAULT 0 COMMENT '民族：0未填',
    native_place int(2) DEFAULT 0 COMMENT '籍贯：0未填',
    politics int(2) DEFAULT 0 COMMENT '政治面貌：0未填',
    marriage int(2) DEFAULT 0 COMMENT '婚否：0未填',
    healthy int(2) DEFAULT 0 COMMENT '健康：0未填',
    height int(2) DEFAULT 0 COMMENT '身高：0未填',
    weight int(2) DEFAULT 0 COMMENT '体重：0未填',
    identity VARCHAR(255) NULL DEFAULT '' COMMENT '身份证号码',
    residence VARCHAR(1000) NULL DEFAULT NULL COMMENT '户口所在地',
    place_residence VARCHAR(1000) NULL DEFAULT NULL COMMENT '现居住地',
    interest VARCHAR(1000) NULL DEFAULT NULL COMMENT '兴趣爱好',
    resume VARCHAR(255) NULL DEFAULT NULL COMMENT '个人简历',
    identity_file VARCHAR(255) NULL DEFAULT NULL COMMENT '身份证复印件',
    job_change int(2) DEFAULT NULL COMMENT '职务变更：0不变、其他的对应job_change表id',
    money_change int(11) DEFAULT NULL COMMENT '工资调整：0不变、其他的对应money_change表id',
    remarks VARCHAR(1000) NULL DEFAULT NULL COMMENT '备注',
    deleted_at datetime NULL DEFAULT NULL COMMENT '删除时间',
    updated_at datetime NULL DEFAULT NULL COMMENT '修改时间',
    created_at datetime NULL DEFAULT NULL COMMENT '创建时间',
    PRIMARY KEY (id)
  )`

  // 部门表
  IDP_ADMIN_GROUP = `CREATE TABLE idp_admin_group (
    id INT(11) NOT NULL AUTO_INCREMENT COMMENT 'id',
    name VARCHAR(64) NULL DEFAULT '' COMMENT '组名称',
    aid INT(11) NULL DEFAULT 1 COMMENT '操作员：1管理员添加、其他对应操作账户的id',
    state INT(2) NULL DEFAULT 1 COMMENT '状态：1启动、2禁用',
    deleted_at datetime NULL DEFAULT NULL COMMENT '删除时间',
    updated_at datetime NULL DEFAULT NULL COMMENT '修改时间',
    created_at datetime NULL DEFAULT NULL COMMENT '创建时间',
    PRIMARY KEY (id)
  )`

  // 角色表
  IDP_ADMIN_ROLE = `CREATE TABLE idp_admin_role (
    id INT(11) NOT NULL AUTO_INCREMENT COMMENT 'id',
    name VARCHAR(64) NULL DEFAULT '' COMMENT '角色名称',
    gid INT(11) NULL DEFAULT NULL COMMENT 'GROUP表关联Id',
    aid INT(11) NULL DEFAULT 1 COMMENT '操作员：1管理员添加、其他对应操作账户的id',
    state INT(2) NULL DEFAULT 1 COMMENT '状态：1启动、2禁用',
    deleted_at datetime NULL DEFAULT NULL COMMENT '删除时间',
    updated_at datetime NULL DEFAULT NULL COMMENT '修改时间',
    created_at datetime NULL DEFAULT NULL COMMENT '创建时间',
    PRIMARY KEY (id)
  )`
)

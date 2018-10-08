package database

// 检查必需表，并创建必需表
func HasTable() error {
  // 查询表
  rows, err := DB.Query("SHOW TABLES")
  if err != nil {
    // 还没选择库，退出
    if err.Error() == "Error 1046: No database selected" {
      return nil
    }
  } else {
    var isHaveLoginLog, isHaveLoginErrorLog, isHaveErrorLogType, isHaveErrorLog, isHaveAdmin bool
    defer rows.Close()
    for rows.Next() {
      var t string
      rows.Scan(&t)
      // 找到登录记录表、登录错误表、错误记录表、错误类型表
      if t == "login_log" {
        isHaveLoginLog = true
      } else if t == "login_error_log" {
        isHaveLoginErrorLog = true
      } else if t == "error_log_type" {
        isHaveErrorLogType = true
      } else if t == "error_log" {
        isHaveErrorLog = true
      } else if t == "admin" {
        isHaveAdmin = true
      }
    }

    if isHaveLoginLog == false {
      // 创建登录记录表
      _, err = DB.Exec(C_T_LOGIN_LOG)
      if err != nil {
        return err
      }
    }
    if isHaveLoginErrorLog == false {
      // 创建登录错误表
      _, err = DB.Exec(C_T_LOGIN_ERROR_LOG)
      if err != nil {
        return err
      }
    }
    if isHaveErrorLogType == false {
      // 创建错误类型表
      _, err = DB.Exec(C_T_ERROR_LOG_TYPE)
      if err != nil {
        return err
      }
    }
    if isHaveErrorLog == false {
      // 创建登录记录表
      _, err = DB.Exec(C_T_ERROR_LOG)
      if err != nil {
        return err
      }
    }
    if isHaveAdmin == false {
      // 创建管理员表
      _, err = DB.Exec(C_T_ADMIN)
      if err != nil {
        return err
      }
    }
  }

  return err
}
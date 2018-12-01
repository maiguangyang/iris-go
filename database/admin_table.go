package database

import (
  "database/sql"
)


type NullInt64 = *sql.NullInt64

type Model struct {
  Id int64 `json:"id" gorm:"primary_key;AUTO_INCREMENT"`
  CreatedAt NullInt64 `json:"created_at" gorm:"type:int(11);null;default:null"`
  UpdatedAt NullInt64 `json:"updated_at" gorm:"type:int(11);null;default:null"`
  DeletedAt NullInt64 `json:"deleted_at" gorm:"type:int(11);null;default:null"`
}

type ModelAt struct {
  CreatedAt NullInt64 `json:"created_at" gorm:"type:int(11);null;default:null"`
  UpdatedAt NullInt64 `json:"updated_at" gorm:"type:int(11);null;default:null"`
  DeletedAt NullInt64 `json:"deleted_at" gorm:"type:int(11);null;default:null"`
}


type AdminAuth23 struct {
  Model
  Rid int64 `json:"rid" gorm:"type:int(11);unique_index;default:null;comment:'角色id';"`
  Sid string `json:"sid" gorm:"type:varchar(1000);unique_index;default:null;comment:'权限表对应id';"`
  Content string `json:"content" gorm:"type:text(60000);default:null;comment:'json格式配置文件';"`
  Auth string `json:"auth" gorm:"type:int(2);default:2;comment:'跨部门查看：1/是，2/否';"`
  // ModelAt
}


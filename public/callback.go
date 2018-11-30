package public

import (
  "fmt"
  "time"
  "github.com/jinzhu/gorm"
)


func addExtraSpaceIfExist(str string) string {
  if str != "" {
    return " " + str
  }
  return ""
}

// updateTimeStampForCreateCallback will set `CreatedAt`, `UpdatedAt` when creating
func updateTimeStampForCreateCallback(scope *gorm.Scope) {
  if !scope.HasError() {
    now := time.Now().Unix()

    if createdAtField, ok := scope.FieldByName("CreatedAt"); ok {
      if createdAtField.IsBlank {
        createdAtField.Set(now)
      }
    }

    // if updatedAtField, ok := scope.FieldByName("UpdatedAt"); ok {
    //   if updatedAtField.IsBlank {
    //     updatedAtField.Set(now)
    //   }
    // }
  }
}

// updateTimeStampForUpdateCallback will set `UpdatedAt` when updating
func updateTimeStampForUpdateCallback(scope *gorm.Scope) {
  if _, ok := scope.Get("gorm:update_column"); !ok {
    scope.SetColumn("UpdatedAt", time.Now().Unix())
  }
}

// deleteCallback used to delete data from database or set deleted_at to current time (when using with soft delete)
func deleteCallback(scope *gorm.Scope) {
  if !scope.HasError() {
    var extraOption string
    if str, ok := scope.Get("gorm:delete_option"); ok {
      extraOption = fmt.Sprint(str)
    }

    deletedAtField, hasDeletedAtField := scope.FieldByName("DeletedAt")

    if !scope.Search.Unscoped && hasDeletedAtField {
      scope.Raw(fmt.Sprintf(
        "UPDATE %v SET %v=%v%v%v",
        scope.QuotedTableName(),
        scope.Quote(deletedAtField.DBName),
        // scope.AddToVars(NowFunc()),
        time.Now().Unix(),
        addExtraSpaceIfExist(scope.CombinedConditionSql()),
        addExtraSpaceIfExist(extraOption),
      )).Exec()
    } else {
      scope.Raw(fmt.Sprintf(
        "DELETE FROM %v%v%v",
        scope.QuotedTableName(),
        addExtraSpaceIfExist(scope.CombinedConditionSql()),
        addExtraSpaceIfExist(extraOption),
      )).Exec()
    }
  }
}

func InitGorm(db *gorm.DB) {
  db.Callback().Create().Remove("gorm:update_time_stamp")
  db.Callback().Update().Remove("gorm:update_time_stamp")
  // db.Callback().Update().Replace("gorm:update_time_stamp", updateCallback)
  db.Callback().Create().Replace("gorm:update_time_stamp", updateTimeStampForCreateCallback)
  db.Callback().Update().Replace("gorm:update_time_stamp", updateTimeStampForUpdateCallback)
  db.Callback().Delete().Replace("gorm:delete", deleteCallback)
}


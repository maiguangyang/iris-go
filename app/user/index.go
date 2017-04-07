package user

import (
	"encoding/json"
	"gopkg.in/kataras/iris.v6"
)

var responseData interface{}

func UserSay(ctx *iris.Context) {

	jsonStr := `{
    "code": 200,
    "data": [{"id": 1, "name": "test1","age": 18}, {"id": 2, "name": "test2","age": 24}],
    "msg" : "success"
  }`

	err := json.Unmarshal([]byte(jsonStr), &responseData)

	if err == nil {
		ctx.JSON(iris.StatusOK, responseData)
	}

}

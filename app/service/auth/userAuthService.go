package auth

import (
	"context"
	"log"
	"zhyu/app/common"
	"zhyu/app/dao"
)

// UserAuth 验证用户是否有效
func UserAuth(ctx context.Context, userId int64) map[string]any {

	// 创建返回信息
	data := make(map[string]map[string]any)
	code := common.INVALID_PARAMS

	// 验证用户
	resultData, isExist := dao.UserGetById(userId)
	log.Printf("resultData: %#v", resultData)
	if isExist {
		data["userInfo"] = make(map[string]any)
		data["userInfo"]["id"] = resultData.Id
		data["userInfo"]["userName"] = resultData.Username
		code = common.SUCCESS
	} else {
		code = common.ERROR_AUTH
	}

	return common.Result(code, common.GetMsg(code), data)
}

package v1

import (
	"context"
	"log"
	"strconv"
	"zhyu/app"
	"zhyu/app/common"
	"zhyu/app/dao"
	"zhyu/app/service"
)

type userData struct {
	Id        int64  `json:"id"`
	Name      string `json:"name"`
	Age       uint8  `json:"age"`
	Address   string `json:"address"`
	Phone     string `json:"phone"`
	Photo     string `json:"photo"`
	Status    uint8  `json:"status"`
	CreatedAt int64  `json:"created"`
}

// UserListService 获取用户列表
func UserListService(ctx context.Context, c *app.Context) any {

	// 创建返回信息
	resultData := make(map[string]any)
	code := common.SUCCESS

	// 参数
	var page int
	var err error
	if page, err = strconv.Atoi(c.PostForm("page")); err != nil {
		page = 1
	}
	if page <= 0 {
		page = 1
	}
	var pageSize int
	if pageSize, err = strconv.Atoi(c.PostForm("pageSize")); err != nil {
		pageSize = 3
	}

	// 查询列表
	where := "status = ? and password = ?"
	args := []any{
		"1", "123456",
	}
	offset := (page - 1) * pageSize
	limit := pageSize
	userList, total := dao.UserGetList(limit, offset, where, args...)
	log.Printf("userlist: %#v, total: %#v", userList, total)

	length := len(*userList)
	resultList := make([]userData, length)
	for key, item := range *userList {
		//resultData2 := userData{
		//	Id:      item.Id,
		//	Name:    item.Name,
		//	Age:     item.Age,
		//	Address: item.Address,
		//	Phone:   item.Phone,
		//	Photo:   item.Photo,
		//	Status:  item.Status,
		//}
		//resultList[key] = resultData2

		resultList[key].Id = item.Id
		resultList[key].Name = item.Name
		resultList[key].Age = item.Age
		resultList[key].Address = item.Address
		resultList[key].Phone = item.Phone
		resultList[key].Photo = item.Photo
		resultList[key].Status = item.Status
		resultList[key].CreatedAt = item.CreatedAt
	}
	resultData["list"] = resultList
	resultData["page"] = service.PagesService(page, pageSize, total)
	log.Printf("pages: %#v", resultData)

	return common.Result(code, common.GetMsg(code), resultData)
}

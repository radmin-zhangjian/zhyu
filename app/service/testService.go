package service

import (
	"log"
	"strconv"
	"time"
	"zhyu/app"
	"zhyu/app/dao"
)

type result struct {
	Id      int64
	Name    string
	Age     uint8
	Address string
	User    userInfo
}

type userInfo struct {
	Phone string
	Photo string
}

func SayRedisService(c *app.Context) any {
	a := c.Query("a")
	key := "say:" + a

	// set
	err := dao.SetString(key, "我是测试redis的setString!!"+a)
	if err != nil {
		return ""
	}

	// get
	res, err := dao.GetString(key)
	if err != nil {
		log.Println(err)
	}

	return res
}

func SayDbService(c *app.Context) any {

	Id, _ := strconv.ParseInt(c.Query("id"), 10, 64)

	if Id == 11 {
		time.Sleep(3 * time.Second)
	}

	log.Printf("query id: %v", Id)
	log.Printf("c.UserInfo.id: %v", c.UserInfo["id"])
	log.Printf("c.UserInfo.userName: %v", c.UserInfo["userName"])

	// 单条插入
	//a := c.Query("a")
	//user := model.User{Username: "小张", Password: "111222", Phone: "13133330000", Address: a}
	//dao.UserCreate(&user)
	//id := user.Id

	// 多条插入
	//user := []model.User{
	//	{Username: "go_yyds", Password: "111222", Phone: "13133330000", Address: "小岛"},
	//	{Username: "go_uuid", Password: "333666", Phone: "13133330001", Address: "小岛"},
	//}
	//dao.UserCreateBatch(&user)

	// ====== 原生查询 ====== todo
	sql := "SELECT * FROM zhyu_user WHERE phone = ? AND user_name = ? limit 10"
	info := dao.QueryRawDao(sql, "13133330001", "go_uuid")
	log.Printf("QueryExecDao: ========== %#v", info)

	// 获取第一个，默认查询第一个
	user, _ := dao.UserGetById(Id)
	log.Printf("user: ========== %v", user)

	// 获取第一个，默认查询第一个
	userData := dao.UserGetFirst()
	log.Printf("userData: %v", userData)
	//userDataJson, err := json.Marshal(userData)
	//if err == nil {
	//	fmt.Println(string(userDataJson))
	//} else {
	//	fmt.Println(err)
	//}

	// 列表
	userlist := dao.UserGetPage(3, 0)
	log.Printf("userlist: %v", userlist)

	// 更新
	dao.UserUpdateUsername(1, "张ss")

	// 拼装数据
	// 1、map类型
	resultData := map[string]any{
		"Id":      userData.Id,
		"Name":    userData.Name,
		"Age":     userData.Age,
		"Address": userData.Address,
		"User": map[string]any{
			"Phone": userData.Phone,
			"Photo": userData.Photo,
		},
	}
	resultData["xxx"] = "xxx"
	log.Printf("resultData: %v", resultData)
	// 2、普通 struct类型
	resultDataStruct := result{
		Id:      resultData["Id"].(int64),
		Name:    resultData["Name"].(string),
		Age:     userData.Age,
		Address: userData.Address,
		User: userInfo{
			Phone: userData.Phone,
			Photo: userData.Photo,
		},
	}
	log.Printf("resultDataStruct: %v", resultDataStruct)
	// 3、匿名 struct类型
	resultDataAny := struct {
		Id      int64
		Name    string
		Age     uint8
		Address string
		User    struct {
			Phone string
			Photo string
		}
	}{
		Id:      userData.Id,
		Name:    userData.Name,
		Age:     userData.Age,
		Address: userData.Address,
		User: struct {
			Phone string
			Photo string
		}{
			Phone: userData.Phone,
			Photo: userData.Photo,
		},
	}
	log.Printf("resultDataAny: %v", resultDataAny)
	// 4、匿名 struct类型2
	resultDataAny2 := struct {
		Id      int64
		Name    string
		Age     uint8
		Address string
		User    struct {
			Phone string
			Photo string
		}
		User2 any
	}{}
	resultDataAny2.Id = userData.Id
	resultDataAny2.Name = userData.Name
	resultDataAny2.Age = userData.Age
	resultDataAny2.Address = userData.Address
	resultDataAny2.User.Phone = userData.Phone
	resultDataAny2.User.Photo = userData.Photo
	resultDataAny2.User2 = resultDataAny2.User
	log.Printf("resultDataAny2: %v", resultDataAny2)

	return resultDataStruct
}

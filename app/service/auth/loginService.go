package auth

import (
	"context"
	"fmt"
	"github.com/go-playground/validator/v10"
	"log"
	"strconv"
	"strings"
	"zhyu/app"
	"zhyu/app/common"
	"zhyu/app/dao"
	"zhyu/app/model"
	"zhyu/utils"
)

type userAuth struct {
	Username string `from:"username" validate:"required,min=5,max=20"`
	Password string `from:"password" validate:"required,min=6,max=25"`
}

// LoginAuthService 用户登陆
func LoginAuthService(ctx context.Context, c *app.Context) any {

	// 创建返回信息
	data := make(map[string]interface{})
	code := common.INVALID_PARAMS

	username := c.PostForm("username")
	password := c.PostForm("password")

	// 验证参数
	valid := validator.New()
	a := userAuth{
		Username: username,
		Password: password,
	}
	ok := valid.Struct(a)
	if ok != nil {
		// 翻译成中文
		trans := utils.ValidateTransInit(valid)
		verrs := ok.(validator.ValidationErrors)
		errs := make(map[string]string)
		for key, value := range verrs.Translate(trans) {
			errs[key[strings.Index(key, ".")+1:]] = value
		}
		fmt.Println(errs)
		//fmt.Printf("Err(s):\n%+v\n", ok)
		return common.Result(code, common.GetMsg(code), data)
	}

	// 验证用户
	user, isExist := dao.UserGetOneNamePass(username, password)
	log.Printf("user: %#v,isExist: %#v ", user, isExist)
	if isExist {
		token, err := utils.GenerateToken(user.Id, username, password)
		if err != nil {
			code = common.ERROR_AUTH_TOKEN
		} else {
			data["token"] = token
			code = common.SUCCESS
		}
	} else {
		code = common.ERROR_AUTH
	}

	return common.Result(code, common.GetMsg(code), data)
}

type userRegister struct {
	Username  string `from:"username" validate:"required,min=5,max=20"`
	Password  string `from:"password" validate:"required,min=6,max=25"`
	Password2 string `from:"password2" validate:"required,min=6,max=25"`
	Age       string `from:"age"`
	Address   string `from:"address"`
}

// RegisterAuthService 用户注册
func RegisterAuthService(ctx context.Context, c *app.Context) any {

	// 创建返回信息
	data := make(map[string]interface{})
	code := common.INVALID_PARAMS

	// 接收参数
	var username string
	var password string
	var password2 string
	var age string
	var address string
	if c.GetHeader("Content-Type") == "application/json" {
		// 接收 application/json 类型的参数
		var paramJson userRegister
		if err := c.ShouldBindJSON(&paramJson); err != nil {
			return common.Result(code, common.GetMsg(code), data)
		}
		username = paramJson.Username
		password = paramJson.Password
		password2 = paramJson.Password2
		age = paramJson.Age
		address = paramJson.Address
	} else {
		// form 形式参数
		username = c.PostForm("username")
		password = c.PostForm("password")
		password2 = c.PostForm("password2")
		age = c.PostForm("age")
		address = c.PostForm("address")
	}

	// 验证参数
	valid := validator.New()
	a := userRegister{
		Username:  username,
		Password:  password,
		Password2: password2,
	}
	ok := valid.Struct(a)
	if ok != nil {
		// 翻译成中文
		trans := utils.ValidateTransInit(valid)
		verrs := ok.(validator.ValidationErrors)
		errs := make(map[string]string)
		for key, value := range verrs.Translate(trans) {
			errs[key[strings.Index(key, ".")+1:]] = value
		}
		fmt.Println(errs)
		return common.Result(code, common.GetMsg(code), data)
	}

	if password != password2 {
		return common.Result(common.ERROR_AUTH, common.GetMsg(common.ERROR_AUTH), data)
	}

	// 查询用户
	sql := "SELECT * FROM " + new(model.User).TableName() + " WHERE user_name = ? limit 1"
	info := dao.QueryRawDao(sql, username)
	if info != nil {
		return common.Result(common.ERROR, "账号已存在", data)
	}

	// 插入数据
	ageInt, _ := strconv.ParseInt(age, 10, 0)
	ageInt8 := uint8(ageInt)
	user := model.User{Username: username, Password: password, Phone: "", Age: ageInt8, Address: address}
	dao.UserCreate(&user)
	if user.Id == 0 {
		return common.Result(common.ERROR, "账号创建失败", data)
	}
	code = common.SUCCESS

	return common.Result(code, common.GetMsg(code), data)
}

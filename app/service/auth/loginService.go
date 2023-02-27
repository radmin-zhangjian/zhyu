package auth

import (
	"context"
	"fmt"
	"github.com/go-playground/validator/v10"
	"log"
	"strings"
	"zhyu/app"
	"zhyu/app/common"
	"zhyu/app/dao"
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

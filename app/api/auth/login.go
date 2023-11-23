package auth

import (
	"net/http"
	"zhyu/app/service/auth"
)

// Login 用户登陆
func (c *App) Login() {
	// valid

	// login auth
	data := auth.LoginAuthService(c.Ctx, c.Context)
	c.JSON(http.StatusOK, data)
}

// Register 用户注册
func (c *App) Register() {
	// valid

	// register
	data := auth.RegisterAuthService(c.Ctx, c.Context)
	c.JSON(http.StatusOK, data)
}

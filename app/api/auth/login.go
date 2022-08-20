package auth

import (
	"net/http"
	"zhyu/app/service/auth"
)

// Login 用户登陆
func (c *App) Login() {
	// valid

	// login auth
	data := auth.LoginAuth(c.Ctx, c.Context)
	c.JSON(http.StatusOK, data)
}

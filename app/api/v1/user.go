package v1

import (
	"net/http"
	v1 "zhyu/app/service/v1"
)

func (c *App) UserList() {
	data := v1.UserListService(c.Ctx, c.Context)
	c.JSON(http.StatusOK, data)
}

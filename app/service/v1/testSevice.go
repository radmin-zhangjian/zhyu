package v1

import (
	"zhyu/app"
)

func Say(c *app.Context) any {
	a := c.Query("a")
	return a
}

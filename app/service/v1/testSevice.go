package v1

import (
	"zhyu/app"
)

func SayService(c *app.Context) any {
	a := c.Query("a")
	return a
}

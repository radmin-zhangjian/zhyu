package v2

import (
	. "zhyu/app"
)

type App struct {
	*Context
}

func New() *App {
	return new(App)
}

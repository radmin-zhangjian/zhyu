package auth

import (
	. "zhyu/app"
)

type App struct {
	*Context
}

func New() *App {
	return new(App)
}

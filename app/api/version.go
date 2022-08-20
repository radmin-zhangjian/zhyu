package api

import (
	"zhyu/app/api/auth"
	v1 "zhyu/app/api/v1"
	v2 "zhyu/app/api/v2"
)

type AppVersion struct {
	Version string
	Object  interface{}
}

// Version 初始化版本
func Version() []AppVersion {
	return []AppVersion{
		{"auth", auth.New()},
		{"v1", v1.New()},
		{"v2", v2.New()},
	}
}

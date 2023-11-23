package common

var (
	// RouteAuth 不需要token认证的路由
	RouteAuth = map[string]struct{}{
		"/api/auth/login":      struct{}{},
		"/api/auth/register":   struct{}{},
		"/api/auth/middleware": struct{}{},
	}
)

package common

var (
	SUCCESS = 1000
	ERROR   = 1002
	FATAL   = 1004

	LOGOOUT                        = 1011
	ERROR_AUTH                     = 1020
	INVALID_PARAMS                 = 1021
	INVALID_RESULT                 = 1022
	ERROR_AUTH_TOKEN               = 1050
	ERROR_AUTH_NO_TOKRN            = 1051
	ERROR_AUTH_CHECK_TOKEN_FAIL    = 1052
	ERROR_AUTH_CHECK_TOKEN_TIMEOUT = 1053

	message = map[int]string{
		SUCCESS:                        "success",
		ERROR:                          "error",
		FATAL:                          "fatal",
		LOGOOUT:                        "退出登陆",
		ERROR_AUTH:                     "验证失败",
		INVALID_PARAMS:                 "无效的参数",
		INVALID_RESULT:                 "数据为空或nil",
		ERROR_AUTH_TOKEN:               "token error",
		ERROR_AUTH_NO_TOKRN:            "token为空",
		ERROR_AUTH_CHECK_TOKEN_FAIL:    "token认证失败",
		ERROR_AUTH_CHECK_TOKEN_TIMEOUT: "token过期",
	}
)

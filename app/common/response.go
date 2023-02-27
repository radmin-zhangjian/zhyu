package common

func GetMsg(key int) string {
	return message[key]
}

func Result(code int, message string, data any) map[string]any {
	if data == false || data == "" {
		data = map[string]any{}
	}
	return map[string]any{
		"code": code,
		"msg":  message,
		"data": data,
	}
}

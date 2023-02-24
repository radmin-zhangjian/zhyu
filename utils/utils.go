package utils

import (
	"zhyu/setting"
	"zhyu/utils/uuid"
)

var (
	Uuid, _ = uuid.NewSnowWorker(setting.Server.WorkerID)
)

// IsArray 查询数组是否存在某个值
func IsArray(target interface{}, obj []interface{}) bool {
	for i := 0; i < len(obj); i++ {
		if target == obj[i] {
			return true
		}
	}
	return false
}

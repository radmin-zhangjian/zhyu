package utils

import (
	"zhyu/setting"
	"zhyu/utils/uuid"
)

var (
	Uuid, _ = uuid.NewSnowWorker(setting.Server.WorkerID)
)

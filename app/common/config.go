package common

import (
	"sync"
)

var (
	// RequestIdMap utils.Uuid.GetId()
	RequestIdMap = new(sync.Map)

	// CallSecondaryService 降级开关 1=不降级 0=降级
	CallSecondaryService int = 1
	// SecondaryKey 降级开关的key
	SecondaryKey string = "/setting/reduce/rank"
)

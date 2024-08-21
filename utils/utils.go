package utils

import (
	"bytes"
	"encoding/json"
	"golang.org/x/crypto/bcrypt"
	"zhyu/setting"
	"zhyu/utils/logger"
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

// StructToString 结构体转字符串 包含空字段
func StructToString(item any) string {
	// 把struct转换成json byte
	emJson, err := json.Marshal(item)
	if err != nil {
		logger.Error("EsTest json marshal err:%v", err)
	}

	docData := string(emJson)

	return docData
}

// StructToMap 结构体转map 去除空字段
func StructToMap(item any) map[string]interface{} {
	// 把struct转换成json byte
	emJson, err := json.Marshal(item)
	if err != nil {
		logger.Error("EsTest json marshal err:%v", err)
	}

	var docJson map[string]interface{}
	var docData = make(map[string]interface{})
	d := json.NewDecoder(bytes.NewReader(emJson))
	d.UseNumber() // 设置将float64转为一个number
	if err := d.Decode(&docJson); err != nil {
		logger.Info("emJson err", err)
	} else {
		for k, v := range docJson {
			if v != "" {
				docData[k] = v
			}
		}
	}

	return docData
}

// HashPassword 加密密码
func HashPassword(password string) (string, error) {
	// bcrypt.GenerateFromPassword 使用 bcrypt 算法生成哈希值，第二个参数是工作因子（cost）
	// 10 是默认推荐的工作因子，数值越大，计算越耗时，安全性越高
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

// CheckPassword 验证密码
func CheckPassword(hashedPassword, password string) bool {
	// bcrypt.CompareHashAndPassword 会将加密的哈希密码和未加密的密码进行比较
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err == nil
}

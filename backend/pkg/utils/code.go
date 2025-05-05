package utils

import (
	"encoding/base64"
	"encoding/json"
	"gopkg.in/yaml.v3"
)

// IsYAML 判断是否为 yaml
func IsYAML(data string) bool {
	var yml map[string]interface{}
	return yaml.Unmarshal([]byte(data), &yml) == nil
}

// IsJSON 判断字符串是否为合法 JSON 格式
func IsJSON(data string) bool {
	var js json.RawMessage
	err := json.Unmarshal([]byte(data), &js)
	return err == nil
}

// IsBase64 判断字符串是否是 Base64 编码
func IsBase64(data string) bool {
	// 检查是否符合 RawStdEncoding 格式
	_, errRaw := base64.RawStdEncoding.DecodeString(data)
	if errRaw == nil {
		return true
	}

	// 检查是否符合 StdEncoding 格式
	_, errStd := base64.StdEncoding.DecodeString(data)
	if errStd == nil {
		return true
	}

	return false
}

package tools

import (
	"bytes"
	"crypto/aes"
	"encoding/base64"
	"errors"
)

// AesEncryptByECB 加密
func AesEncryptByECB(data, key string) (string, error) {
	// 判断key长度
	keyLenMap := map[int]struct{}{16: {}, 24: {}, 32: {}}
	if _, ok := keyLenMap[len(key)]; !ok {
		return "", errors.New("key长度必须是 16、24、32 其中一个")
	}
	// 密钥和待加密数据转成[]byte
	originByte := []byte(data)
	keyByte := []byte(key)
	// 创建密码组，长度只能是16、24、32字节
	block, _ := aes.NewCipher(keyByte)
	// 获取密钥长度
	blockSize := block.BlockSize()
	// 补码
	originByte = pkcs7padding(originByte, blockSize)
	// 创建保存加密变量
	encryptResult := make([]byte, len(originByte))
	// CEB是把整个明文分成若干段相同的小段，然后对每一小段进行加密
	for bs, be := 0, blockSize; bs < len(originByte); bs, be = bs+blockSize, be+blockSize {
		block.Encrypt(encryptResult[bs:be], originByte[bs:be])
	}
	return base64.StdEncoding.EncodeToString(encryptResult), nil
}

// pkcs7padding 补码
func pkcs7padding(originByte []byte, blockSize int) []byte {
	// 计算补码长度
	padding := blockSize - len(originByte)%blockSize
	// 生成补码
	padText := bytes.Repeat([]byte{byte(padding)}, padding)
	// 追加补码
	return append(originByte, padText...)
}

// AesDecryptByECB 解密
func AesDecryptByECB(data, key string) (string, error) {
	// 判断key长度
	keyLenMap := map[int]struct{}{16: {}, 24: {}, 32: {}}
	if _, ok := keyLenMap[len(key)]; !ok {
		return "", errors.New("key长度必须是 16、24、32 其中一个")
	}
	// 反解密码base64
	originByte, _ := base64.StdEncoding.DecodeString(data)
	// 密钥和待加密数据转成[]byte
	keyByte := []byte(key)
	// 创建密码组，长度只能是16、24、32字节
	block, _ := aes.NewCipher(keyByte)
	// 获取密钥长度
	blockSize := block.BlockSize()
	// 创建保存解密变量
	decrypted := make([]byte, len(originByte))
	for bs, be := 0, blockSize; bs < len(originByte); bs, be = bs+blockSize, be+blockSize {
		block.Decrypt(decrypted[bs:be], originByte[bs:be])
	}
	// 解码
	return string(pkcs7UnPadding(decrypted)), nil
}

// pkcs7UnPadding 解码
func pkcs7UnPadding(originDataByte []byte) []byte {
	length := len(originDataByte)
	unfading := int(originDataByte[length-1])
	return originDataByte[:(length - unfading)]
}

package tools

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha1"
	"crypto/x509"
	"encoding/pem"
)

// RsaOAEPEncrypt 加密
func RsaOAEPEncrypt(origData []byte, publicKey string) ([]byte, error) {
	block, _ := pem.Decode([]byte(publicKey))
	pub, _ := x509.ParsePKCS1PublicKey(block.Bytes)
	return rsa.EncryptOAEP(sha1.New(), rand.Reader, pub, origData, nil)
}

// RsaOAEPDecrypt 解密
func RsaOAEPDecrypt(ciphertext []byte, privateKey string) ([]byte, error) {
	block, _ := pem.Decode([]byte(privateKey))
	riv, _ := x509.ParsePKCS1PrivateKey(block.Bytes)
	return rsa.DecryptOAEP(sha1.New(), rand.Reader, riv, ciphertext, nil)
}

// RsaEncrypt 加密
func RsaEncrypt(origData []byte, publicKey []byte) ([]byte, error) {
	block, _ := pem.Decode(publicKey)
	pub, _ := x509.ParsePKCS1PublicKey(block.Bytes)
	return rsa.EncryptOAEP(sha1.New(), rand.Reader, pub, origData, nil)
}

// RsaDecrypt 解密
func RsaDecrypt(ciphertext []byte, privateKey []byte) ([]byte, error) {
	block, _ := pem.Decode(privateKey)
	riv, _ := x509.ParsePKCS1PrivateKey(block.Bytes)
	return rsa.DecryptOAEP(sha1.New(), rand.Reader, riv, ciphertext, nil)
}

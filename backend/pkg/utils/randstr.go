package utils

import (
	"bytes"
	"crypto/md5"
	"crypto/rand"
	"encoding/binary"
	"encoding/hex"
)

// RandBytes generates n random bytes
func RandBytes(n int) []byte {
	b := make([]byte, n)
	_, err := rand.Read(b)
	if err != nil {
		panic(err)
	}
	return b
}

const Base64Chars = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ+/"
const Base62Chars = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
const HexChars = "0123456789abcdef"
const DecChars = "0123456789"

// RandBase64 generates a random Base64 string with length of n
//
// Example: X02+jDDF/exDoqPg9/aXlzbUCN93GIQ5
func RandBase64(n int) string { return RandString(n, Base64Chars) }

// RandBase62 generates a random Base62 string with length of n
//
// Example: 1BsNqB61o4ztSqLC6labKGNf4MYy352X
func RandBase62(n int) string { return RandString(n, Base62Chars) }

// RandDec generates a random decimal number string with length of n
//
// Example: 37110235710860781655802098192113
func RandDec(n int) string { return RandString(n, DecChars) }

// RandHex generates a random Hexadecimal string with length of n
//
// Example: 67aab2d956bd7cc621af22cfb169cba8
func RandHex(n int) string { return RandString(n, HexChars) }

// list of default letters that can be used to make a random string when calling String
// function with no letters provided
var defLetters = []rune(Base62Chars)

// RandString generates a random string using only letters provided in the letters parameter.
//
// If user omits letters parameter, this function will use Base62Chars instead.
func RandString(n int, letters ...string) string {
	var letterRunes []rune
	if len(letters) == 0 {
		letterRunes = defLetters
	} else {
		letterRunes = []rune(letters[0])
	}

	var bb bytes.Buffer
	bb.Grow(n)
	l := uint32(len(letterRunes))
	// on each loop, generate one random rune and append to output
	for i := 0; i < n; i++ {
		bb.WriteRune(letterRunes[binary.BigEndian.Uint32(RandBytes(4))%l])
	}
	return bb.String()
}

// MD5 计算md5
func MD5(v string) string {
	d := []byte(v)
	m := md5.New()
	m.Write(d)
	return hex.EncodeToString(m.Sum(nil))
}

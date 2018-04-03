package util

import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/hex"
	"fmt"
	"math/rand"
	"strings"
)

/*
	generate encrypted password
*/
func SHA512(text string) string {
	hasher := sha512.New()
	hasher.Write([]byte(text))
	return strings.ToUpper(hex.EncodeToString(hasher.Sum(nil)))
}

func GenerateSessionId(seed string) string {
	sid := sha256.Sum256([]byte(seed + fmt.Sprint(rand.Int())))
	return strings.ToUpper(hex.EncodeToString(sid[:]))
}

func Md5Hash(text string) string {
	hasher := md5.New()
	hasher.Write([]byte(text))
	return hex.EncodeToString(hasher.Sum(nil))
}

func Sha1Encrypt(str string) string {
	hash := sha1.New()
	hash.Write([]byte(str))
	md := hash.Sum(nil)
	return hex.EncodeToString(md)
}

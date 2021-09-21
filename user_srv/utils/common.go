package utils

import (
	"crypto/md5"
	"encoding/hex"

	"user_srv/global"
)

func MD5(str string) string {
	salt := global.Config.MD5Salt
	hash := md5.New()
	hash.Write([]byte(str + salt))
	return hex.EncodeToString(hash.Sum(nil))
}

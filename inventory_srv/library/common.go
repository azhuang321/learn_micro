package library

import (
	"crypto/md5"
	"encoding/hex"
	"github.com/gookit/ini/v2"
)

func MD5(str string) string {
	salt := ini.String("MD5_SALT")
	hash := md5.New()
	hash.Write([]byte(str + salt))
	return hex.EncodeToString(hash.Sum(nil))
}

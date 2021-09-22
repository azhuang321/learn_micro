package initialize

import (
	"github.com/gookit/validate"
	"github.com/gookit/validate/locales/zhcn"
)

func InitTranslate() {
	zhcn.RegisterGlobal()

	validate.AddGlobalMessages(map[string]string{
		"len":      "{field}长度为%d",
		"minLen":   "{field}最小长度为%d",
		"maxLen":   "{field}最大长度为%d",
		"required": "{field}必须填写",
	})
}

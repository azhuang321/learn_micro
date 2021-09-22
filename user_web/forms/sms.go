package forms

import (
	"github.com/gookit/validate"
	"regexp"
)

type SendSmsForm struct {
	Mobile string `form:"mobile" json:"mobile" validate:"required|MobileValidate"`
	Type   uint   `form:"type" json:"type" validate:"required"`
}

// Messages 您可以自定义验证器错误消息
func (s SendSmsForm) Messages() map[string]string {
	return validate.MS{
		"Mobile.MobileValidate": "{field}输入不正确",
	}
}

// Translates 你可以自定义字段翻译
func (s SendSmsForm) Translates() map[string]string {
	return validate.MS{
		"Mobile": "手机号码",
		"Type":   "密码",
	}
}

// MobileValidate 定义在结构体中的自定义验证器
func (s SendSmsForm) MobileValidate(val string) bool {
	ok, _ := regexp.MatchString(`^(13[0-9]|14[5|7]|15[0|1|2|3|4|5|6|7|8|9]|18[0|1|2|3|5|6|7|8|9])\d{8}$`, val)
	if !ok {
		return false
	}
	return true
}

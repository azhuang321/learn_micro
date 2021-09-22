package user

import (
	"regexp"

	"github.com/gookit/validate"
)

type RegisterFrom struct {
	Mobile   string `json:"mobile" validate:"required|MobileValidate"`
	Password string `json:"password" validate:"required|minLen:3|maxLen:20"`
	Code     string `json:"code" validate:"required|len:6"`
}

// Messages 您可以自定义验证器错误消息
func (f RegisterFrom) Messages() map[string]string {
	return validate.MS{
		"Mobile.MobileValidate": "{field}输入不正确",
	}
}

// Translates 你可以自定义字段翻译
func (f RegisterFrom) Translates() map[string]string {
	return validate.MS{
		"Mobile":   "手机号码",
		"Password": "密码",
		"Code":     "验证码",
	}
}

func (f RegisterFrom) MobileValidate(val string) bool {
	ok, _ := regexp.MatchString(`^(13[0-9]|14[5|7]|15[0|1|2|3|4|5|6|7|8|9]|18[0|1|2|3|5|6|7|8|9])\d{8}$`, val)
	if !ok {
		return false
	}
	return true
}

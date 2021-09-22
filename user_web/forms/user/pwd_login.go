package user

import (
	"github.com/gookit/validate"
	"regexp"
)

type PasswordLoginFrom struct {
	Mobile    string `json:"mobile" validate:"required|MobileValidate"`
	Password  string `json:"password" validate:"required|minLen:3|maxLen:20"`
	CaptchaId string `json:"captcha_id" validate:"required"`
	Captcha   string `json:"captcha" validate:"required|len:5"`
}

// Messages 您可以自定义验证器错误消息
func (f PasswordLoginFrom) Messages() map[string]string {
	return validate.MS{
		"Mobile.MobileValidate": "{field}输入不正确",
	}
}

// Translates 你可以自定义字段翻译
func (f PasswordLoginFrom) Translates() map[string]string {
	return validate.MS{
		"Mobile":    "手机号码",
		"Password":  "密码",
		"CaptchaId": "验证码ID",
		"Captcha":   "验证码",
	}
}

// MobileValidate 定义在结构体中的自定义验证器
func (f PasswordLoginFrom) MobileValidate(val string) bool {
	ok, _ := regexp.MatchString(`^(13[0-9]|14[5|7]|15[0|1|2|3|4|5|6|7|8|9]|18[0|1|2|3|5|6|7|8|9])\d{8}$`, val)
	if !ok {
		return false
	}
	return true
}

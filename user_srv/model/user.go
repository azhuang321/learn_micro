package model

type User struct {
	ID       int64  `json:"id"`
	Mobile   string `json:"mobile"`
	Password string `json:"password"`
	Nickname string `json:"nick_name"`
}

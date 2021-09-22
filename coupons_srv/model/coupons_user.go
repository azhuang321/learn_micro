package model

type CouponsUser struct {
	ID     uint   `json:"id" gorm:"column:id"`
	Mobile string `json:"mobile" gorm:"column:mobile"`
	Num    uint32 `json:"num" gorm:"column:num"`
}

func (m *CouponsUser) TableName() string {
	return "coupons_user"
}

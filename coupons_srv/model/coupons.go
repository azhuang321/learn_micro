package model

type Coupons struct {
	ID  uint   `json:"id" gorm:"column:id"`
	Num uint32 `json:"num" gorm:"column:num"`
}

func (m *Coupons) TableName() string {
	return "coupons"
}

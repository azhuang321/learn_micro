package model

type CouponsHistory struct {
	ID     uint   `json:"id" gorm:"column:id"`
	Mobile string `json:"mobile" gorm:"column:mobile"`
	Num    uint32 `json:"num" gorm:"column:num"`
	Status int8   `json:"status" gorm:"column:status"` // (1:未扣减,2:已归还)
}

func (m *CouponsHistory) TableName() string {
	return "coupons_history"
}

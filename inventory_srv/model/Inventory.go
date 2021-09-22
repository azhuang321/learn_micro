package model

type Inventory struct {
	ID         uint `json:"id" gorm:"column:id"`
	AddTime    int  `json:"add_time" gorm:"column:add_time"`
	IsDeteted  int8 `json:"is_deteted" gorm:"column:is_deteted"`
	UpdateTime int  `json:"update_time" gorm:"column:update_time"`
	Goods      int  `json:"goods" gorm:"column:goods"`     // 商品id
	Stocks     int  `json:"stocks" gorm:"column:stocks"`   // 库存数量
	Version    int  `json:"version" gorm:"column:version"` // 版本号(分布式锁的乐观锁)
}

func (m *Inventory) TableName() string {
	return "inventory"
}

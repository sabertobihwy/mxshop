package model

import (
	"database/sql/driver"
	"github.com/goccy/go-json"
)

type GoodsInfo struct {
	GoodsId int32
	Num     int32
}

type GooodsInfoList []GoodsInfo

func (g *GooodsInfoList) Scan(value interface{}) error {
	return json.Unmarshal(value.([]byte), &g)
}
func (g GooodsInfoList) Value() (driver.Value, error) {
	return json.Marshal(g)
}

type Inventory struct {
	BaseModel
	Goods  int32 `gorm:"type:int;index"`
	Stocks int32 `gorm:"type:int;"`
	Verson int32 `gorm:"type:int;"`
}

type InventoryHistory struct {
	OrderSn string         `gorm:"type:varchar(20);index:idx_order_sn,unique"`
	Status  int32          `gorm:"type:int;"` // 1. refduced 2. rebacked
	Details GooodsInfoList `gorm:"type:varchar(200);"`
}

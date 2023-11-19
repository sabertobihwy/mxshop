package model

type Inventory struct {
	BaseModel
	Goods  int32 `gorm:"type:int;index"`
	Stocks int32 `gorm:"type:int;"`
	Verson int32 `gorm:"type:int;"`
}

//type InventoryHistory struct {
//	user   int32
//	goods  int32
//	nums   int32
//	order  int32
//	status int32 // 1: withholding stock, idempotence 2: paid
//	// 如果要归还stock，需要先查看有没有归还stock的流水，避免重复
//}

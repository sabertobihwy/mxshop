package model

import (
	"database/sql/driver"
	"github.com/goccy/go-json"
	"time"
)
import "gorm.io/gorm"

type BaseModel struct {
	ID        int32     `gorm:"primarykey"`
	CreatedAt time.Time `gorm:"column:add_time"`
	UpdatedAt time.Time `gorm:"column:update_time"`
	DeletedAt gorm.DeletedAt
	IsDeleted bool
}
type GormList []string

func (g *GormList) Scan(value interface{}) error {
	return json.Unmarshal(value.([]byte), &g)
}
func (g GormList) Value() (driver.Value, error) {
	return json.Marshal(g)
}

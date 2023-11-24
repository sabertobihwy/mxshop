package model

import (
	"database/sql/driver"
	"github.com/goccy/go-json"
	"time"
)
import "gorm.io/gorm"

type BaseModel struct {
	ID        int32          `gorm:"primarykey" json:"id"`
	CreatedAt time.Time      `gorm:"column:add_time" json:"-"`
	UpdatedAt time.Time      `gorm:"column:update_time" json:"-"`
	DeletedAt gorm.DeletedAt `json:"-"`
	IsDeleted bool           `json:"-"`
}
type GormList []string

func (g *GormList) Scan(value interface{}) error {
	return json.Unmarshal(value.([]byte), &g)
}
func (g GormList) Value() (driver.Value, error) {
	return json.Marshal(g)
}

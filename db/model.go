package db

import (
	"fmt"
	"github.com/ebar-go/ego/library"
	"database/sql/driver"
	"time"
)

// 基础模型
type Model struct {
	CreatedAt Timestamp `gorm:"column:createtime" json:"created_at"`
	UpdatedAt Timestamp `gorm:"column:updatetime" json:"updated_at"`
}

// Timestamp 满足gorm的时间戳格式
type Timestamp struct {
	time.Time
}

// MarshalJSON 解析
func (t Timestamp) MarshalJSON() ([]byte, error) {
	formatted := fmt.Sprintf("\"%s\"", t.Format(library.GetDefaultTimeFormat()))
	return []byte(formatted), nil
}

// Value
func (t Timestamp) Value() (driver.Value, error) {
	var zeroTime time.Time
	if t.Time.UnixNano() == zeroTime.UnixNano() {
		return nil, nil
	}
	return t.Time, nil
}

// Scan 转换时间戳
func (t *Timestamp) Scan(v interface{}) error {
	value, ok := v.(time.Time)
	if ok {
		*t = Timestamp{Time: value}
		return nil
	}
	return fmt.Errorf("can not convert %v to timestamp", v)
}

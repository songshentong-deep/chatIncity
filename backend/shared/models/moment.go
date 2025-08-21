package models

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"time"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// StringArray 自定义类型用于处理 PostgreSQL JSONB 数组
type StringArray []string

// Value 实现 driver.Valuer 接口
func (sa StringArray) Value() (driver.Value, error) {
	if len(sa) == 0 {
		return "[]", nil
	}
	return json.Marshal(sa)
}

// Scan 实现 sql.Scanner 接口
func (sa *StringArray) Scan(value interface{}) error {
	if value == nil {
		*sa = StringArray{}
		return nil
	}
	
	var bytes []byte
	switch v := value.(type) {
	case []byte:
		bytes = v
	case string:
		bytes = []byte(v)
	default:
		return fmt.Errorf("cannot scan %T into StringArray", value)
	}
	
	return json.Unmarshal(bytes, sa)
}

type Moment struct {
	ID        string         `json:"id" gorm:"primaryKey;type:varchar(36)"`
	UserID    string         `json:"user_id" gorm:"type:varchar(36);not null;index"`
	Content   string         `json:"content" gorm:"type:text"`
	Images    StringArray    `json:"images" gorm:"type:jsonb"`
	VideoURL  *string        `json:"video_url" gorm:"type:varchar(500)"`
	Type      string         `json:"type" gorm:"type:varchar(20);default:'text'"` // text, image, video
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
	
	// 关联用户信息
	User User `json:"user" gorm:"foreignKey:UserID"`
}

func (m *Moment) BeforeCreate(tx *gorm.DB) error {
	if m.ID == "" {
		m.ID = uuid.New().String()
	}
	return nil
}
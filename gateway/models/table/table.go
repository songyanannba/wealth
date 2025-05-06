package table

import "time"

type GVA_MODEL struct {
	ID        int        `gorm:"primarykey"`               // 主键ID
	CreatedAt *time.Time `gorm:"column:created_at;size:0"` // 创建时间
	UpdatedAt *time.Time `gorm:"column:updated_at;size:0"` // 更新时间
}

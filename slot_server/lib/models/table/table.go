package table

type GVA_MODEL struct {
	ID        int   `gorm:"primarykey"`               // 主键ID
	CreatedAt int64 `gorm:"column:createtime;size:0"` // 创建时间
	UpdatedAt int64 `gorm:"column:updatetime;size:0"` // 更新时间
}

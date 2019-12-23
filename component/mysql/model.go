package mysql

// 基础模型
type Model struct {
	// 这里的time无法直接通过json.Unmarshal解析为实体，因为时间格式不含T
	CreatedAt Timestamp `gorm:"column:createtime" json:"created_at"`
	UpdatedAt Timestamp `gorm:"column:updatetime" json:"updated_at"`
}

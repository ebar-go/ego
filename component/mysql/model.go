package mysql

// 基础模型
type Model struct {
	// 这里的time无法直接通过json.Unmarshal解析为实体，因为时间格式不含T
	// TODO 更换为时间戳，因为timestamp会涉及时区问题，对于跨地区应用，时间戳更为恰当
	CreatedAt Timestamp `gorm:"column:createtime" json:"created_at"`
	UpdatedAt Timestamp `gorm:"column:updatetime" json:"updated_at"`
}

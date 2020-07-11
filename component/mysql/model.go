package mysql

// 基础模型
type Model struct {
	Id        int `gorm:"primary_key;AUTO_INCREMENT;column:id" json:"id"`
	CreatedAt int `gorm:"column:created_at" json:"created_at"`
	UpdatedAt int `gorm:"column:updated_at" json:"updated_at"`
}

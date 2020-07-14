// Dao database access object
// usage
/**
type UserDao struct {
	mysql.Dao
}
 */
package mysql

import "github.com/jinzhu/gorm"

// Dao
type Dao struct {
	DB *gorm.DB
}

// Create
func (dao *Dao) Create(entity interface{}) error {
	return dao.DB.Create(entity).Error
}

// Update
func (dao *Dao) Update(entity interface{}) error {
	return dao.DB.Save(entity).Error
}

// GetByPrimaryKey
func (dao *Dao) GetByPrimaryKey(id int, entity interface{}) error  {
	return dao.DB.First(entity, id).Error
}


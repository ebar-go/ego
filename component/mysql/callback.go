// callback register gorm callback event
// usage:
/*
func init()  {
	event.Listen(event.AfterDatabaseConnect, func(ev event.Event) {
		mysql.SetCreateCallback(app.DB())
		mysql.SetUpdateCallback(app.DB())
	})
}
 */
package mysql

import (
	"github.com/ebar-go/ego/utils/date"
	"github.com/jinzhu/gorm"
)


// SetCreateCallback
func SetCreateCallback(db *gorm.DB)  {
	db.Callback().Create().Register("update_created_at", func(scope *gorm.Scope) {
		now := date.GetTimeStamp()
		if scope.HasColumn(ColumnCreatedAt) {
			_ = scope.SetColumn(ColumnCreatedAt, now)
		}
		if scope.HasColumn(ColumnUpdatedAt) {
			_ = scope.SetColumn(ColumnUpdatedAt, now)
		}
	})
}

// SetUpdateCallback
func SetUpdateCallback(db *gorm.DB)  {
	db.Callback().Update().Register("update_updated_at", func(scope *gorm.Scope) {
		if scope.HasColumn(ColumnUpdatedAt) {
			_ = scope.SetColumn(ColumnUpdatedAt, date.GetTimeStamp())
		}
	})
}

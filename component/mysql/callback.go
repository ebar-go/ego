// callback register gorm callback event
// usage:
/*
func init()  {
	event.Listen(event.AfterDatabaseConnect, func(ev event.Event) {
		mysql.RegisterCreateCallback(app.DB())
		mysql.RegisterUpdateCallback(app.DB())
	})
}
 */
package mysql

import (
	"github.com/ebar-go/egu"
	"github.com/jinzhu/gorm"
)


// RegisterCreateCallback
func RegisterCreateCallback(db *gorm.DB)  {
	db.Callback().Create().Register("update_created_at", func(scope *gorm.Scope) {
		now := egu.GetTimeStamp()
		if scope.HasColumn(ColumnCreatedAt) {
			_ = scope.SetColumn(ColumnCreatedAt, now)
		}
		if scope.HasColumn(ColumnUpdatedAt) {
			_ = scope.SetColumn(ColumnUpdatedAt, now)
		}
	})
}

// RegisterUpdateCallback
func RegisterUpdateCallback(db *gorm.DB)  {
	db.Callback().Update().Register("update_updated_at", func(scope *gorm.Scope) {
		if scope.HasColumn(ColumnUpdatedAt) {
			_ = scope.SetColumn(ColumnUpdatedAt, egu.GetTimeStamp())
		}
	})
}

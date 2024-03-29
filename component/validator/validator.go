package validator

import (
	"github.com/go-playground/validator/v10"
	"reflect"
	"sync"
)

type Instance struct {
	once     sync.Once
	validate *validator.Validate
}

// getKindOf return the kind of data
func getKindOf(data interface{}) reflect.Kind {
	value := reflect.ValueOf(data)
	valueType := value.Kind()

	if valueType == reflect.Ptr {
		valueType = value.Elem().Kind()
	}
	return valueType
}

// ValidateStruct validate struct
func (v *Instance) ValidateStruct(obj interface{}) error {
	if getKindOf(obj) == reflect.Struct {
		v.lazyInit()
		return v.validate.Struct(obj)
	}

	return nil
}

// lazyInit
func (v *Instance) lazyInit() {
	v.once.Do(func() {
		v.validate.SetTagName("binding")

		// define filed name
		v.validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
			return fld.Tag.Get("comment")
		})
	})
}

func New() *Instance {
	return &Instance{validate: validator.New()}
}

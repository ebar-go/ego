package validator

import (
	"reflect"
	"sync"

	"errors"
	"github.com/ebar-go/ego/helper"
	"github.com/go-playground/locales/zh"
	zh_translations "github.com/go-playground/validator/v10/translations/zh"
)

// trans use single pattern
var trans = GetZhTranslator()

func GetZhTranslator() ut.Translator {
	//中文翻译器
	zh_ch := zh.New()
	uni := ut.New(zh_ch)
	trans, _ := uni.GetTranslator("zh")
	return trans
}

// Validator 自定义验证器
type Validator struct {
	once     sync.Once
	validate *validator.Validate
}

// ValidateStruct 验证
func (v *Validator) ValidateStruct(obj interface{}) error {

	if helper.GetKindOf(obj) == reflect.Struct {
		v.lazyInit()

		if err := v.validate.Struct(obj); err != nil {
			//验证器
			for _, err := range err.(validator.ValidationErrors) {
				return errors.New(err.Translate(trans))
			}

		}
	}

	return nil
}

// Engine 获取一个实例
func (v *Validator) Engine() interface{} {
	v.lazyInit()
	return v.validate
}

// lazyInit 懒加载
func (v *Validator) lazyInit() {
	v.once.Do(func() {
		v.validate = validator.New()
		v.validate.SetTagName("binding")

		// 自定义名称字段
		v.validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
			return fld.Tag.Get("comment")
		})

		// 使用中文
		_ = zh_translations.RegisterDefaultTranslations(v.validate, trans)
	})
}

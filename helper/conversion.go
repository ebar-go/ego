package helper

import (
	"errors"
	"reflect"
	json "github.com/pquerna/ffjson/ffjson"
	"fmt"
)

// Struct2Map 支持结构体转化为map，在嵌套结构中不支持interface{}传值的结构体
func Struct2Map(obj interface{}) map[string]interface{} {
	var node map[string]interface{}
	objT := reflect.TypeOf(obj)
	if objT.Kind() != reflect.Struct {
		panic(errors.New("argument is not of the expected type"))
	}
	objV := reflect.ValueOf(obj)
	var data = make(map[string]interface{})
	for i := 0; i < objT.NumField(); i++ {
		switch objV.Field(i).Type().Kind() {
		case reflect.Struct:
			node = Struct2Map(objV.Field(i).Interface())
			data[objT.Field(i).Name] = node
		case reflect.Slice:
			target := objV.Field(i).Interface()
			tmp := make([]map[string]interface{}, reflect.ValueOf(target).Len())
			for j := 0; j < reflect.ValueOf(target).Len(); j++ {
				if reflect.ValueOf(target).Index(j).Kind() == reflect.Struct {
					node = Struct2Map(reflect.ValueOf(target).Index(j).Interface())
					tmp[j] = node
				}
			}
			data[objT.Field(i).Name] = tmp
		default:
			data[objT.Field(i).Name] = objV.Field(i).Interface()
		}
	}
	return data
}

// Map2Struct 数组转结构体
func Map2Struct(mapInstance map[string]interface{}, obj interface{}) error {
	buf, err := json.Marshal(mapInstance)
	if err != nil {
		return err
	}

	if err := json.Unmarshal(buf, obj); err!= nil {
		return err
	}

	return nil
}

// Float2String float转string
func Float2String(a float64) string {
	return fmt.Sprintf("%.f", a)
}

package conv

import (
	"encoding/json"
	"errors"
	"reflect"
	"strconv"
)

// Struct2Map return map
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

// Map2Struct return interface
func Map2Struct(mapInstance map[string]interface{}, obj interface{}) error {
	buf, err := json.Marshal(mapInstance)
	if err != nil {
		return err
	}

	if err := json.Unmarshal(buf, obj); err != nil {
		return err
	}

	return nil
}

// TransformInterface transform the source interface to target,target should be a pointer
func TransformInterface(source interface{}, target interface{}) error {
	buf, err := json.Marshal(source)
	if err != nil {
		return err
	}

	return json.Unmarshal(buf, target)
}

// String2Int return int of number string
func String2Int(s string) int {
	i, err := strconv.Atoi(s)
	if err != nil {
		return 0
	}

	return i
}

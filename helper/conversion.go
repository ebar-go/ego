package helper

import (
	"errors"
	"reflect"
	json "github.com/pquerna/ffjson/ffjson"
	"fmt"
	"strconv"
	"net/http"
	"io/ioutil"
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

	if err := json.Unmarshal(buf, obj); err != nil {
		return err
	}

	return nil
}

// FormatInterface 格式化interface
func FormatInterface(source interface{}, target interface{}) error {
	jsonStr, err := JsonEncode(source)
	if err != nil {
		return err
	}

	return JsonDecode([]byte(jsonStr), target)
}

// Float2String float转string
func Float2String(a float64) string {
	return fmt.Sprintf("%.f", a)
}

// Interface2Int interface转int
func Interface2Int(i interface{}) int {
	f, ok := i.(float64)
	if !ok {
		return 0
	}

	str := fmt.Sprintf("%.f", f)
	if result, err := strconv.Atoi(str); err == nil {
		return result
	}

	return 0
}

// StringifyResponse 将response序列化
func StringifyResponse(response *http.Response) (string, error) {
	if response == nil {
		return "", errors.New("没有响应数据")
	}

	if response.StatusCode != 200 {
		return "", errors.New("非200的上游返回")
	}

	data, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return "", err
	}

	// 关闭响应
	defer func() {
		response.Body.Close()
	}()

	return string(data), nil
}

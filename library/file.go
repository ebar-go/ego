package library

import (
	"os"
	"fmt"
)

//IsPathExist 判断文件或文件夹是否存在
func IsPathExist(path string)(bool){
	_, err := os.Stat(path)
	if err != nil{
		if os.IsExist(err){
			return true
		}
		if os.IsNotExist(err){
			return false
		}
		fmt.Println(err)
		return false
	}
	return true
}
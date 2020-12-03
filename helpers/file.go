package helpers

import (
	"fmt"
	"io/ioutil"
	"os"
)

//读取整个文件
func ReadJsonFile(filePth string) []byte {
	data, err := ioutil.ReadFile(filePth)
	if err != nil {
		fmt.Println("error" + err.Error())
	}
	return data
}

//判断文件夹是否存在
func FileIsExist(path string) bool {
	_, err := os.Stat(path)
	if err != nil {
		if os.IsExist(err) {
			return true
		}
		if os.IsNotExist(err) {
			return false
		}
		fmt.Println(err)
		return false
	}
	return true
}

//递归创建文件夹
func RecursiveCreationFolder(path string, perm os.FileMode) error {
	if !FileIsExist(path) {
		err := os.MkdirAll(path, perm)
		return err
	}
	return nil
}

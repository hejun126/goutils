package helpers

import (
	"errors"
	"io"
	"mime/multipart"
	"os"
	"strings"
	"time"
)

//文件路径
var path = "./storage"
//允许的后缀
var allowSuffix = [] string{"jpg", "gif","jpeg", "png", "wbmp"}
//允许的mime类型
var allowMime = []string{"image/jpg", "image/jpeg", "image/gif", "image/wbmp", "image/png"}
//最大尺寸
var maxSize int64 = 2 * 1024 * 1024
//文件名前缀
var prefix = "up_"

//图片上传结构体
type Image struct {
	//上传路径
	Path string
	//文件头信息
	FileHeader *multipart.FileHeader
	//文件信息
	File multipart.File
}

//上传
func (image *Image) Upload() (string, error) {
	//验证
	mine := image.FileHeader.Header["Content-Type"][0]
	suffix := strings.Split(mine, "image/")[1]
	bse, err := checkSize(image.FileHeader.Size)
	if !bse {
		return "", err
	}
	bsx, err := checkSuffix(suffix)
	if !bsx {
		return "", err
	}
	bme, err := checkMine(mine)
	if !bme {
		return "", err
	}
	//文件夹创建
	b, truePath, err := checkFolder(image.Path)
	if !b {
		return "", err
	}
	newName := createNewName(suffix)
	//保存移动
	savePath := truePath + "/" + newName
	out, err := os.Create(savePath)
	if err != nil {
		return "", errors.New("图片上传错误")
	}
	_, err = io.Copy(out, image.File)
	if err != nil {
		return "", errors.New("图片移动错误")
	}
	return savePath[1:], nil
}

//检查后缀
func checkSuffix(suffix string) (bool, error) {
	index := InArray(allowSuffix, suffix)
	if index > -1 {
		return true, nil
	}
	err := errors.New("不允许的图片类型")
	return false, err
}

//检查mine类型
func checkMine(mine string) (bool, error) {
	index := InArray(allowMime, mine)
	if index > -1 {
		return true, nil
	}
	err := errors.New("不允许的mine类型")
	return false, err
}

//检测尺寸
func checkSize(size int64) (bool, error) {
	if size > maxSize {
		err := errors.New("允许上传的文件最大为2M")
		return false, err
	}
	return true, nil
}

//创建新名称
func createNewName(suffix string) string {
	var newName string
	uniqueId := GetUUID()
	newName = prefix + uniqueId + "." + suffix
	return newName
}

//检测文件夹是否存在，若不存在则创建
func checkFolder(requestPath string) (bool, string, error) {
	currentTime := time.Now()
	truePath := path + "/" + requestPath + "/" + currentTime.Format("2006") + "/" + time.Now().Format("01") + "/" + currentTime.Format("02")
	if FileIsExist(truePath) {
		return true, truePath, nil
	}
	err := RecursiveCreationFolder(truePath, 0777)
	if err != nil {
		err := errors.New("文件系统错误")
		return false, "", err
	}
	return true, truePath, nil
}

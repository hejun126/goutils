package helpers

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io"
)

//生产32位的md5字符串
func Md5V1(str string) string {
	h := md5.New()
	h.Write([]byte(str))
	return hex.EncodeToString(h.Sum(nil))
}



func Md5V2(str string) string {
	data := []byte(str)
	md5Result := md5.Sum(data)
	return fmt.Sprintf("%x", md5Result)
}

func Md5V3(str string) string {
	h := md5.New()

	io.WriteString(h, str)
	return fmt.Sprintf("%x", h.Sum(nil))
}
package helpers

import (
	"bytes"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"net/http"
	"strconv"
)

//计算表名
func GetTableIndex(keyword string) (int,error){
	md5Str := Md5V2(keyword)
	index, err  := strconv.ParseUint(md5Str[28:], 16, 32)
	if err != nil {
		return 0, err
	}
	return int(index % 100), nil
}

//唯一标识
func UniqueId()(string,error){
	b := make([]byte, 48)
	if _, err := io.ReadFull(rand.Reader, b); err != nil {
		return "",errors.New("唯一标识创建失败")
	}
	return Md5V1(base64.URLEncoding.EncodeToString(b)),nil
}
//post方式请求接口
func PostJson(url string,data interface{}) (string,error){
	content := ""
	jsonStr,err := json.Marshal(data)
	if err != nil {
		return content,err
	}
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonStr))
	if err != nil{
		return content,err
	}
	defer resp.Body.Close()
	tmp,err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return content,err
	}
	return string(tmp),nil
}



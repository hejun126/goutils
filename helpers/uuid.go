package helpers

import (
	"fmt"
	"math/rand"
	"time"
)

//读取整个文件
func GetUUID() string {
	rand.Seed(time.Now().UnixNano())
	randNumberA := rand.Intn(10000)
	randNumberB := rand.Intn(10000)
	uuid := fmt.Sprintf("%d%d%d",time.Now().Unix(),randNumberA, randNumberB)
	return uuid
}
package helpers

import "time"

func GetNowTime() string{
	var cstZone = time.FixedZone("CST", 8*3600)       // 东八
	return time.Now().In(cstZone).Format("2006-01-02 15:04:05")
}

//获取时间戳秒
func GetNowTimestamp() int64{
	var cstZone = time.FixedZone("CST", 8*3600)
	return time.Now().In(cstZone).UnixNano()/1000000000
}

func GetNowByFormat(format string) string{
	var cstZone = time.FixedZone("CST", 8*3600)       // 东八
	return time.Now().In(cstZone).Format(format)
}

func GetLastMonth() string{
	var cstZone = time.FixedZone("CST", 8*3600)       // 东八
	return time.Now().AddDate(0, -1, 0).In(cstZone).Format("2006-01")
}
package consul

import (
	"main/src/common"
	"time"
)

func AutoRegister() {
	// 定时把consul里的rds信息重新注册
	for {
		common.LogWriterLn(common.LogPath, common.INFO, "=======================定时任务AutoRegister开始==================")
		rdsList := GetDBInfoFromConsul()
		statusCode, err := Register(rdsList)
		if statusCode == 200 {
			common.LogWriterLn(common.LogPath, common.INFO, "定时任务，每日重新自动注册consul内的rds实例信息, Success!")
		} else {
			common.LogWriterLn(common.LogPath, common.ERROR, "定时任务，每日重新自动注册consul内的rds实例信息, Failed!"+err.Error())
		}
		common.LogWriterLn(common.LogPath, common.INFO, "=======================定时任务AutoRegister结束==================")
		time.Sleep(24 * time.Hour)
	}
}

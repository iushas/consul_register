package common

import "encoding/json"

func GetRegTemplate() RegTemplateJson {
	// 获取注册模板信息

	regTemplate := FileReadAll(RegTemplatePath)

	config := RegTemplateJson{}

	err := json.Unmarshal([]byte(regTemplate), &config)

	if err != nil {
		LogWriterLn(LogPath, ERROR, "Get regTemplate err: "+err.Error())
	}
	return config
}

func GetServerConf() ServerConfJson {
	//获取服务启动参数

	serverConf := FileReadAll(ServerRegTemplatePath)

	config := ServerConfJson{}

	err := json.Unmarshal([]byte(serverConf), &config)

	if err != nil {
		LogWriterLn(LogPath, ERROR, "Get serverConf err: "+err.Error())
	}
	return config
}

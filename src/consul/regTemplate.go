package consul

import (
	"encoding/json"
	"io"
	"io/ioutil"
	_ "io/ioutil"
	"main/src/common"
	"os"
	"strconv"
)

//此文件函数为consul注册模板配置和获取，可通过http接口操作

func SetConsulTemplateJson(configTemplate string) (StatusCode int, err error) {
	config := common.RegTemplateJson{}
	err = json.Unmarshal([]byte(configTemplate), &config)
	if err != nil {
		common.LogWriterLn(common.LogPath, common.INFO, "Parse configTemplate error: "+err.Error())
		return 500, err
	}

	/*判断模板正确性*/
	/*tag*/
	if len(config.ConsulTemplate.Tags) == 0 || config.ConsulTemplate.Tags == nil {
		common.LogWriterLn(common.LogPath, common.ERROR, "ConfigTemplate Tags IS Error, Please input one tag at least!")
		return 500, nil
	}
	/*for _, tag := range config.ConsulTemplate.Tags {
		if tag != "region" && tag != "business" && tag != "module" {
			common.LogWriterLn(common.LogPath, common.ERROR, "ConfigTemplate Tags IS Error, tag must in region | business| module, error tag is:"+tag)
			return 500, nil
		}
	}*/
	/*check*/
	if config.ConsulTemplate.Check.DeregisterCriticalServiceAfter == "" || config.ConsulTemplate.Check.Interval == "" || config.ConsulTemplate.Check.Timeout == "" {
		common.LogWriterLn(common.LogPath, common.ERROR, "ConsulTemplate.Check has empty attribute!")
		return 500, nil
	} else {
		DeregisterCriticalServiceAfter, err := strconv.Atoi(config.ConsulTemplate.Check.DeregisterCriticalServiceAfter[0 : len(config.ConsulTemplate.Check.DeregisterCriticalServiceAfter)-1])
		if err != nil || DeregisterCriticalServiceAfter == 0 {
			common.LogWriterLn(common.LogPath, common.ERROR, "ConsulTemplate.Check.DeregisterCriticalServiceAfter is wrong!, Input is:"+config.ConsulTemplate.Check.DeregisterCriticalServiceAfter)
			return 500, nil
		}
		Interval, err := strconv.Atoi(config.ConsulTemplate.Check.Interval[0 : len(config.ConsulTemplate.Check.Interval)-1])
		if err != nil || Interval == 0 {
			common.LogWriterLn(common.LogPath, common.ERROR, "ConsulTemplate.Check.Interval is wrong!, Input is:"+config.ConsulTemplate.Check.Interval)
			return 500, nil
		}
		Timeout, err := strconv.Atoi(config.ConsulTemplate.Check.Timeout[0 : len(config.ConsulTemplate.Check.Timeout)-1])
		if err != nil || Timeout == 0 {
			common.LogWriterLn(common.LogPath, common.ERROR, "ConsulTemplate.Check.Timeout is wrong!, Input is:"+config.ConsulTemplate.Check.Timeout)
			return 500, nil
		}
	}

	f, err := os.OpenFile(common.RegTemplatePath, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0644)
	if err != nil {
		common.LogWriterLn(common.LogPath, common.INFO, "file create failed. err: "+err.Error())
		return 500, err
	}
	n, _ := f.Seek(0, io.SeekEnd)
	_, err = f.WriteAt([]byte(configTemplate), n)
	defer f.Close()
	return 200, nil
}

func GetConsulTemplateJson() (StatusCode int, err error, configTemplate string) {
	config, err := ioutil.ReadFile(common.RegTemplatePath)
	if err != nil {
		common.LogWriterLn(common.LogPath, common.ERROR, "Get ConsulJson err: "+err.Error())
		return 500, err, ""
		panic(err)
	}
	configTemplate = string(config)
	return 200, nil, configTemplate
}

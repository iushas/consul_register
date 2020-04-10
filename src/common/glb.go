package common

import (
	"time"
)

const (
	INFO  string = "[INFO]"
	ERROR string = "[ERROR]"
	DEBUG string = "[DEBUG]"
)

//服务路径配置, 后面代码修改到服务自身配置
//TODO
var BasePath = ""
var RegTemplate = RegTemplateJson{}
var RegTemplatePath = "conf/reg_tpl.json" //对应consul.json
var LogPath = "log/" + time.Now().Format("2006-01-02") + ".log"

var ServerConfig = ServerConfJson{} //对应conf.json
var ServerRegTemplatePath = "conf/conf.json"

// 服务注册rds信息配置
type ConsulStruct struct {
	ID      string   `json:"ID"`
	Name    string   `json:"Name"`
	Address string   `json:"Address"`
	Port    int      `json:"Port"`
	Tags    []string `json:"Tags"`
	Meta    Meta     `json:"Meta"`
	Check   Check    `json:"Check"`
}
type Meta struct {
	Alias         string `json:"alias"`
	IP            string `json:"ip"`
	Region        string `json:"region"`
	ServiceDetail string `json:"service_detail"`
	ServiceType   string `json:"service_type"`
}
type Check struct {
	DeregisterCriticalServiceAfter string `json:"DeregisterCriticalServiceAfter"`
	Interval                       string `json:"Interval"`
	Name                           string `json:"Name"`
	TCP                            string `json:"TCP"`
	Timeout                        string `json:"Timeout"`
}

//注册模板可变变量的配置
type RegTemplateJson struct {
	ConsulAddress  string `json:"ConsulAddress"`
	ConsulTemplate struct {
		Tags  []string `json:"Tags"`
		Check struct {
			DeregisterCriticalServiceAfter string `json:"DeregisterCriticalServiceAfter"`
			Interval                       string `json:"Interval"`
			Timeout                        string `json:"Timeout"`
		} `json:"Check"`
	} `json:"ConsulTemplate"`
}

//服务自身配置

type ServerConfJson struct {
	Consul struct {
		HTTPReqIsAuth string `json:"HttpReqIsAuth"`
		UserName      string `json:"UserName"`
		Password      string `json:"Password"`
	} `json:"Consul"`
	Server struct {
		Address string `json:"Address"`
		Port    int    `json:"Port"`
	} `json:"Server"`
	Aliyun struct {
		RegionId        string
		AccessKeyId     string
		AccessKeySecret string
	} `json:"Aliyun"`
}

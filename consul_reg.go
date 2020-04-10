package main

import (
	"flag"
	"main/src/common"
	"main/src/consul"
	_ "main/src/consul"
	"main/src/handler"
	"os"
	"strconv"
	"sync"
	"time"
)

func init() {
	//读取相关配置
	common.BasePath, _ = os.Getwd()
	common.RegTemplatePath = common.BasePath + "/" + common.RegTemplatePath
	common.LogPath = common.BasePath + "/" + common.LogPath
	common.RegTemplate = common.GetRegTemplate()
	common.ServerConfig = common.GetServerConf()
}

func main() {
	var wg sync.WaitGroup
	wg.Add(1)
	//http server go-routine
	go func() {
		defer wg.Done()
		HttpServer()
	}()

	//sleep 一下，等服务启动在启动定时钟
	time.Sleep(60 * time.Second)

	//定时钟 go-routine
	//暂时取消自动注册
	/*go func() {
		defer wg.Done()
		ticking()
	}()*/
	wg.Wait()

	common.LogWriterLn(common.LogPath, common.INFO, "Server Terminated ")
}

func HttpServer() {
	var address string
	flag.StringVar(&address, "address", "", "server address")
	var port string
	flag.StringVar(&port, "port", "", "server port")
	flag.Usage()
	flag.Parse()
	ListenAddress := ""
	if address != "" && port != "" {
		ListenAddress = address + ":" + port
	} else if address != "" {
		ListenAddress = address + ":" + strconv.Itoa(common.ServerConfig.Server.Port)
	} else if port != "" {
		ListenAddress = common.ServerConfig.Server.Address + ":" + port
	} else if address == "" && port == "" {
		ListenAddress = common.ServerConfig.Server.Address + ":" + strconv.Itoa(common.ServerConfig.Server.Port)
	}
	common.LogWriterLn(common.LogPath, common.INFO, "Server BasePath: "+common.BasePath)
	common.LogWriterLn(common.LogPath, common.INFO, "Server LogPath: "+common.LogPath)
	common.LogWriterLn(common.LogPath, common.INFO, "Consul RegTemplatePath: "+common.RegTemplatePath)
	common.LogWriterLn(common.LogPath, common.INFO, "Server Started: "+ListenAddress)
	handler.HandleRequest(ListenAddress)
}

func ticking() {
	consul.AutoRegister()
}

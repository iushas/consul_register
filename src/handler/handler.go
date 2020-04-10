package handler

import (
	"encoding/json"
	"io/ioutil"
	"main/src/common"
	"main/src/consul"
	"net/http"
)

func HandleRequest(serverAddress string) {
	server := http.NewServeMux()
	server.HandleFunc("/register", handleRegister)
	server.HandleFunc("/deregister", handleDeRegister)
	server.HandleFunc("/setConsulTemplate", handleSetConsulConfig)
	server.HandleFunc("/getConsulTemplate", handleGetConsulConfig)
	_ = http.ListenAndServe(serverAddress, server)
}

func handleRegister(w http.ResponseWriter, r *http.Request) {
	common.LogWriterLn(common.LogPath, common.INFO, r.URL)
	body, err := ioutil.ReadAll(r.Body)

	if err != nil {
		w.WriteHeader(500)
		w.Header()
		_, _ = w.Write([]byte("Parse Request Expectation Failed"))
		common.LogWriterLn(common.LogPath, common.ERROR, "Parse Request Expectation Failed!")
		return
	}

	type RdsList []string
	list := RdsList{}
	_ = json.Unmarshal(body, &list)
	common.LogWriterLn(common.LogPath, common.INFO, "Register request rds list： "+string(body))

	if len(list) == 0 {
		w.WriteHeader(400)
		w.Header()
		_, _ = w.Write([]byte("Parameters is NULL OR format is wrong!"))
		common.LogWriterLn(common.LogPath, common.ERROR, "Parameters is NULL OR format is wrong!")
		return
	}
	statusCode, _ := consul.Register(list)
	if statusCode == 200 {
		w.WriteHeader(statusCode)
		w.Header()
		_, _ = w.Write([]byte("Register request  Success"))
		common.LogWriterLn(common.LogPath, common.INFO, "Register request  Success!")
	} else {
		w.WriteHeader(statusCode)
		w.Header()
		_, _ = w.Write([]byte("Register request Failed"))
		common.LogWriterLn(common.LogPath, common.ERROR, "Register request Failed!")
		return
	}
}

func handleDeRegister(w http.ResponseWriter, r *http.Request) {
	common.LogWriterLn(common.LogPath, common.INFO, r.URL)
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(500)
		w.Header()
		_, _ = w.Write([]byte("Parse Request Expectation Failed"))
		common.LogWriterLn(common.LogPath, common.ERROR, "Parse Request Expectation Failed!")
		return
	}

	type RdsList []string
	list := RdsList{}
	_ = json.Unmarshal(body, &list)

	common.LogWriterLn(common.LogPath, common.INFO, "Deregister request rds list： "+string(body))
	if len(list) == 0 {
		w.WriteHeader(400)
		w.Header()
		_, _ = w.Write([]byte("Parameters is NULL OR format is wrong!"))
		common.LogWriterLn(common.LogPath, common.ERROR, "Parameters is NULL OR format is wrong!")
		return
	}

	statusCode, _ := consul.DeRegister(list)

	if statusCode == 200 {
		w.WriteHeader(statusCode)
		w.Header()
		_, _ = w.Write([]byte("Deregister request success"))
		common.LogWriterLn(common.LogPath, common.INFO, "Deregister request success！")
	} else {
		w.WriteHeader(statusCode)
		w.Header()
		_, _ = w.Write([]byte("Deregister request Failed"))
		common.LogWriterLn(common.LogPath, common.ERROR, "Deregister request Failed！")
	}

}

func handleSetConsulConfig(w http.ResponseWriter, r *http.Request) {
	common.LogWriterLn(common.LogPath, common.INFO, r.URL)
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(500)
		w.Header()
		_, _ = w.Write([]byte(err.Error()))
		common.LogWriterLn(common.LogPath, common.ERROR, err.Error())
	}

	if len(body) == 0 {
		w.WriteHeader(400)
		w.Header()
		_, _ = w.Write([]byte("Expectation Failed"))
		common.LogWriterLn(common.LogPath, common.ERROR, err.Error())
		return
	}
	configString := string(body)

	common.LogWriterLn(common.LogPath, common.INFO, "New consul config template："+configString)

	statusCode, err := consul.SetConsulTemplateJson(configString)

	if statusCode == 200 {
		w.WriteHeader(statusCode)
		w.Header()
		_, _ = w.Write([]byte("New consul config template write success"))
		common.LogWriterLn(common.LogPath, common.INFO, "New consul config template write success！")
	} else {
		w.WriteHeader(statusCode)
		w.Header()
		_, _ = w.Write([]byte("New consul config template write error"))
		common.LogWriterLn(common.LogPath, common.ERROR, "New consul config template write error!")
	}
}

func handleGetConsulConfig(w http.ResponseWriter, r *http.Request) {
	common.LogWriterLn(common.LogPath, common.INFO, r.URL)

	statusCode, err, configString := consul.GetConsulTemplateJson()

	if statusCode == 200 && err == nil {
		w.WriteHeader(statusCode)
		w.Header()
		_, _ = w.Write([]byte(configString))
		common.LogWriterLn(common.LogPath, common.INFO, "Get consul config template  success！")
	} else {
		w.WriteHeader(statusCode)
		w.Header()
		_, _ = w.Write([]byte("Get consul config template  error"))
		common.LogWriterLn(common.LogPath, common.ERROR, "Get consul config template  error!")
	}
}

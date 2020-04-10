package consul

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"main/src/common"
	"net/http"
)

func GetDBInfoFromConsul() (RdsList []string) {
	client := http.Client{}
	//只拉服务db规则的数据
	reqDetail := "/v1/agent/services?filter=Service%20matches%20db"

	req, err := http.NewRequest(http.MethodGet, common.RegTemplate.ConsulAddress+reqDetail, nil)

	if common.ServerConfig.Consul.HTTPReqIsAuth == "True" {
		//add http basic auth info
		user := common.ServerConfig.Consul.UserName
		password, _ := common.AesDecrypt(common.ServerConfig.Consul.Password)
		req.SetBasicAuth(user, password)
	}

	if err != nil {
		common.LogWriterLn(common.LogPath, common.ERROR, err.Error())

	}
	//req.Header.Add("Content-Type", "application/json; charset=utf-8")

	common.LogWriterLn(common.LogPath, common.INFO, fmt.Sprintf("Get rdsinfo from consul, url: %s  body:%s", req.URL, req.Body))

	resp, err := client.Do(req)
	if err != nil {
		common.LogWriterLn(common.LogPath, common.ERROR, err.Error())
	}

	if resp.StatusCode == 200 {
		common.LogWriterLn(common.LogPath, common.INFO, fmt.Sprintf("Get rdsinfo from consul, code: %d,  body:%s", resp.StatusCode, resp.Status))
		strBody, _ := ioutil.ReadAll(resp.Body)

		//解析字符传，获取rdsId信息
		var data interface{}
		json.Unmarshal([]byte(strBody), &data)
		dataV := data.(map[string]interface{})
		for i := range dataV {
			RdsList = append(RdsList, i)
		}
		common.LogWriterLn(common.LogPath, common.INFO, "Turn resp body to rdsList success!")
	} else {
		common.LogWriterLn(common.LogPath, common.ERROR, fmt.Sprintf("Turn resp body to rdsList Error, code: %d,  body:%s", resp.StatusCode, resp.Status))
	}

	fmt.Println(RdsList)
	return RdsList
}

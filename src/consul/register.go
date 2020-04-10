package consul

import (
	"bytes"
	"encoding/json"
	"fmt"
	"main/src/common"
	"main/src/rds"
	"net/http"
	"strconv"
)

func GetFromFile(filePath string) (dbConsulArr []common.ConsulStruct) {

	config := common.GetRegTemplate()
	rdslist := common.FileReadAll(filePath)

	var rdsRegInfo map[string][]string
	json.Unmarshal([]byte(rdslist), &rdsRegInfo)

	for k, v := range rdsRegInfo {
		_, _, rdsinfo := rds.GetRdsInfo([]string{k})
		item := rdsinfo[0]
		//处理 meta和check部分
		tMetaMap := common.Meta{}
		tMetaMap.Alias = item.DBInstanceDescription
		tMetaMap.IP = item.ConnectionString
		tMetaMap.Region = item.RegionId
		tMetaMap.ServiceDetail = item.Engine + item.EngineVersion + "|" + item.DBInstanceCPU + "C|" + strconv.FormatInt(item.DBInstanceMemory/1024, 10) + "G|" + strconv.Itoa(item.DBInstanceStorage) + "G"
		tMetaMap.ServiceType = item.DBInstanceType

		tCheckMap := common.Check{}
		tCheckMap.Name = item.DBInstanceId
		tCheckMap.DeregisterCriticalServiceAfter = config.ConsulTemplate.Check.DeregisterCriticalServiceAfter
		tCheckMap.Interval = config.ConsulTemplate.Check.Interval
		tCheckMap.TCP = item.ConnectionString + ":" + item.Port
		tCheckMap.Timeout = config.ConsulTemplate.Check.Timeout
		port, _ := strconv.Atoi(item.Port)

		//处理通用部分
		dbConsul := common.ConsulStruct{
			ID:      item.DBInstanceId,
			Name:    "db-mysql",
			Address: item.ConnectionString,
			Port:    port,
			Tags:    v,
			Meta:    tMetaMap,
			Check:   tCheckMap,
		}
		//处理Tags
		/*for _, value := range config.ConsulTemplate.Tags {
			switch value {
			case "region":
				if item.RegionId == "cn-hangzhou" {
					dbConsul.Tags = append(dbConsul.Tags, "hz-ali")
				}
			case "business":
				//TODO
			case "module":
				//TODO
			}
		}*/

		dbConsulArr = append(dbConsulArr, dbConsul)
	}
	return dbConsulArr
}

func Get(rdsID []string) (int, error, []common.ConsulStruct) {
	//再读取一次注册模板配置
	common.RegTemplate = common.GetRegTemplate()

	dbConsulArr := []common.ConsulStruct{}
	status, err, rdsInfoList := rds.GetRdsInfo(rdsID)

	if status != 200 {
		return status, err, nil
	}

	tags := common.RegTemplate.ConsulTemplate.Tags

	for _, item := range rdsInfoList {

		//处理 meta和check部分
		tMetaMap := common.Meta{}
		tMetaMap.Alias = item.DBInstanceDescription
		tMetaMap.IP = item.ConnectionString
		tMetaMap.Region = item.RegionId
		tMetaMap.ServiceDetail = item.Engine + item.EngineVersion + "|" + item.DBInstanceCPU + "C|" + strconv.FormatInt(item.DBInstanceMemory/1024, 10) + "G|" + strconv.Itoa(item.DBInstanceStorage) + "G"
		tMetaMap.ServiceType = item.DBInstanceType

		tCheckMap := common.Check{}
		tCheckMap.Name = item.DBInstanceId
		tCheckMap.DeregisterCriticalServiceAfter = common.RegTemplate.ConsulTemplate.Check.DeregisterCriticalServiceAfter
		tCheckMap.Interval = common.RegTemplate.ConsulTemplate.Check.Interval
		tCheckMap.TCP = item.ConnectionString + ":" + item.Port
		tCheckMap.Timeout = common.RegTemplate.ConsulTemplate.Check.Timeout
		port, _ := strconv.Atoi(item.Port)

		//处理通用部分
		dbConsul := common.ConsulStruct{
			ID:      item.DBInstanceId,
			Name:    "db-mysql",
			Address: item.ConnectionString,
			Port:    port,
			Tags:    tags,
			Meta:    tMetaMap,
			Check:   tCheckMap,
		}

		//处理Tags
		/*for _, value := range common.RegTemplate.ConsulTemplate.Tags {
			switch value {
			case "region":
				if item.RegionId == "cn-hangzhou" {
					dbConsul.Tags = append(dbConsul.Tags, "hz-ali")
				}
			case "business":
				//TODO
			case "module":
			//TODO

			default:
				common.LogWriterLn(common.LogPath, common.ERROR, "不正确的Tag输入："+value)
				common.LogWriterLn(common.LogPath, common.ERROR, "Tag只能包含region|business|module 一个或多个")

			}
		}*/

		dbConsulArr = append(dbConsulArr, dbConsul)
	}
	return 200, nil, dbConsulArr
}

func Put(dbConsulArr []common.ConsulStruct) (StatusCode int, err error) {

	for _, item := range dbConsulArr {
		jsonReq, err := json.Marshal(item)

		if err != nil {
			common.LogWriterLn(common.LogPath, common.ERROR, err.Error())
			return 500, err
		}

		client := http.Client{}
		reqDetail := "/v1/agent/service/register"
		req, err := http.NewRequest(http.MethodPut, common.RegTemplate.ConsulAddress+reqDetail, bytes.NewBuffer(jsonReq))
		if common.ServerConfig.Consul.HTTPReqIsAuth == "True" {
			//add http basic auth info
			user := common.ServerConfig.Consul.UserName
			password, _ := common.AesDecrypt(common.ServerConfig.Consul.Password)
			req.SetBasicAuth(user, password)
		}
		if err != nil {
			common.LogWriterLn(common.LogPath, common.ERROR, err.Error())
			return 501, err
		}
		req.Header.Add("Content-Type", "application/json; charset=utf-8")

		common.LogWriterLn(common.LogPath, common.INFO, fmt.Sprintf("Register request url: %s  body:%s", req.URL, req.Body))

		resp, err := client.Do(req)
		if err != nil {
			common.LogWriterLn(common.LogPath, common.ERROR, err.Error())
			return resp.StatusCode, err
		}
		defer resp.Body.Close()

		if resp.StatusCode == 200 {
			common.LogWriterLn(common.LogPath, common.INFO, fmt.Sprintf("Register response code: %d,  body:%s", resp.StatusCode, resp.Status))
		} else {
			common.LogWriterLn(common.LogPath, common.ERROR, fmt.Sprintf("Register response code: %d,  body:%s", resp.StatusCode, resp.Status))
			return resp.StatusCode, err
		}

	}
	return 200, nil

}

func Delete(rdsList []string) (StatusCode int, err error) {

	for _, item := range rdsList {

		client := http.Client{}
		reqDetail := "/v1/agent/service/deregister/" + item
		req, err := http.NewRequest(http.MethodPut, common.RegTemplate.ConsulAddress+reqDetail, nil)
		if err != nil {
			return 500, err
		}

		if common.ServerConfig.Consul.HTTPReqIsAuth == "True" {
			//add http basic auth info
			user := common.ServerConfig.Consul.UserName
			password, _ := common.AesDecrypt(common.ServerConfig.Consul.Password)
			req.SetBasicAuth(user, password)
		}

		common.LogWriterLn(common.LogPath, common.INFO, fmt.Sprintf("DeRegister request url: %s  body:%s", req.URL, req.Body))

		resp, err := client.Do(req)
		if err != nil {
			common.LogWriterLn(common.LogPath, common.ERROR, err.Error())
			return resp.StatusCode, err
		}
		defer resp.Body.Close()

		if resp.StatusCode == 200 {
			common.LogWriterLn(common.LogPath, common.INFO, fmt.Sprintf("DeRegister response code: %d,  body:%s", resp.StatusCode, resp.Status))
		} else {
			common.LogWriterLn(common.LogPath, common.ERROR, fmt.Sprintf("DeRegister response code: %d,  body:%s", resp.StatusCode, resp.Status))
			return resp.StatusCode, nil
		}
	}
	return 200, nil
}

func Register(rdsList []string) (StatusCode int, err error) {
	status, err, regRequest := Get(rdsList)
	if status != 200 {
		return status, err
	}
	StatusCode, err = Put(regRequest)
	return StatusCode, err
}

func DeRegister(rdsList []string) (StatusCode int, err error) {
	status, err, _ := Get(rdsList)
	if status != 200 {
		return status, err
	}
	StatusCode, err = Delete(rdsList)
	return StatusCode, err
}

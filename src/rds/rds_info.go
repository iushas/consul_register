package rds

import (
	"github.com/aliyun/alibaba-cloud-sdk-go/services/rds"
	"main/src/common"
)

func GetRdsInfo(rdsList []string) (statusCode int, err error, response []rds.DBInstanceAttributeInDescribeDBInstanceAttribute) {
	regionId := common.ServerConfig.Aliyun.RegionId
	accessKey, _ := common.AesDecrypt(common.ServerConfig.Aliyun.AccessKeyId)
	accessSecret, _ := common.AesDecrypt(common.ServerConfig.Aliyun.AccessKeySecret)
	client, err := rds.NewClientWithAccessKey(
		regionId,
		accessKey,
		accessSecret)

	if err != nil {
		common.LogWriterLn(common.LogPath, common.ERROR, err)
		return 590, err, nil
		panic(err)
	}
	request := rds.CreateDescribeDBInstanceAttributeRequest()
	for _, item := range rdsList {
		request.DBInstanceId = item

		resp, err := client.DescribeDBInstanceAttribute(request)

		if err != nil {
			common.LogWriterLn(common.LogPath, common.ERROR, err)
			return resp.GetHttpStatus(), err, nil
			panic(err)
		}
		if resp.GetHttpStatus() == 200 {
			// results  merge
			response = MergeResponse(response, resp.Items.DBInstanceAttribute)
		}
	}
	return 200, err, response
}

func MergeResponse(response1 []rds.DBInstanceAttributeInDescribeDBInstanceAttribute, response2 []rds.DBInstanceAttributeInDescribeDBInstanceAttribute) (response []rds.DBInstanceAttributeInDescribeDBInstanceAttribute) {
	response = make([]rds.DBInstanceAttributeInDescribeDBInstanceAttribute, len(response1)+len(response2))
	if len(response1) != 0 && len(response2) != 0 {
		copy(response, response1)
		copy(response[len(response1):], response2)

	} else {
		if len(response1) != 0 {
			copy(response, response1)
		} else {
			copy(response, response2)
		}
	}
	return response
}

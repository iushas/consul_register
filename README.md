# consul_register
目前用来实现数据库信息注册到consul， 服务里定义好了注册模板，并且支持特定注册参数的修改

1、注册模板

模板举例：

s_id=''    // 实例ID

s_name='db-mysql'

s_ip='rm-xxxxxxxxxxxxxx.mysql.rds.aliyuncs.com'    //实例内网域名

s_port=3306

s_tcp='rm-xxxxxxxxxxxxxx.mysql.rds.aliyuncs.com:3306'   //域名+端口

s_region='rds'                //根据实例所在区来注册， 比如sh_idc, hz_idc ...

s_put="""

{

\"ID\": \"$s_id\",

\"Name\": \"$s_name\",

\"Address\": \"$s_ip\",

\"Port\": $s_port,

\"Tags\": [\"$s_region\"],

\"Meta\": {

\"region\": \"$s_region\",

\"alias\": \"$s_id\",

\"ip\": \"$s_ip\",

\"service_type\": \"mysql-master\",        //根据rds 是主库还是从库，从库为mysql-slave

\"service_detail\": \"杭州作为测试使用\"

} ,

\"Check\": {

\"Name\": \"$s_id\",

\"DeregisterCriticalServiceAfter\": \"1440m\",

\"TCP\": \"$s_tcp\",

\"Interval\": \"60s\",

\"Timeout\": \"5s\"

}

}

"""

echo $s_put

curl http://consul地址/v1/agent/service/register -H 'Accept: application/json' -X PUT  -d "$s_put"

2、注销模板

curl --request PUT http://consul地址/v1/agent/service/deregister/id

id 为要注销的rds 实例id

3、HTTP请求

目前ops运维平台的功能集成方式，为负责人开发相应功能， 运维开发集成http接口。

故consul的注册和注销功能也是服务化的， 服务采用go语言开发，端口暂时采用9160端口。

这里主要涉及到4个接口。

1）consul 地址和模板

curl -i http://127.0.0.1:9160/setConsulTemplate-H 'Accept: application/json' -X PUT -d "{\"ConsulAddress\":\"http://consul地址\",\"ConsulTemplate\": {\"Tags\":[\"region\",\"businessdfg\",\"module\"],\"Check\":{ \"DeregisterCriticalServiceAfter\":\"1400m\",\"Interval\":\"60s\",\"Timeout\": \"5s\" } }}"

Tags值为具体的tag列表

2) 获取consul 模板

curl -i http://127.0.0.1:9160/getConsulTemplate-XGET

3）注册请求

curl -i http://127.0.0.1:9160/register -H  'Accept: application/json'  -X PUT  -d"[\"rm-xxxxxxx1\",\"rm-xxxxxxx2\"]"

4）注销请求

curl -i http://127.0.0.1:9160/deregister -H 'Accept: application/json'  -X PUT  -d "[\"rm-xxxxxxx1\",\"rm-xxxxxxx2\"]"
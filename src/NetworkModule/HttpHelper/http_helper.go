package HttpHelper

// help页面
func Help() string {
	info := `                    帮助页面 
一：基本接口
1:登陆主机服务
  请求url:/host/login/service

2:登出主机服务
  请求url:/host/logout/service

3:修改密码
  请求url:/host/change/password

4:获取服务器的网络配置(包括IP,子网掩码,网关)
  请求url: /host/get/network/configure

5:更新服务器的网络配置(包括IP,子网掩码,网关)
  请求url: /host/update/network/configure

6:控制服务组件(启动/停止/删除/启用/禁用)
  请求url: /host/operate/service

7:查看服务组件信息
  请求url: /host/get/service/info

8:下载服务组件的日志
  请求url: /host/download/service/log

9:升级服务组件
  请求url: /host/update/service

10:获取系统时间
  请求url: /host/get/time

11:设置系统时间
  请求url: /host/update/time

12:获取ntp配置
  请求url: /host/get/ntp

13:更新ntp配置
  请求url: /host/update/ntp

14:操作服务器(关机/重启)
  请求url: /host/operate/server

15:获取服务器标识码
  请求url: /host/get/server/id

16:更新license
  请求url: /host/update/license

17:获取license信息
  请求url: /host/get/license

18:获取所有的时区信息
  请求url: /host/get/time/zones/info

19:获取当前服务器的时区信息
  请求url: /host/get/server/time/zone

20:更新时区
  请求url: /host/update/server/time/zone

21:查看服务器的系统日志
  请求url: /host/get/server/log

22:查看服务组件的日志
  请求url: /host/get/service/log

23: 下载服务器的系统日志
  请求url: /host/download/server/log

24: 统计CPU的使用情况
  请求url: /host/get/cpu/statistic

25: 统计磁盘的使用情况
  请求url: /host/get/disk/statistic

26: 统计网络的使用情况
  请求url: /host/get/network/statistic

27: 统计内存的使用情况
  请求url: /host/get/memory/statistic

28: 统计服务组件的资源使用情况
  请求url: /host/get/service/statistic

29: 查询磁盘的存储配置(获取存储的预留空间)
  请求url: /host/get/storage/configure

30: 设置磁盘的存储配置(设置存储的预留空间)
  请求url: /host/storage/configure

31: 查看磁盘状态的信息
  请求url: /host/get/disk/status/info

31: DES加密
  请求url: /host/encrypt/des

31: DES解密
  请求url: /host/decrypt/des

二：页面接口
1:获取时间信息
  请求url:/host/get/time/info/page

2:设置时间信息
  请求url:/host/update/time/info/page
`
	return info
}

// 获取网络配置
func NetworkConfigureHelp() (info string) {
	info = `
说明: 获取服务器的网络配置信息
请求url: /host/get/network/configure
请求body:无

返回示例:
{
 "network": [				// 此处为数组，如果有多块网卡，则返回多个对象
  {
   "MAC": "ac:1f:6b:b1:40:92",		// MAC地址
   "name": "enp3s0f0",				// 网口名称
   "host_name": "",					// 主机名称
   "IPv4": {						// IPv4信息
    "ip": "192.168.1.230",
    "netmask": "255.255.255.0",
    "gateway": "192.168.1.1",
    "auto_configure": false,
    "enable": true
   },
   "IPv6": {						// IPv6信息
    "ip": "",
    "netmask": "",
    "gateway": "",
    "auto_configure": true,
    "enable": true
   },
   "DNS": {							// DNS信息
    "major_dns": "114.114.114.114",
    "minor_dns": "8.8.8.8"
   },
   "enable": true,				// 网卡是否插了网线，如果有插网线，则为true
   "bandwidth": 105485241,		// 网卡的带宽，单位为byte(注意，我们平时说的百兆网卡，指的是bit)
   "bandwidth_readable": "125.0M" // 可读性好的网卡的带宽
  }
 ]
}`
	return info
}

// 更新网络配置
func UpdateNetworkConfigureHelp() (info string) {
	info = `
说明: 更新服务器的网络配置信息
请求url: /host/update/network/configure
请求body:
{
   "MAC": "ac:1f:6b:b1:40:93",		// MAC地址
   "name": "enp3s0f1",				// 网口名称
   "host_name": "",					// 主机名称
   "IPv4": {		// IPv4信息
    "ip": "192.168.5.5",
    "netmask": "255.255.255.0",
    "gateway": "192.168.5.10",
    "auto_configure": false,
    "enable": true
   },
   "IPv6": {		// IPv6信息
    "ip": "",
    "netmask": "",
    "gateway": "",
    "auto_configure": true,
    "enable": true
   },
   "DNS": {		// DNS信息
    "major_dns": "114.114.114.114",
    "minor_dns": "8.8.8.8"
   }
}

返回示例:
成功：
{
 "status": 200,
 "reason": "OK"
}

失败
{
 "status": 400,			// 或者其他非200的数组
 "reason": "some error"	// 错误描述
}
 `
	return info
}

// 登陆的help
func LoginServiceHelp() (info string) {
	info = `
说明: 登陆主机服务
请求url: /host/login/service
请求body:
{
   "user": "admin",		// 用户名
   "password": "21232f297a57a5a743894a0e4a801fc3",				// 用户密码的md5
} 

返回示例:
成功：
{
 "status": 200,
 "reason": "OK"
 "token": "fesdfale9sdf"			// 以后所有调用主机服务的http接口，head中都要包含Token字段
}

失败
{
 "status": 403,			// 或者其他非200的数组
 "reason": "invalid user or password"	// 错误描述
}
 `
	return info
}

func LogoutServiceHelp() (info string) {
	info = `
说明: 退出主机服务
请求url: /host/logout/service
请求body:
{
   "user": "admin",		// 用户名
   "password": "21232f297a57a5a743894a0e4a801fc3",				// 用户密码的md5
} 

返回示例:
成功：
{
 "status": 200,
 "reason": "OK"
}

失败
{
 "status": 403,			// 或者其他非200的数组
 "reason": "invalid user or password"	// 错误描述
}
 `
	return info
}

// 修改密码
func ChangePasswordHelp() (info string) {
	info = `
说明: 修改密码
请求url: /host/change/password
请求body:
{
   "user": "admin",		// 用户名
   "old_password": "21232f297a57a5a743894a0e4a801fc3",				// 用户旧密码的md5
   "new_password": "e10adc3949ba59abbe56e057f20f883e",				// 用户新密码的md5
} 

返回示例:
成功：
{
 "status": 200,
 "reason": "OK"
}

失败
{
 "status": 403,			// 或者其他非200的数组
 "reason": "invalid user or password"	// 错误描述
}
 `
	return info
}

// 服务控制
func ControlServiceHelper() (info string) {
	info = `
说明: 操作服务组件
请求url: /host/operate/service
请求body:
{
    "name":"service",		// 服务组件的名称(参考 /host/get/service/info 的返回结果)
    "control":"stop", 		// 操作类型(start/stop/restart/remove/disable/enable)
    "time":"now"			// 操作时间，now为立即起效，如果为定时操作，则传入起效时间，格式为:2010-10-10 02:03
} 

返回示例:
成功:
{
 "status": 200,
 "reason": "OK"
}

失败:
{
 "status": 400,			// 或者其他非200的数组
 "reason": "err reason"	// 错误描述
}
 `
	return info
}

// 服务模块信息
func ServiceInfoHelper() (info string) {
	info = `
说明: 获取服务组件的信息
请求url: /host/get/service/info
请求body: 无

返回示例:
成功:
{
 "service": [
  {
   "path": "/mars/elasticsearch",	// 服务的安装路径
   "name": "elasticsearch",			// 服务组件的名称
   "status": "stop",				// 服务组件的状态
   "has_status": true,			// 服务组件是否有运行状态(web, flyway等没有运行状态)
   "has_log": true,				// 服务组件是否有日志(web, flyway等没有运行状态)
   "has_update": true,			// 服务组件是否可以升级(mysql是不可以升级的)
   "git": {				// 服务组件的git信息
    "branch": "master",	
    "commit": "098cbad5",
    "version": "v0.60001"
   }
  }
 ]
}

失败:
{
 "status": 400,			// 或者其他非200的数组
 "reason": "err reason"	// 错误描述
}
 `
	return info
}

// 下载服务组件的日志
func DownloadServiceLogHelp() (info string) {
	info = `
说明: 下载服务组件的日志(下载整个文件夹)
请求url: /download/service/log
请求body: 
{
  "name":"collect"		// 服务组件名称(参考 /host/get/service/info 的返回结果)
}

返回示例:
成功:
{
 "status": 200,
 "reason": "OK"	
 "path": "/html/download/log/collect.zip"  // 日志的下载路径
}

失败:
{
 "status": 400,			// 或者其他非200的数组
 "reason": "err reason"	// 错误描述
}
 `
	return info
}

func UpdateServiceHelper() (info string) {
	info = `
说明: 安装/升级服务组件
请求url: /host/update/service
请求body(body由两部分组成，第一部分为json格式的升级/安装信息,第二部分为升级包的内容): 
----------------------------864422039439417963841331
Content-Disposition: form-data; name="first"
Content-Type: application/json

{
 "update_time":"now" 	// 升级时间，是立马升级还是定时升级，如果为定时升级，则格式为:2010-01-02 03:04
}
----------------------------864422039439417963841331
Content-Disposition: form-data; name="second"; filename="collect_update_v0.6001.bin"
Content-Type: video/mp4

... ftypisom....isomiso2avc1mp41...%mo // 升级包的内容

返回示例:
成功:
{
 "status": 200,
 "reason": "OK"
}

失败:
{
 "status": 400,			// 或者其他非200的数组
 "reason": "err reason"	// 错误描述
}
 `
	return info
}

func GetTimeHelp() (info string) {
	info = `
说明: 获取系统时间
请求url: /host/get/time
请求body: 空

返回示例:
成功:
{
 "time": "2019-05-31 11:21:28"
}

失败:
{
 "status": 400,			// 或者其他非200的数组
 "reason": "err reason"	// 错误描述
}
 `
	return info
}

// 更新系统时间的help
func UpdateTimeHelp() (info string) {
	info = `
说明: 更新系统时间
请求url: /host/update/time
请求body: 
{
 "time": "2019-05-31 11:21:28"
}

返回示例:
成功:
{
 "status": 200,			
 "reason": "OK"
}

失败:
{
 "status": 400,			// 或者其他非200的数组
 "reason": "err reason"	// 错误描述
}
 `
	return info
}

// 获取ntp配置
func GetNtpConfitureHelp() (info string) {
	info = `
说明: 获取ntp配置
请求url: /host/get/ntp
请求body: 无

返回示例:
成功:
{
 "ip": "127.0.0.0",		// ntp服务器的ip或者域名
 "port": 123,			// ntp服务器的端口
 "proofread_interval": 100,	// 校时间隔，每隔多少分钟进行一次校时，单位是分钟
 "enable": false		// 是否启用ntp校时功能
}

失败:
{
 "status": 400,			// 或者其他非200的数组
 "reason": "err reason"	// 错误描述
}
 `
	return info
}

// 更新ntp配置
func UpdateNtpConfigureHelp() (info string) {
	info = `
说明: 更新ntp配置
请求url: /host/update/ntp
请求body: 
{
	"ntp": {
		"ip": "127.0.0.0",		// ntp服务器的ip或者域名
		"port": 123,			// ntp服务器的端口
		"proofread_interval": 100, // 校时间隔，单位为分钟
		"enable": false			// 是否启用ntp校时功能
	},
	"operate": "test"		// 操作类型，如果是测试ntp服务器，则为"test",如果为设置ntp配置，则为"set"
}

返回示例:
成功:
{
 "status": 200,
 "reason": "OK"
}

失败:
{
 "status": 400,			// 或者其他非200的数组
 "reason": "err reason"	// 错误描述
}
 `
	return info
}

func ServerOperateHelp() (info string) {
	info = `
说明: 操作服务器(关机/重启)
请求url: /host/operate/server
请求body: 
{
 "operate":"stop",	// 操作类型，支持"stop"和"restart"	
 "time":"now"		// 操作时间，"now"为立即执行，如果要定时执行，则传入具体时间，格式为"2010-01-02 03:04"
}

返回示例:
成功:
{
 "status": 200,
 "reason": "OK"
}

失败:
{
 "status": 400,			// 或者其他非200的数组
 "reason": "err reason"	// 错误描述
}
 `
	return info
}

// 获取服务器的标识符
func ServerIdHelp() (info string) {
	info = `
说明: 获取服务器的标识符(为了生成license)
请求url: /host/operate/server
请求body: 空

返回示例:
成功:
{
 "id": "host=xyzdfskhewroia45",		// 服务器标识符
 "download_path": "/html/download/license/hardware.conf"	// 服务器标识文件下载路径
 "zip_download_path": "/html/download/license/hardware.zip"	// 服务器标识压缩文件下载路径
}

失败:
{
 "status": 400,			// 或者其他非200的数组
 "reason": "err reason"	// 错误描述
}
 `
	return info
}

// 更新license信息
func UpdateLicenseHelp() (info string) {
	info = `
说明: 更新license
请求url: /host/update/license
请求body: 
----------------------------431011755631093163772100
Content-Disposition: form-data; name="dfsdfsd"; filename="license"
Content-Type: application/octet-stream

license=Ojhq089rovwLxgxX7SOBqkSpOQFaoFMYcNrhRcQkI53IdySJt38O87JJM0dvRghtjHunpR5YVhMSxpViCTmutUwbHIn2qsQJYY602i+qZWYBoC93jsD1HKZi4MRxcYceq+RfxomqA2oEJjyYAxHUHsIqPOvVOB4bYFG/hkmJnQ3z2v55cA9gNclwB3MUYRXcVn/3PckAtT+gcLcDkSuvXsM8YRL5y08zZM5hI1bCKi08yDYrro6WDWKI6gijOZFqiOzl3ts/9VRoP4/SPKNHPeU+GM6M18Kwj0Zz6GAQRfp4mRqS5Pz+qakmvePL77XWmRSKU3uAVy8Pn1Lo/KBD2Q==

----------------------------431011755631093163772100--

返回示例:
成功:
{
 "status": 200,
 "reason": "OK"
}

失败:
{
 "status": 400,			// 或者其他非200的数组
 "reason": "err reason"	// 错误描述
}
 `
	return info
}

// 获取license信息
func GetLicenseInfoHelp() (info string) {
	info = `
说明: 更新license
请求url: /host/get/license
请求body: 无

返回示例:
成功:
{
 "device_capacity": 20,		// 最多接入的设备数量
 "expire": "2050-10-02",	// license过期时间
 "multi_instance": true,	// license是否支持多实例
 "valid": true,				// license是否有效
 "error": ""				// 如果license过期，则此处为过期原因
}

失败:
{
 "status": 400,			// 或者其他非200的数组
 "reason": "err reason"	// 错误描述
}
 `
	return info
}

// 获取所有时区信息
func TimeZoneInfoHelp() (info string) {
	info = `
说明: 获取所有时区信息
请求url: /host/get/time/zones/info
请求body: 无

返回示例:
成功:
{
 "time_zone":[{
  "offset":28800,	// 相对于格林威治时间的偏移(单位秒)，例如东八区，则偏移为8 * 60 * 60
  "city:"Asia/Shanghai",	// 该时区所在的城市，如果是中国标准时区，则有Asia/Shanghai和Asia/Chongqing
  "name":"GMT+08",  // 该时区在界面上显示的名字
  "code":"CST"	// 该时区的代码，中国标准时区为CST
 }]
}

失败:
{
 "status": 400,			// 或者其他非200的数组
 "reason": "err reason"	// 错误描述
}
 `
	return info
}

// 获取所有时区信息
func ServerTimeZoneHelp() (info string) {
	info = `
说明: 获取当前服务器的时区信息
请求url: /host/get/server/time/zone
请求body: 无

返回示例:
成功:
{
 "offset":28800,	// 相对于格林威治时间的偏移(单位秒)，例如东八区，则偏移为8 * 60 * 60
 "city:"Asia/Shanghai",	// 该时区所在的城市，如果是中国标准时区，则有Asia/Shanghai和Asia/Chongqing
 "name":"GMT+08",  // 该时区在界面上显示的名字
 "code":"CST"	// 该时区的代码，中国标准时区为CST
}

失败:
{
 "status": 400,			// 或者其他非200的数组
 "reason": "err reason"	// 错误描述
}
 `
	return info
}

// 获取所有时区信息
func UpdateTimeZoneHelp() (info string) {
	info = `
说明: 更新当前服务器的时区信息
请求url: /host/update/server/time/zone
请求body: 
{
 "offset":28800,	// 相对于格林威治时间的偏移(单位秒)，例如东八区，则偏移为8 * 60 * 60
 "city:"Asia/Shanghai",	// 该时区所在的城市，如果是中国标准时区，则有Asia/Shanghai和Asia/Chongqing
 "name":"GMT+08",  // 该时区在界面上显示的名字
 "code":"CST"	// 该时区的代码，中国标准时区为CST
}

返回示例:
成功:
{
 "status": 200,
 "reason": "OK"
}

失败:
{
 "status": 400,			// 或者其他非200的数组
 "reason": "err reason"	// 错误描述
}
 `
	return info
}

// 获取时间信息页面
func TimeInfoPageHelp() (info string) {
	info = `
说明: 更新当前服务器的时区信息
请求url: /host/get/time/info/page
请求body: 无

返回示例:
成功:
{
 "time_zone": {		// 当前时区信息
  "offset": 28800,
  "city": "Asia/Shanghai",
  "name": "GMT+8.00",
  "code": "CST"
 },
 "ntp": {		// ntp配置信息
  "ip": "118.24.4.66",
  "port": 123,
  "proofread_interval": 60,
  "enable": true
 },
 "time": "2019-06-11 19:21:51"	// 系统当前时间信息
}

失败:
{
 "status": 400,			// 或者其他非200的数组
 "reason": "err reason"	// 错误描述
}
 `
	return info
}

// 获取时间信息页面
func UpdateTimeInfoPageHelp() (info string) {
	info = `
说明: 更新当前服务器的时区信息
请求url: /host/update/time/info/page
请求body: 
{
 "time_zone": {		// 当前时区信息
  "offset": 28800,
  "city": "Asia/Shanghai",
  "name": "GMT+8.00",
  "code": "CST"
 },
 "ntp": {		// ntp配置信息
  "ip": "118.24.4.66",
  "port": 123,
  "proofread_interval": 60,
  "enable": true
 },
 "time": "2019-06-11 19:21:51"	// 系统要更新的时间,如果ntp中的enable为false,则此配置不起效
}

返回示例:
成功:
{
 "status": 200,		
 "reason": "OK"	
}

失败:
{
 "status": 400,			// 或者其他非200的数组
 "reason": "err reason"	// 错误描述
}
 `
	return info
}

// 查看服务器的系统日志
func GetServerLogHelp() (info string) {
	info = `
说明: 查看服务器的系统日志
请求url: /host/get/server/log
请求body: 
{
"begin":0,		// 查看的开始行数,如果是第一次请求，则填0
"end":0			// 查看的结束行数,可以一直填0
}
返回示例:
成功:
{
 "begin": 141,		// 开始的行数	
 "end": 241,		// 结束的行数
 "log": [			// 日志的内容
  "Jun 19 09:00:01 localhost systemd: Started Session 102 of user root.",
  "Jun 19 09:01:02 localhost systemd: Started Session 103 of user root.",
 ]
}

失败:
{
 "status": 400,			// 或者其他非200的数组
 "reason": "err reason"	// 错误描述
}
 `
	return info
}

// 查看服务组件的日志
func GetServiceLogHelp() (info string) {
	info = `
说明: 查看服务组件的日志
请求url: /host/get/service/log
请求body: 
{
"service":"business",	// 服务组件的名称
"begin":0,		// 查看的开始行数,如果是第一次请求，则填0
"end":0			// 查看的结束行数,可以一直填0
}
返回示例:
成功:
{
 "begin": 141,		// 开始的行数	
 "end": 241,		// 结束的行数
 "log": [			// 日志的内容
  "Jun 19 09:00:01 localhost systemd: Started Session 102 of user root.",
  "Jun 19 09:01:02 localhost systemd: Started Session 103 of user root.",
 ]
}

失败:
{
 "status": 400,			// 或者其他非200的数组
 "reason": "err reason"	// 错误描述
}
 `
	return info
}

// 下载服务器的系统日志
func DownloadServerLogHelp() (info string) {
	info = `
说明: 查看服务器的系统日志
请求url: /host/download/server/log
请求body: 空

返回示例:
成功:
{
 "log": "/html/download/log/message",			// 日志下载路径	
 "zip_log": "/html/download/log/message.zip"		// 压缩后的日志的下载路径
}

失败:
{
 "status": 400,			// 或者其他非200的数组
 "reason": "err reason"	// 错误描述
}
 `
	return info
}

// 统计CPU的使用情况
func CpuStatisticHelp() (info string) {
	info = `
说明: 统计CPU的使用情况
请求url: /host/get/cpu/statistic
请求body: 
{
 "begin": "",	// 起始时间，格式为"2019-07-02 14:51:20",也可以为空
 "end": ""		// 结束时间时间，格式为"2019-07-02 14:51:20",也可以为空
}

返回示例:
成功:
{
 "cpu": [
  {
   "time": "2019-07-02 14:33:35",	// 采样的时间
   "statistic": {
    "total_capacity": 56,			// CPU的总核数
    "total_capacity_readable": "",
    "used_capacity": 2,				// 使用了多少CPU(总共为100,如果此值为2，则使用了2%的cpu)
    "used_capacity_readable": "",
    "available_capacity": 98,		// 多少空闲的CPU
    "available_capacity_readable": ""
   }
  }]
}

失败:
{
 "status": 400,			// 或者其他非200的数组
 "reason": "err reason"	// 错误描述
}
 `
	return info
}

// 统计磁盘的使用情况
func DiskStatisticHelp() (info string) {
	info = `
说明: 统计磁盘的使用情况
请求url: /host/get/disk/statistic
请求body: 
{
 "begin": "",	// 起始时间，格式为"2019-07-02 14:51:20",也可以为空
 "end": ""		// 结束时间时间，格式为"2019-07-02 14:51:20",也可以为空
}

返回示例:
成功:
{
 "disk": [
  {
   "time": "2019-07-02 14:53:55",	//采样时间
   "statistic": [
    {
     "read_count": 0,				// 一秒内读取磁盘的次数
     "merged_read_count": 0,		// 一秒内合并读取磁盘的次数
     "write_count": 6,				// 一秒内写入磁盘的次数
     "merged_write_count": 0,		// 一秒内合并写入磁盘的次数
     "read_byte": 0,				// 一秒内读取的字节数
     "read_byte_readable": "0.0B",
     "write_byte": 552960,			// 一秒内写入的字节数
     "write_byte_readable": "540.0K",
     "read_time": 0,
     "write_time": 16,
     "iops_in_progress": 0,			// 平均每次IO操作的数据量
     "io_time": 16,					// 平均每次IO请求等待时间
     "weighted_IO": 3043041,		// 平均等待处理的IO请求队列长度
     "name": "sda",					// 磁盘的名称
     "serialNumber": "",
     "label": "",
     "util": 99                     // 磁盘的利用率百分比
    },
    {
     "read_count": 0,
     "merged_read_count": 0,
     "write_count": 0,
     "merged_write_count": 0,
     "read_byte": 0,
     "read_byte_readable": "0.0B",
     "write_byte": 0,
     "write_byte_readable": "0.0B",
     "read_time": 0,
     "write_time": 0,
     "iops_in_progress": 0,
     "io_time": 0,
     "weighted_IO": 228,
     "name": "sda2",
     "serialNumber": "",
     "label": ""
    },
    {
     "read_count": 0,
     "merged_read_count": 0,
     "write_count": 6,
     "merged_write_count": 0,
     "read_byte": 0,
     "read_byte_readable": "0.0B",
     "write_byte": 552960,
     "write_byte_readable": "540.0K",
     "read_time": 0,
     "write_time": 16,
     "iops_in_progress": 0,
     "io_time": 16,
     "weighted_IO": 3038707,
     "name": "sda5",
     "serialNumber": "",
     "label": ""
    },
    {
     "read_count": 0,
     "merged_read_count": 0,
     "write_count": 0,
     "merged_write_count": 0,
     "read_byte": 0,
     "read_byte_readable": "0.0B",
     "write_byte": 0,
     "write_byte_readable": "0.0B",
     "read_time": 0,
     "write_time": 0,
     "iops_in_progress": 0,
     "io_time": 0,
     "weighted_IO": 134857,
     "name": "sdb",
     "serialNumber": "",
     "label": ""
    }]
  }]
}

失败:
{
 "status": 400,			// 或者其他非200的数组
 "reason": "err reason"	// 错误描述
}
 `
	return info
}

// 统计网络的使用情况
func NetworkStatisticHelp() (info string) {
	info = `
说明: 统计网络的使用情况
请求url: /host/get/network/statistic
请求body: 
{
 "begin": "",	// 起始时间，格式为"2019-07-02 14:51:20",也可以为空
 "end": ""		// 结束时间时间，格式为"2019-07-02 14:51:20",也可以为空
}

返回示例:
成功:
{
 "network": [
  {
   "time": "2019-07-02 14:33:37",		// 采样的时间
   "statistic": [
    {
     "name": "enp3s0f0",			// 网卡名称
     "receive_byte": 607,			// 一秒内接收的字节数
     "receive_byte_readable": "607.0B",
     "send_byte": 332,				// 一秒内发送的字节数
     "send_byte_readable": "332.0B",
     "capacity_byte": 131072000,
     "capacity_byte_readable": "125.0M",	// 网卡的最大吞吐量
     "receive_packet": 5,		// 一秒内接收的数据包
     "send_packet": 3,			// 一秒内发送的数据包
     "receive_error": 0,		// 一秒内接收到的错误数据包的个数
     "send_error": 0,			// 一秒内发送的错误数据包的个数
     "receive_drop": 0,			// 一秒内丢弃调接收数据包的个数
     "send_drop": 0				// 一秒内丢弃掉发送数据包的个数
    },
    {
     "name": "enp3s0f1",
     "receive_byte": 0,
     "receive_byte_readable": "0.0B",
     "send_byte": 0,
     "send_byte_readable": "0.0B",
     "capacity_byte": 0,
     "capacity_byte_readable": "0.0B",
     "receive_packet": 0,
     "send_packet": 0,
     "receive_error": 0,
     "send_error": 0,
     "receive_drop": 0,
     "send_drop": 0
    }]
  }]
}

失败:
{
 "status": 400,			// 或者其他非200的数组
 "reason": "err reason"	// 错误描述
}
 `
	return info
}

// 统计内存的使用情况
func MemoryStatisticHelp() (info string) {
	info = `
说明: 统计内存的使用情况
请求url: /host/get/memory/statistic
请求body: 
{
 "begin": "",	// 起始时间，格式为"2019-07-02 14:51:20",也可以为空
 "end": ""		// 结束时间时间，格式为"2019-07-02 14:51:20",也可以为空
}

返回示例:
成功:
{
 "memory": [
  {
   "time": "2019-07-02 14:33:35",		// 采样时间
   "statistic": {
    "total_capacity": 134907924480,		// 内存的总容量
    "total_capacity_readable": "125.6G",
    "used_size": 14672211968,		// 已经使用的内存容量
    "used_size_readable": "13.7G",
    "available_size": 120235712512,	// 剩余可用的内存容量
    "available_size_readable": "112.0G"
    "buffer_size": 14672211968,		// buff的容量
    "buffer_size_readable": "13.7G",
    "cache_size": 120235712512,	// cache的容量
    "cache_size_readable": "112.0G"
   }]
}

失败:
{
 "status": 400,			// 或者其他非200的数组
 "reason": "err reason"	// 错误描述
}
 `
	return info
}

// 统计服务组件的资源使用情况
func ServiceStatisticHelp() (info string) {
	info = `
说明: 统计服务组件的资源使用
请求url: /host/get/service/statistic
请求body: 
{
 "begin": "",	// 起始时间，格式为"2019-07-02 14:51:20",也可以为空
 "end": ""		// 结束时间时间，格式为"2019-07-02 14:51:20",也可以为空
}

返回示例:
成功:
{
 "services": [
  {
   "time": "2019-07-02 15:09:17",		// 采样时间
   "statistic": [
    {
     "service": "analysis",			// 服务组件的名称
     "pid": 14025,					// 服务组件的进程ID
     "thread_count": 292,			// 服务组件的线程个数
     "cpu": 1.0982724022578325,		// 服务组件的CPU使用率(最高为100)
     "startup": "2019-07-01 14:38:14",	// 服务组件的启动时间
     "command": "",
     "memory": {
      "used_precent": 1.0254033,	// 服务组件的内存使用率(最高为100)
      "rss": 1383350272,			// 服务组件使用的物理内存
      "rss_readable": "1.3G",
      "vms": 30209196032,			// 服务组件使用的虚拟内存
      "vms_readable": "28.1G",
      "data": 0,
      "data_readable": "0.0B"
     },
     "disk": {
      "read_count": 0,				// 服务组件在一秒内读取磁盘的次数
      "write_count": 0,				// 服务组件在一秒内写入磁盘的次数
      "read_byte": 0,				// 服务组件在一秒内写入磁盘的字节数
      "read_byte_readable": "0.0B",
      "write_byte": 0,				// 服务组件在一秒内读取磁盘的字节数
      "write_byte_readable": "0.0B"
     },
     "network": {
      "name": "all",
      "receive_byte": 0,			// 服务组件在一秒内通过网络接收到的字节数
      "receive_byte_readable": "0.0B",
      "send_byte": 0,				// 服务组件在一秒内通过网络发送的字节数
      "send_byte_readable": "0.0B",
      "capacity_byte": 0,
      "capacity_byte_readable": "",
      "receive_packet": 0,			// 服务组件在一秒内通过网络接收到的数据包的个数
      "send_packet": 0,				// 服务组件在一秒内通过网络发送的数据包的个数
      "receive_error": 0,
      "send_error": 0,
      "receive_drop": 0,
      "send_drop": 0
     }
    }]
  }]
}

失败:
{
 "status": 400,			// 或者其他非200的数组
 "reason": "err reason"	// 错误描述
}
 `
	return info
}

// 获取存储配置
func GetStorageConfigure() (info string) {
	info = `
说明: 查询磁盘的存储配置(获取存储的预留空间)
请求url: /host/get/storage/configure
请求body: 无

返回示例:
成功:
{
 "storage": {
  "remove_threshold": 1099511627776,		// 存储预留的空间，单位为byte
  "remove_threshold_readable": "1.0T"
 },
 "last_remove": []							// 最近一次删除的图片的目录
}

失败:
{
 "status": 400,			// 或者其他非200的数组
 "reason": "err reason"	// 错误描述
}
 `
	return info
}

// 设置存储配置
func StorageConfigure() (info string) {
	info = `
说明: 设置磁盘的存储配置(获取存储的预留空间)
请求url: /host/storage/configure
请求body: 
{
  "remove_threshold": 1099511627776		// 设置存储的预留空间，单位为byte
}

返回示例:
成功:
{
 "status": 200,	
 "reason": "OK"	
}

失败:
{
 "status": 400,			// 或者其他非200的数组
 "reason": "err reason"	// 错误描述
}
 `
	return info
}

// 查看磁盘的状态的信息
func GetDiskStatusInfoHelp() (info string) {
	info = `
说明: 查看磁盘的状态的信息
请求url: /host/get/disk/status/info
请求body: 空

返回示例:
成功:
[
 {
  "name": "/dev/sda1",
  "type": 1,                     //磁盘的类型，1表示本地，2表示网络磁盘，3表示云存储
  "mount_point": "/boot",
  "status": true,				//磁盘的状态，true是正常，false是异常
  "mount_status": true,			//磁盘挂载的属性，true是已挂载、可读写，false是未挂载
  "total_capacity": 0,
  "total_capacity_readable": "1014M",
  "used_capacity": 0,
  "used_capacity_readable": "",
  "available_capacity": 0,
  "available_capacity_readable": "869M"
 }
]

失败:
{
 "status": 400,			// 或者其他非200的数组
 "reason": "err reason"	// 错误描述
}
 `
	return info
}

// des加密的帮助
func EncrpytDESHelp()(info string) {
	info = `
说明: DES加密
请求url: /host/encrypt/des
请求body: 
{
  "info": "123456"		// 要加密的数据
}

返回示例:
成功:
{
  "info": "uwd8sdfellkil"		// 加密后数据的base64编码
}

失败:
{
 "status": 400,			// 或者其他非200的数组
 "reason": "err reason"	// 错误描述
}
 `
	return info
}

// des解密的帮助
func DecryptDESHelp()(info string) {
	info = `
说明: DES加密
请求url: /host/decrypt/des
请求body: 
{
  "info": "uwd8sdfellkil"		// des加密后数据的base64编码
}

返回示例:
成功:
{
  "info": "123456"		// 解密之后的数据
}

失败:
{
 "status": 400,			// 或者其他非200的数组
 "reason": "err reason"	// 错误描述
}
 `
	return info
}

// 更改kafka配置的帮助
func UpdateKafkaConfigHelp()(info string) {
	info = `
说明: 更改kafka配置
请求url: /host/update/kafka/config
请求body: 
{
  "ip": "192.168.1.100",  // 需要修改的IP
  "port": 9092            // 需要修改的端口
}

返回示例:
成功:
{
  "status": 200,	
  "reason": "OK"	
}

失败:
{
 "status": 400,			// 或者其他非200的数组
 "reason": "err reason"	// 错误描述
}
 `
	return info
}

// 修改地图模式的帮助
func UpdateMapModeHelp()(info string) {
	info = `
说明: 修改地图模式
请求url: /host/update/map/mode
请求body(body由两部分组成，第一部分为json格式的在线/离线模式,第二部分为压缩包的内容): 
----------------------------864422039439417963841331
Content-Disposition: form-data; name="first"
Content-Type: application/json

{
 "mode": 1	// 1表示在线，2表示离线
}
----------------------------864422039439417963841331
Content-Disposition: form-data; name="second"; filename="amap.zip"
Content-Type: application/octet-stream

... ftypisom....isomiso2avc1mp41...%mo // 压缩包的内容

返回示例:
成功:
{
  "status": 200,	
  "reason": "OK"	
}

失败:
{
 "status": 500,			// 或者其他非200的数组
 "reason": "err reason"	// 错误描述
}
 `
	return info
}

// 操作磁盘挂载格式化的帮助
func OperateDiskHelp()(info string) {
	info = `
说明: 更改kafka配置
请求url: /host/operate/disk
请求body: 
{
  "name": "/dev/sdb",
  "type": "mount"    // mount挂载，format是格式化
}

返回示例:
成功:
{
  "status": 200,	
  "reason": "OK"	
}

失败:
{
 "status": 400,			// 或者其他非200的数组
 "reason": "err reason"	// 错误描述
}
 `
	return info
}

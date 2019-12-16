package network

import (
	"CommonModule"
	"CommonModule/message"
	"NetworkModule/HttpHelper"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/sirupsen/logrus"
	"io"
	"io/ioutil"
	"net/http"
	_ "net/http/pprof"
	"os"
	"strconv"
	"strings"
)

// 时间页面
type timePageInfo struct {
	TimeZone common.TimeZone `json:"time_zone"` // 当前时区信息
	NTP      common.NTPInfo  `json:"ntp"`
	Time     string          `json:"time"`
}

// 查看的日志信息
type logInfo struct {
	Begin int64    `json:"begin"` // 开始行数
	End   int64    `json:"end"`   // 结束行数
	Log   []string `json:"log"`   // 日志
}

// 时间范围
type timeRange struct {
	Begin string `json:"begin"` // 开始时间
	End   string `json:"end"`   // 结束时间
}

// http
type HttpService struct {
	Network     *NetworkManager        //
	httpService *common.ServiceInfo    //  http服务
	token       HttpHelper.TokenHelper // token管理
}

// 启动http服务
func (h *HttpService) Startup(info *common.ServiceInfo) error {
	h.httpService = info
	h.token = *HttpHelper.InitTokenHelper()

	// 构建监听的主机和端口
	listenInfo := ":" + strconv.Itoa(int(h.httpService.HttpPort))

	// 开始监听
	h.handle()
	err := http.ListenAndServe(listenInfo, nil)
	if err != nil {
		logrus.Errorf("listen http port %d fail, error reason:%s", h.httpService.Port, err.Error())
		os.Exit(0)
	}
	return err
}

// 设置 http handle
func (h *HttpService) handle() {
	http.HandleFunc("/host/", h.processHttp)
	http.HandleFunc("/html/", h.html)
	http.HandleFunc("/", h.other)
}

// 处理http请求
func (h *HttpService) processHttp(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")

	// 处理options
	if req.Method == "OPTIONS" {
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		w.Header().Set("Vary", "Origin, Access-Control-Request-Method, Access-Control-Request-Headers")
		if _, ok := req.Header["Access-Control-Request-Method"]; ok {
			w.Header().Set("Access-Control-Allow-Methods", req.Header.Get("Access-Control-Request-Method"))
		}
		if _, ok := req.Header["Access-Control-Request-Headers"]; ok {
			w.Header().Set("Access-Control-Allow-Headers", req.Header.Get("Access-Control-Request-Headers"))
		}
		return
	}

	// 记录请求
	h.recordEvent(w, req)

	//  升级需要特殊处理(因为body有多个multipart)
	if req.URL.Path == "/host/update/service" {

		// 检查token是否合法
		if _, err := h.checkToken(req); err != nil {
			response := h.response(401, err.Error())
			w.Write(response)
			return
		}

		h.updateService(w, req)
		return
	}

	//  上传license需要特殊处理(上传license是上传文件)
	if req.URL.Path == "/host/update/license" {

		// 检查token是否合法
		if _, err := h.checkToken(req); err != nil {
			response := h.response(401, err.Error())
			w.Write(response)
			return
		}

		h.updateLicense(w, req)
		return
	}

	//  更新地图模式需要特殊处理(因为body有多个multipart)
	if req.URL.Path == "/host/update/map/mode" {

		// 检查token是否合法
		if _, err := h.checkToken(req); err != nil {
			response := h.response(401, err.Error())
			w.Write(response)
			return
		}

		h.updateMapMode(w, req)
		return
	}

	var body string
	var response []byte
	info, err := ioutil.ReadAll(req.Body)
	if err != nil {
		logrus.Errorf("get %s body fail, error reason:%s", req.URL.Path, err.Error())
		w.Write(h.response(400, "get http body fail"))
		return
	}
	defer req.Body.Close()
	body = string(info[:])

	// 登陆需要特殊处理(因为登陆没有token)
	if req.URL.Path == "/host/login/service" {
		response = h.loginService(&body, req.RemoteAddr, req.UserAgent())
		w.Write(response)
		return
	}

	// 获取token
	var token string
	if token, err = h.checkToken(req); err != nil {
		response = h.response(401, err.Error())
		w.Write(response)
		return
	}

	// 记录
	logrus.Infof("http request, url:%s, token:%s, body:%s", req.URL.Path, token, body)

	switch req.URL.Path {
	case "/host/index": // index页面
		response = h.index(&body)
	case "/host/help": // 帮助页面
		response = h.help(&body)
	case "/host/get/ntp": // 获取ntp配置信息
		response = h.ntpServerInfo(&body)
	case "/host/get/ntp/help":
		response = h.ntpServerInfoHelp(&body)
	case "/host/update/ntp": // 更新ntp配置信息
		response = h.updateNtpServerInfo(&body)
	case "/host/update/ntp/help":
		response = h.updateNtpServerInfoHelp(&body)
	case "/host/operate/server": // 控制服务器（重启，停止等）
		response = h.serverOperate(&body)
	case "/host/operate/server/help":
		response = h.serverOperateHelp(&body)
	case "/host/get/server/log": // 查看服务器的系统日志
		response = h.getServerLog(&body)
	case "/host/get/server/log/help":
		response = h.getServerLogHelp(&body)
	case "/host/download/server/log": // 查看服务器的系统日志
		response = h.downloadServerLog(&body)
	case "/host/download/server/log/help":
		response = h.downloadServerLogHelp(&body)
	case "/host/get/service/log": // 查看服务组件的日志
		response = h.getServiceLog(&body)
	case "/host/get/service/log/help":
		response = h.getServiceLogHelp(&body)
	case "/host/update/network/configure/help":
		response = h.updateNetworkConfigureHelp(&body)
	case "/host/login/service/help":
		response = h.loginServiceHelp(&body)
	case "/host/get/time": // 获取本机时间
		response = h.getTime(&body)
	case "/host/get/time/help":
		response = h.getTimeHelp(&body)
	case "/host/update/time": // 修改本机时间
		response = h.updateTime(&body)
	case "/host/update/time/help":
		response = h.updateTimeHelp(&body)
	case "/host/stop": // 停止进程
		response = h.stop(&body)
	case "/host/get/time/info/page": // 时间信息的页面
		response = h.timeInfoPage(&body)
	case "/host/get/time/info/page/help":
		response = h.timeInfoPageHelp(&body)
	case "/host/update/time/info/page": // 更新时间信息的页面
		response = h.updateTimeInfoPage(&body)
	case "/host/update/time/info/page/help":
		response = h.updateTimeInfoPageHelp(&body)
	case "/host/service/info": // 获取host服务的信息
		response = h.hostServiceInfo(&body)
	case "/host/get/cpu/statistic": // 获取cpu的使用情况
		response = h.getCpuStatistic(&body)
	case "/host/get/cpu/statistic/help":
		response = h.getCpuStatisticHelp(&body)
	case "/host/get/disk/statistic": // 获取磁盘的使用情况
		response = h.getDiskStatistic(&body)
	case "/host/get/disk/statistic/help":
		response = h.getDiskStatisticHelp(&body)
	case "/host/get/network/statistic": // 获取网络的使用情况
		response = h.getNetworkStatistic(&body)
	case "/host/get/network/statistic/help":
		response = h.getNetworkStatisticHelp(&body)
	case "/host/get/memory/statistic": // 获取内存络的使用情况
		response = h.getMemoryStatistic(&body)
	case "/host/get/memory/statistic/help":
		response = h.getMemoryStatisticHelp(&body)
	case "/host/get/service/statistic": // 获取服务组件的资源使用情况
		response = h.getServiceStatistic(&body)
	case "/host/get/service/statistic/help":
		response = h.getServiceStatisticHelp(&body)
	case "/host/storage/configure": // 设置存储配置
		response = h.storageConfigure(&body)
	case "/host/storage/configure/help": // 设置存储配置
		response = h.storageConfigureHelp(&body)
	case "/host/get/storage/configure": // 获取存储配置
		response = h.getStorageConfigure(&body)
	case "/host/get/storage/configure/help": // 获取存储配置
		response = h.getStorageConfigureHelp(&body)
	case "/host/get/disk/status/info": // 获取磁盘状态的信息
		response = h.diskStatusInfo(&body)
	case "/host/get/disk/status/info/help":
		response = h.getDiskStatusInfoHelp(&body)
	case "/host/set/debug/info": // 修改调试配置
		response = h.setDebugInfo(&body)
	case "/host/get/debug/info": // 获取调试配置
		response = h.getDebugInfo(&body)
	case "/host/decrypt/des/help":
		response = h.decryptDESHelp(&body)
	case "/host/send/kafka/message":
		response = h.sendKafkaMessage(&body)
	default:
		response = h.response(404, "not find processor of path "+req.URL.Path)
	}
	// 记录
	logrus.Infof("http response, url:%s, token:%s, response:%s", req.URL.Path, token, string(response))
	w.Write(response)
}

// 处理index
func (h *HttpService) index(body *string) []byte {
	return []byte("Welcome to host service, I am queen!!!")
}

// help 返回帮助页面
func (h *HttpService) help(body *string) []byte {
	return []byte(HttpHelper.Help())
}

// 处理除了html,host的其他
func (h *HttpService) other(w http.ResponseWriter, req *http.Request) {
	url := req.RequestURI
	if len(url) <= 0 || url == "/" {
		url = "/index.html"
	}
	fileName := "./html" + url
	h.processHtml(fileName, w, req)
}

// 处理html
func (h *HttpService) html(w http.ResponseWriter, req *http.Request) {
	fileName := "." + req.RequestURI
	h.processHtml(fileName, w, req)
}

// 处理html
func (h *HttpService) processHtml(fileName string, w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")

	// 处理options
	if req.Method == "OPTIONS" {
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		w.Header().Set("Vary", "Origin, Access-Control-Request-Method, Access-Control-Request-Headers")
		if _, ok := req.Header["Access-Control-Request-Method"]; ok {
			w.Header().Set("Access-Control-Allow-Methods", req.Header.Get("Access-Control-Request-Method"))
		}
		if _, ok := req.Header["Access-Control-Request-Headers"]; ok {
			w.Header().Set("Access-Control-Allow-Headers", req.Header.Get("Access-Control-Request-Headers"))
		}
		return
	}
	http.ServeFile(w, req, fileName)
}

// 获取服务器硬件信息
func (h *HttpService) hardwareInfo(body *string) []byte {
	return h.response(400, "not process it")
}

// 获取服务器状态（运行状态,启动时间等)
func (h *HttpService) serverStatus(body *string) []byte {
	return h.response(400, "not process it")
}

// 获取服务组件的信息(各个服务组件的安装路径,版本信息等)
func (h *HttpService) serviceInfo(body *string) []byte {
	info, err := h.Network.processServiceInfo()
	if err == nil {
		ret, _ := json.MarshalIndent(info, "", " ")
		return ret
	}
	return h.response(400, fmt.Sprintf("get service info fail, error reason:%s", err.Error()))
}

// 获取服务组件信息的帮助
func (h *HttpService) serviceInfoHelp(body *string) []byte {
	return []byte(HttpHelper.ServiceInfoHelper())
}

// 获取服务组件的运行状态(服务组件的内存，cpu,网络的使用情况)
func (h *HttpService) serviceStatus(body *string) []byte {
	return h.response(400, "not process it")
}

//  获取host信息(生成license使用)
func (h *HttpService) serverID(body *string) []byte {
	id, file, zipPath, err := h.Network.getServerId()
	if err != nil {
		return h.response(400, err.Error())
	}

	type ServerInfo struct {
		ID      string `json:"id"`
		Path    string `json:"download_path"`     // 服务器标识符的下载路径
		ZipPath string `json:"zip_download_path"` // zip文件的下载路径
	}

	info := ServerInfo{ID: id, Path: file, ZipPath: zipPath}
	rsp, _ := json.MarshalIndent(info, "", " ")
	return rsp
}

// 获取服务器标识符的帮助
func (h *HttpService) serverIDHelp(body *string) []byte {
	return []byte(HttpHelper.ServerIdHelp())
}

// 处理下发下来的license信息
func (h *HttpService) updateLicense(w http.ResponseWriter, req *http.Request) {
	reader, err := req.MultipartReader()
	if err != nil {
		logrus.Errorf("process update service fail, read multi part fail, error:", err.Error())
		w.Write(h.response(400, err.Error()))
		return
	}

	var fileData []byte

	// 解析body
	for {
		part, err := reader.NextPart()
		if err != nil {
			if err == io.EOF {
				break
			}

			logrus.Errorf("process update license fail, parse body fail, error:", err.Error())
			w.Write(h.response(400, err.Error()))
			return
		}

		// 处理升级包内容
		if part.Header.Get("Content-Type") == "application/octet-stream" {

			// 读取文件名和文件内容
			fileData, _ = ioutil.ReadAll(part)
			break
		}
	}

	// 判断上传的license是否合法
	license := string(fileData[:])
	if len(license) > 1024 || strings.Index(license, "license") < 0 {
		w.Write(h.response(400, "invalid license"))
		return
	}

	// 更新license信息
	err = h.Network.updateLicense(license)
	if err != nil {
		w.Write(h.response(500, err.Error()))
	} else {
		w.Write(h.response(200, "OK"))
	}
}

// 下发下来的license信息的帮助
func (h *HttpService) updateLicenseHelp(body *string) []byte {
	return []byte(HttpHelper.UpdateLicenseHelp())
}

// 获取license信息(获取license中的是否有效，过期时间，支持的设备数量等)
func (h *HttpService) licenseInfo(body *string) []byte {
	license, err := h.Network.getLicenseInfo()
	if err != nil {
		return h.response(500, err.Error())
	}

	info, err := json.MarshalIndent(license, "", " ")
	return info
}

// 获取license信息的帮助
func (h *HttpService) licenseInfoHelp(body *string) []byte {
	return []byte(HttpHelper.GetLicenseInfoHelp())
}

// 删除指定日期的图片
func (h *HttpService) removeImage(body *string) []byte {
	return h.response(400, "not process it")
}

// 安装/升级服务模块
func (h *HttpService) 	updateService(w http.ResponseWriter, req *http.Request) {

	reader, err := req.MultipartReader()
	if err != nil {
		logrus.Errorf("process update service fail, read multi part fail, error:", err.Error())
		w.Write(h.response(400, err.Error()))
		return
	}

	var fileName string
	var fileData []byte
	var updateInfo message.UpdateInfo

	// 解析body
	for {
		part, err := reader.NextPart()
		if err != nil {
			if err == io.EOF {
				break
			}

			logrus.Errorf("process update service fail, parse body fail, error:", err.Error())
			w.Write(h.response(400, err.Error()))
			return
		}

		// 处理升级包内容
		if part.Header.Get("Content-Type") == "application/octet-stream" {

			// 读取文件名和文件内容
			fileName = part.FileName()
			fileData, _ = ioutil.ReadAll(part)
		} else if part.Header.Get("Content-Type") == "application/json" {

			// 读取json信息
			data, _ := ioutil.ReadAll(part)
			err = json.Unmarshal(data, &updateInfo)
			if err != nil {
				updateInfo.UpdateTime = "now"
			}
		}
	}

	// 判断传过来的数据是否合法
	if len(updateInfo.UpdateTime) <= 0 {
		w.Write(h.response(400, "invalid update info, not find update time"))
		return
	}
	if len(fileName) <= 0 {
		w.Write(h.response(400, "invalid update info, not find file name"))
		return
	}

	err = h.Network.ProcessUpdateService(updateInfo, fileName, fileData)
	if err != nil {
		if err.Error() == "success" {
			w.Write(h.response(200, "OK"))
		} else {
			w.Write(h.response(500, err.Error()))
		}
	}
}

// 安装/升级服务模块的帮助
func (h *HttpService) updateServiceHelp(body *string) []byte {
	return []byte(HttpHelper.UpdateServiceHelper())
}

// 主机信息(包括leader和worker的主机的ip,hostname, server name等)
func (h *HttpService) hostInfo(body *string) []byte {
	return h.response(400, "not process it")
}

// 修改主机信息(ip, server name等)
func (h *HttpService) updateHostInfo(body *string) []byte {
	return h.response(400, "not process it")
}

// 获取ntp服务器信息
func (h *HttpService) ntpServerInfo(body *string) []byte {
	ntp, err := h.Network.getNtpServerInfo()
	if err != nil {
		return h.response(500, err.Error())
	}
	rsp, _ := json.MarshalIndent(ntp, "", " ")
	return rsp
}

// 获取ntp服务器信息的帮助
func (h *HttpService) ntpServerInfoHelp(body *string) []byte {
	return []byte(HttpHelper.GetNtpConfitureHelp())
}

// 设置ntp服务器信息
func (h *HttpService) updateNtpServerInfo(body *string) []byte {
	type NtpControl struct {
		NTP     common.NTPInfo        `json:"ntp"`
		Operate common.NtpControlType `json:"operate"`
	}

	// 解析json
	var ntp NtpControl
	err := json.Unmarshal([]byte(*body), &ntp)
	if err != nil {
		return h.response(400, err.Error())
	}

	//  更新ntp配置
	err = h.Network.updateNtpConfigure(ntp.NTP, ntp.Operate)
	if err != nil {
		return h.response(500, err.Error())
	}
	return h.response(200, "OK")
}

// 设置ntp服务器信息的帮助
func (h *HttpService) updateNtpServerInfoHelp(body *string) []byte {
	return []byte(HttpHelper.UpdateNtpConfigureHelp())
}

//  服务器控制（服务器的关机/重启,包括定时关机/重启）
func (h *HttpService) serverOperate(body *string) []byte {
	type SystemOperate struct {
		Operate common.OperateType `json:"operate"` // 操作类型
		Time    string             `json:"time"`    // 操作时间
	}

	// 反序列化json
	var operate SystemOperate
	err := json.Unmarshal([]byte(*body), &operate)
	if err != nil {
		return h.response(400, err.Error())
	}

	err = h.Network.systemOperate(operate.Operate, operate.Time)
	if err != nil {
		return h.response(500, err.Error())
	}
	return h.response(200, "OK")
}

func (h *HttpService) serverOperateHelp(body *string) []byte {
	return []byte(HttpHelper.ServerOperateHelp())
}

// 服务模块的控制(服务模块的关机/重启，包括定时关机/重启,注意，也可以控制本服务)
func (h *HttpService) serviceOperate(body *string) []byte {
	var control common.EntityControl
	err := json.Unmarshal([]byte(*body), &control)
	if err != nil {
		return h.response(400, fmt.Sprintf("parse json fail, error:%s", err.Error()))
	}

	err = h.Network.processControlService(control)
	if err != nil {
		return h.response(500, fmt.Sprintf("process control service fail, error:%s", err.Error()))
	} else {
		return h.response(200, "OK")
	}
}

// 服务控制的帮助文档
func (h *HttpService) serviceOperateHelp(body *string) []byte {
	return []byte(HttpHelper.ControlServiceHelper())
}

// 查看服务器的系统日志
func (h *HttpService) getServerLog(body *string) []byte {
	type LogRange struct {
		Begin int64 `json:"begin"`
		End   int64 `json:"end"`
	}

	var rng LogRange
	err := json.Unmarshal([]byte(*body), &rng)
	if err != nil {
		return h.response(400, err.Error())
	}

	// 获取日志
	begin, end, logs, err := h.Network.getServerLog(rng.Begin, rng.End)
	if err != nil {
		return h.response(500, err.Error())
	}

	// 生成json
	log := logInfo{Begin: begin, End: end, Log: logs}
	data, _ := json.MarshalIndent(log, "", " ")
	return data
}

func (h *HttpService) getServerLogHelp(body *string) []byte {
	return []byte(HttpHelper.GetServerLogHelp())
}

// 下载服务器的系统日志
func (h *HttpService) downloadServerLog(body *string) []byte {
	type ServerLog struct {
		Log    string `json:"log"`
		ZipLog string `json:"zip_log"`
	}

	// 获取日志
	log, zipLog, err := h.Network.downloadServerLog()
	if err != nil {
		return h.response(500, err.Error())
	}

	// 生成json
	info := ServerLog{Log: log, ZipLog: zipLog}
	data, _ := json.MarshalIndent(info, "", " ")
	return data
}

func (h *HttpService) downloadServerLogHelp(body *string) []byte {
	return []byte(HttpHelper.DownloadServerLogHelp())
}

// 查看服务组件的系统日志
func (h *HttpService) getServiceLog(body *string) []byte {
	type Range struct {
		Service string `json:"service"`
		Begin   int64  `json:"begin"`
		End     int64  `json:"end"`
	}

	var rng Range
	err := json.Unmarshal([]byte(*body), &rng)
	if err != nil {
		return h.response(400, err.Error())
	}

	// 获取日志
	begin, end, logs, err := h.Network.getServiceLog(rng.Service, rng.Begin, rng.End)
	if err != nil {
		return h.response(500, err.Error())
	}

	// 生成json
	log := logInfo{Begin: begin, End: end, Log: logs}
	data, _ := json.MarshalIndent(log, "", " ")
	return data
}

func (h *HttpService) getServiceLogHelp(body *string) []byte {
	return []byte(HttpHelper.GetServiceLogHelp())
}

// 下载服务模块的日志
func (h *HttpService) downloadServiceLog(body *string) []byte {
	type DownloadInfo struct {
		ServiceName string `json:"name"`
	}

	// 解析json
	var download DownloadInfo
	err := json.Unmarshal([]byte(*body), &download)
	if err != nil {
		return h.response(200, "parse json fail")
	}

	// 下载日志
	log, err := h.Network.processDownloadServiceLog(download.ServiceName)
	if err != nil {

		// 获取日志失败
		return h.response(500, fmt.Sprintf("download %s fail, error reason:%s", download.ServiceName, err.Error()))
	} else {

		// 生成回应
		type DownloadResponse struct {
			Status int    `json:"status"` // 状态码
			Reason string `json:"reason"` //
			Path   string `json:"path"`
		}
		data, _ := json.MarshalIndent(DownloadResponse{Status: 200, Reason: "OK", Path: log}, "", " ")
		return data
	}
}

// 下载服务模块日志的帮助
func (h *HttpService) downloadServiceLogHelp(body *string) []byte {
	return []byte(HttpHelper.DownloadServiceLogHelp())
}

// 获取网络配置
func (h *HttpService) networkConfigure(body *string, localIp string) []byte {
	ip, err := h.Network.GetNetworkConfigure()
	if err != nil {
		return h.response(500, err.Error())
	} else {
		// 对数据进行排序(当前的IP排在最前面)
		ip = h.sortIp(ip, localIp)

		info, _ := json.MarshalIndent(ip, "", " ")
		return info
	}
}

// 对ip进行排序
func (h *HttpService) sortIp(confs common.NetworkConfigure, localIp string) common.NetworkConfigure {
	for i, conf := range confs.Network {
		if conf.IPv6.Enable || conf.IPv4.Enable {
			if len(conf.IPv4.IP) > 0 && strings.Index(localIp, conf.IPv4.IP) >= 0 {
				confs.Network[0], confs.Network[i] = confs.Network[i], confs.Network[0]
				break
			}

			if len(conf.IPv6.IP) > 0 && strings.Index(localIp, conf.IPv6.IP) >= 0 {
				confs.Network[0], confs.Network[i] = confs.Network[i], confs.Network[0]
				break
			}
		}
	}
	return confs
}

// 获取网络配置的帮助信息
func (h *HttpService) networkConfigureHelp(body *string) []byte {
	return []byte(HttpHelper.NetworkConfigureHelp())
}

// 设置网络配置
func (h *HttpService) updateNetworkConfigure(body *string) []byte {

	// 反序列化json
	var configure common.NetworkInterface
	err := json.Unmarshal([]byte(*body), &configure)
	if err != nil {
		logrus.Errorf("parse update network configure fail, body:%s", body)
		return h.response(400, err.Error())
	}
	logrus.Infof("%s", configure.String())

	err = HttpHelper.VerifyNetworkConfigure(configure)
	if err != nil {
		return h.response(400, err.Error())
	}
	// 更新ip
	err = h.Network.UpdateNetworkConfigure(configure)
	if err != nil {
		return h.response(500, err.Error())
	} else {
		return h.response(200, "OK")
	}

}

// 更新网络配置的帮助
func (h *HttpService) updateNetworkConfigureHelp(body *string) []byte {
	return []byte(HttpHelper.UpdateNetworkConfigureHelp())
}

// 登陆
func (h *HttpService) loginService(body *string, ip, userAgent string) []byte {
	var user common.UserInfo
	err := json.Unmarshal([]byte(*body), &user)
	if err != nil {
		return h.response(400, "parse json error")
	}

	err = h.Network.processLogin(user)
	if err == nil {
		token := h.token.CreateToken(ip)
		logrus.Infof("add token:%s", token)
		return h.loginResponse(200, "OK", token)
	} else {
		return h.response(403, "invalid user or password")
	}
}

// 登陆的help
func (h *HttpService) loginServiceHelp(body *string) []byte {
	helper := HttpHelper.LoginServiceHelp()
	return []byte(helper)
}

// 登出
func (h *HttpService) logoutService(body *string, token string) []byte {
	logrus.Infof("remove token:%s", token)
	h.token.RemoveToken(token)
	return h.response(200, "OK")
}

// 登出的help
func (h *HttpService) logoutServiceHelp(body *string) []byte {
	return []byte(HttpHelper.LogoutServiceHelp())
}

// 修改密码
func (h *HttpService) changePassword(body *string) []byte {
	var pwd common.ChangePassword
	err := json.Unmarshal([]byte(*body), &pwd)
	if err != nil {
		return h.response(400, "parse json error")
	}

	err = h.Network.processChangePassword(pwd)
	if err == nil {
		return h.response(200, "OK")
	} else {
		return h.response(403, err.Error())
	}
}

// 修改密码的help
func (h *HttpService) changePasswordHelp(body *string) []byte {
	return []byte(HttpHelper.ChangePasswordHelp())
}

// 获取系统时间
func (h *HttpService) getTime(body *string) []byte {
	time, _ := h.Network.getTime()
	data, _ := json.MarshalIndent(struct {
		Time string `json:"time"`
	}{Time: time}, "", " ")
	return data
}

// 获取系统时间的帮助
func (h *HttpService) getTimeHelp(body *string) []byte {
	return []byte(HttpHelper.GetTimeHelp())
}

// 设置系统时间
func (h *HttpService) updateTime(body *string) []byte {
	tmp := struct {
		Time string `json:"time"`
	}{}
	err := json.Unmarshal([]byte(*body), &tmp)
	if err != nil {
		return h.response(400, fmt.Sprintf("parse json fail, error reason:%s", err.Error()))
	}

	err = h.Network.updateTime(tmp.Time)
	if err != nil {
		return h.response(500, err.Error())
	}
	return h.response(200, "OK")
}

// 设置系统时间的帮助
func (h *HttpService) updateTimeHelp(body *string) []byte {
	return []byte(HttpHelper.UpdateTimeHelp())
}

// 获取所有的时区信息
func (h *HttpService) timeZonesInfo(body *string) []byte {
	info, err := h.Network.getTimeZonesInfo()
	if err != nil {
		return h.response(500, err.Error())
	} else {
		type TimeZones struct {
			TimeZone []common.TimeZone `json:"time_zone"`
		}

		zones := TimeZones{TimeZone: info}
		result, _ := json.MarshalIndent(zones, "", " ")
		return result
	}
}

func (h *HttpService) timeZonesInfoHelp(body *string) []byte {
	return []byte(HttpHelper.TimeZoneInfoHelp())
}

// 获取当前服务器的时区
func (h *HttpService) serverTimeZone(body *string) []byte {
	info, err := h.Network.getTimeZone()
	if err != nil {
		return h.response(500, err.Error())
	} else {
		data, _ := json.MarshalIndent(info, "", " ")
		return data
	}
}

func (h *HttpService) serverTimeZoneHelp(body *string) []byte {
	return []byte(HttpHelper.ServerTimeZoneHelp())
}

// 更新时区信息
func (h *HttpService) updateTimeZone(body *string) []byte {

	// 反序列化json
	var zone common.TimeZone
	err := json.Unmarshal([]byte(*body), &zone)
	if err != nil {
		return h.response(400, err.Error())
	}

	// 更新时区信息
	err = h.Network.updateTimeZone(zone)
	if err != nil {
		return h.response(400, err.Error())
	} else {
		return h.response(200, "OK")
	}
}

func (h *HttpService) updateTimeZoneHelp(body *string) []byte {
	return []byte(HttpHelper.UpdateTimeZoneHelp())
}

// 查看所有token
func (h *HttpService) getToken(body *string) []byte {
	return h.token.Json()
}

// 获取时间信息的页面
func (h *HttpService) timeInfoPage(body *string) []byte {

	// 获取当前时区信息
	zone, zoneErr := h.Network.getTimeZone()
	ntp, ntpErr := h.Network.getNtpServerInfo()
	time, timeErr := h.Network.getTime()

	if zoneErr == nil && ntpErr == nil && timeErr == nil {
		timeInfo := timePageInfo{TimeZone: zone, NTP: ntp, Time: time}
		result, _ := json.MarshalIndent(timeInfo, "", " ")
		return result
	}

	return h.response(500, "get time info error")
}

func (h *HttpService) timeInfoPageHelp(body *string) []byte {
	return []byte(HttpHelper.TimeInfoPageHelp())
}

//  更新时间信息的页面
func (h *HttpService) updateTimeInfoPage(body *string) []byte {
	var timeInfo timePageInfo
	err := json.Unmarshal([]byte(*body), &timeInfo)
	if err != nil {
		return h.response(400, err.Error())
	}

	// 更新时区信息
	if len(timeInfo.TimeZone.City) > 0 {
		h.Network.updateTimeZone(timeInfo.TimeZone)
	}

	// 更新ntp信息
	h.Network.updateNtpConfigure(timeInfo.NTP, common.Ntp_Control_Set)

	// 设置时间
	if !timeInfo.NTP.Enable {
		h.Network.updateTime(timeInfo.Time)
	}

	return h.response(200, "OK")
}

func (h *HttpService) updateTimeInfoPageHelp(body *string) []byte {
	return []byte(HttpHelper.UpdateTimeInfoPageHelp())
}

// 获取host服务的信息
func (h *HttpService) hostServiceInfo(body *string) (data []byte) {

	// 获取主机服务的信息
	hostStartup, systemStartup, gitBranch, gitCommit, err := h.Network.getHostServiceInfo()
	if err == nil {
		rsp := struct {
			Startup struct {
				HostService string `json:"host_service"`
				System      string `json:"system"`
			} `json:"startup"`

			Git struct {
				Branch   string `json:"branch"`
				CommitId string `json:"commit_id"`
			} `json:"git"`
		}{}
		rsp.Startup.System = systemStartup.Format("2006-01-02 15:04:05")
		rsp.Startup.HostService = hostStartup.Format("2006-01-02 15:04:05")
		rsp.Git.Branch = gitBranch
		rsp.Git.CommitId = gitCommit

		data, _ = json.MarshalIndent(rsp, "", " ")
		return
	} else {
		return h.response(500, err.Error())
	}
}

// 获取cpu的使用情况
func (h *HttpService) getCpuStatistic(body *string) []byte {

	// 解析时间范围信息
	var times timeRange
	err := json.Unmarshal([]byte(*body), &times)
	if err != nil {
		return h.response(400, err.Error())
	}

	// 获取资源使用情况
	samplingTimes, statistics, err := h.Network.getCpuStatistic(times.Begin, times.End)
	if err != nil {
		return h.response(500, err.Error())
	}

	// 统计信息转换为json
	type statInfo struct {
		Time      string              `json:"time"`
		Statistic common.CapacityInfo `json:"statistic"`
	}
	var rsp = struct {
		Network []statInfo `json:"cpu"`
	}{}

	for index, samplingTime := range samplingTimes {
		rsp.Network = append(rsp.Network, statInfo{Time: samplingTime.Format("2006-01-02 15:04:05"), Statistic: statistics[index]})
	}

	data, _ := json.MarshalIndent(rsp, "", " ")
	return data
}

func (h *HttpService) getCpuStatisticHelp(body *string) []byte {
	return []byte(HttpHelper.CpuStatisticHelp())
}

// 获取磁盘的使用情况
func (h *HttpService) getDiskStatistic(body *string) []byte {

	// 解析时间范围信息
	var times timeRange
	err := json.Unmarshal([]byte(*body), &times)
	if err != nil {
		return h.response(400, err.Error())
	}

	// 获取资源使用情况
	samplingTimes, statistics, err := h.Network.getDiskStatistic(times.Begin, times.End)
	if err != nil {
		return h.response(500, err.Error())
	}

	// 统计信息转换为json
	type statInfo struct {
		Time      string                     `json:"time"`
		Statistic []common.DiskStatisticInfo `json:"statistic"`
	}
	var rsp = struct {
		Network []statInfo `json:"disk"`
	}{}

	for index, samplingTime := range samplingTimes {
		rsp.Network = append(rsp.Network, statInfo{Time: samplingTime.Format("2006-01-02 15:04:05"), Statistic: statistics[index]})
	}

	data, _ := json.MarshalIndent(rsp, "", " ")
	return data
}

func (h *HttpService) getDiskStatisticHelp(body *string) []byte {
	return []byte(HttpHelper.DiskStatisticHelp())
}

// 获取磁网络的使用情况
func (h *HttpService) getNetworkStatistic(body *string) []byte {

	// 解析时间范围信息
	var times timeRange
	err := json.Unmarshal([]byte(*body), &times)
	if err != nil {
		return h.response(400, err.Error())
	}

	// 获取资源使用情况
	samplingTimes, statistics, err := h.Network.getNetworkStatistic(times.Begin, times.End)
	if err != nil {
		return h.response(500, err.Error())
	}

	// 统计信息转换为json
	type statInfo struct {
		Time      string                        `json:"time"`
		Statistic []common.NetworkStatisticInfo `json:"statistic"`
	}
	var rsp = struct {
		Network []statInfo `json:"network"`
	}{}

	for index, samplintTime := range samplingTimes {
		rsp.Network = append(rsp.Network, statInfo{Time: samplintTime.Format("2006-01-02 15:04:05"), Statistic: statistics[index]})
	}

	data, _ := json.MarshalIndent(rsp, "", " ")
	return data
}

func (h *HttpService) getNetworkStatisticHelp(body *string) []byte {
	return []byte(HttpHelper.NetworkStatisticHelp())
}

// 获取内存的使用情况
func (h *HttpService) getMemoryStatistic(body *string) []byte {

	// 解析时间范围信息
	var times timeRange
	err := json.Unmarshal([]byte(*body), &times)
	if err != nil {
		return h.response(400, err.Error())
	}

	// 获取资源使用情况
	samplingTimes, statistics, err := h.Network.getMemoryStatistic(times.Begin, times.End)
	if err != nil {
		return h.response(500, err.Error())
	}

	// 统计信息转换为json
	type statInfo struct {
		Time      string                     `json:"time"`
		Statistic common.MemoryStatisticInfo `json:"statistic"`
	}
	var rsp = struct {
		Network []statInfo `json:"memory"`
	}{}

	for index, samplingTime := range samplingTimes {
		rsp.Network = append(rsp.Network, statInfo{Time: samplingTime.Format("2006-01-02 15:04:05"), Statistic: statistics[index]})
	}

	data, _ := json.MarshalIndent(rsp, "", " ")
	return data
}

func (h *HttpService) getMemoryStatisticHelp(body *string) []byte {
	return []byte(HttpHelper.MemoryStatisticHelp())
}

// 获取服务组件的资源使用情况
func (h *HttpService) getServiceStatistic(body *string) []byte {

	// 解析时间范围信息
	var times timeRange
	err := json.Unmarshal([]byte(*body), &times)
	if err != nil {
		return h.response(400, err.Error())
	}

	// 获取资源使用情况
	samplingTimes, statistics, err := h.Network.getServiceStatistic(times.Begin, times.End)
	if err != nil {
		return h.response(500, err.Error())
	}

	// 统计信息转换为json
	type statInfo struct {
		Time      string               `json:"time"`
		Statistic []common.ProcessInfo `json:"statistic"`
	}
	var rsp = struct {
		Services []statInfo `json:"services"`
	}{}

	for index, samplingTime := range samplingTimes {
		rsp.Services = append(rsp.Services, statInfo{Time: samplingTime.Format("2006-01-02 15:04:05"), Statistic: statistics[index]})
	}

	data, _ := json.MarshalIndent(rsp, "", " ")
	return data
}

func (h *HttpService) getServiceStatisticHelp(body *string) []byte {
	return []byte(HttpHelper.ServiceStatisticHelp())
}

// 设置调试信息
func (h *HttpService) setDebugInfo(body *string) []byte {
	var debug common.DebugInfo
	err := json.Unmarshal([]byte(*body), &debug)
	if err != nil {
		return h.response(400, err.Error())
	}

	// 设置调试信息
	err = h.Network.setDebugInfo(debug)
	if err != nil {
		return h.response(500, err.Error())
	}

	return h.response(200, "OK")
}

// 获取调试信息
func (h *HttpService) getDebugInfo(body *string) []byte {
	debug := h.Network.getDebugInfo()
	data, err := json.MarshalIndent(debug, "", " ")
	if err != nil {
		return h.response(500, err.Error())
	}

	return data
}

// 设置存储配置
func (h *HttpService) storageConfigure(body *string) []byte {
	var storage common.StorageConfigure
	err := json.Unmarshal([]byte(*body), &storage)
	if err != nil {
		return h.response(400, err.Error())
	}

	// 更新存储配置
	h.Network.storageConfigure(storage)
	return h.response(200, "OK")
}

func (h *HttpService) storageConfigureHelp(body *string) []byte {
	return []byte(HttpHelper.StorageConfigure())
}

// 获取存储配置
func (h *HttpService) getStorageConfigure(body *string) []byte {
	threshold, remove, err := h.Network.getStorageConfigure()

	if err != nil {
		return h.response(500, err.Error())
	} else {
		rsp := struct {
			Storage    common.StorageConfigure `json:"storage"`
			LastRemove []string                `json:"last_remove"`
		}{}
		rsp.Storage = common.StorageConfigure{RemoveThreshold: threshold, RemoveThresholdReadable: common.ReabableSize(uint64(threshold))}
		rsp.LastRemove = remove

		// 转换为json
		rspData, _ := json.MarshalIndent(rsp, "", " ")
		return rspData
	}
}

// 获取存储配置
func (h *HttpService) getStorageConfigureHelp(body *string) []byte {
	return []byte(HttpHelper.GetStorageConfigure())
}

//获取磁盘状态的信息
func (h *HttpService) diskStatusInfo(body *string) []byte {
	diskStatus, err := h.Network.getDiskInfo()
	if err != nil {
		return h.response(500, err.Error())
	}
	rsp, _ := json.MarshalIndent(diskStatus, "", " ")
	return rsp
}

func (h *HttpService) getDiskStatusInfoHelp(body *string) []byte {
	return []byte(HttpHelper.GetDiskStatusInfoHelp())
}

//操作磁盘进行格式化挂载
func (h *HttpService) operationDisk(body *string) []byte {
	var operate common.OperationDisk
	err := json.Unmarshal([]byte(*body), &operate)
	if err != nil {
		return h.response(400, err.Error())
	}

	// 操作磁盘配置
	//h.Network.storageConfigure(storage)
	//return h.response(200, "OK")
	err = h.Network.getOperationDisk(operate)
	if err != nil {
		return h.response(500, err.Error())
	}
	return h.response(200, "OK")
	//if err != nil {
	//	return h.response(500, err.Error())
	//}
	//rsp, _ := json.MarshalIndent(operationStatus, "", " ")
	//return rsp
}

func (h *HttpService) operationDiskHelp(body *string) []byte {
	return []byte(HttpHelper.OperateDiskHelp())
}

// des加密
func (h *HttpService) encryptDES(body *string) []byte {
	type JsonInfo struct {
		Info string `json:"info"`
	}

	// 解析json
	var jsonData JsonInfo
	err := json.Unmarshal([]byte(*body), &jsonData)
	if err != nil {
		return h.response(400, "parse json fail")
	}

	// DES加密
	info := h.Network.encryptDES(jsonData.Info)
	var encrypt = JsonInfo{Info: info}
	data, _ := json.MarshalIndent(encrypt, "", "")
	return data
}

func (h *HttpService) encrpytDESHelp(body *string) []byte {
	return []byte(HttpHelper.EncrpytDESHelp())
}

// des 解密
func (h *HttpService) decryptDES(body *string) []byte {
	type JsonInfo struct {
		Info string `json:"info"`
	}

	// 解析json
	var jsonData JsonInfo
	err := json.Unmarshal([]byte(*body), &jsonData)
	if err != nil {
		return h.response(400, "parse json fail")
	}

	// DES解密
	info := h.Network.decryptDES(jsonData.Info)
	var encrypt = JsonInfo{Info: info}
	data, _ := json.MarshalIndent(encrypt, "", "")
	return data
}

// 发送 kafka消息
func (h *HttpService) sendKafkaMessage(body *string) []byte {
	type Message struct {
		Topic string `json:"topic"`
		Body  string `json:"body"`
	}

	var message Message
	err := json.Unmarshal([]byte(*body), &message)
	if err != nil {
		return h.response(400, "parse json fail")
	}

	// base64解码
	data, err := base64.StdEncoding.DecodeString(message.Body)
	if err != nil {
		return h.response(400, "base64 parse body fail")
	}

	h.Network.sendKafkaMessage(message.Topic, string(data[:]))
	return h.response(200, "OK")
}

func (h *HttpService) decryptDESHelp(body *string) []byte {
	return []byte(HttpHelper.DecryptDESHelp())
}

//修改kafka配置文件
func (h *HttpService) updateKafkaConfig(body *string) []byte {
	var update common.KafkaInfo
	err := json.Unmarshal([]byte(*body), &update)
	if err != nil {
		return h.response(400, err.Error())
	}

	err = h.Network.updateKafkaConfig(update)
	if err != nil {
		return h.response(500, err.Error())
	}
	return h.response(200, "OK")

}

func (h *HttpService) updateKafkaConfigHelp(body *string) []byte {
	return []byte(HttpHelper.UpdateKafkaConfigHelp())
}

//获取kafka配置文件的信息
func (h *HttpService) kafkaConfigInfo(body *string) []byte {
	infoStatus, err := h.Network.getKafkaConfigInfo()
	if err != nil {
		return h.response(500, err.Error())
	}
	rsp, _ := json.MarshalIndent(infoStatus, "", " ")
	return rsp
}

//func (h *HttpService) kafkaConfigInfoHelp(body *string) []byte {
//	return []byte(HttpHelper.GetkafkaConfigInfoHelp())
//}

//修改地图模式
func (h *HttpService) updateMapMode(w http.ResponseWriter, req *http.Request) {

	reader, err := req.MultipartReader()
	if err != nil {
		logrus.Errorf("process update map mode fail, read multi part fail, error:", err.Error())
		w.Write(h.response(400, err.Error()))
		return
	}

	var fileName string
	var fileData []byte
	var updateInfo message.MapInfo

	// 解析body
	for {
		part, err := reader.NextPart()
		if err != nil {
			if err == io.EOF {
				break
			}

			logrus.Errorf("process update map mode fail, parse body fail, error:", err.Error())
			w.Write(h.response(400, err.Error()))
			return
		}

		// 处理压缩包内容
		if part.Header.Get("Content-Type") == "application/x-zip-compressed" {

			// 读取文件名和文件内容
			fileName = part.FileName()
			fileData, _ = ioutil.ReadAll(part)
		} else if part.Header.Get("Content-Type") == "application/json" {

			// 读取json信息
			data, _ := ioutil.ReadAll(part)
			json.Unmarshal(data, &updateInfo)
		}
	}

	//检查传的值是否正常
	if updateInfo.Mode > 2 || updateInfo.Mode < 1 {
		w.Write(h.response(400, "invalid mode info"))
		return
	}

	if updateInfo.Mode == 2  && fileName != "amap.zip" {
		w.Write(h.response(400, "invalid zip name"))
		return
	}

	//地图离线
	if updateInfo.Mode == 2 {
		err = h.Network.ProcessUpdateMapMode(updateInfo, fileName, fileData)
		//if err != nil {
		//	if err.Error() == "success" {
		//		w.Write(h.response(200, "OK"))
		//	} else {
		//		w.Write(h.response(500, err.Error()))
		//	}
		//}
		if err != nil {
			w.Write(h.response(500, "fail"))
		} else {
			w.Write(h.response(200, "OK"))
		}
	}

	//地图在线
	if updateInfo.Mode == 1 {
		resultMapMode,_ := common.CommondResult("cat /mars/web/webconfig.js |grep 'window.IS_OFFLINE_MAP' |awk -F'=' '{print $NF}' |awk -F';' '{print $1}'")
		resultMapModeSlice := strings.Split(strings.Trim(resultMapMode, "\n"), "\n")
		_,err = common.CommondResult(fmt.Sprintf(`sed -i "s/window.IS_OFFLINE_MAP =%s/window.IS_OFFLINE_MAP = %s/g" /mars/web/webconfig.js`, resultMapModeSlice[0], "false"))
		if err != nil {
			w.Write(h.response(500, "fail"))
		} else {
			w.Write(h.response(200, "OK"))
		}
	}
}

func (h *HttpService) updateMapModeHelp(body *string) []byte {
	return []byte(HttpHelper.UpdateMapModeHelp())
}

//获取地图模式的信息
func (h *HttpService) mapModeInfo(body *string) []byte {
	infoStatus, err := h.Network.getMapModeInfo()
	if err != nil {
		return h.response(500, err.Error())
	}
	rsp, _ := json.MarshalIndent(infoStatus, "", " ")
	return rsp
}

//func (h *HttpService) kafkaConfigInfoHelp(body *string) []byte {
//	return []byte(HttpHelper.GetkafkaConfigInfoHelp())
//}

// 停止进程
func (h *HttpService) stop(body *string) []byte {
	h.Network.stopHostService()
	return h.response(200, "OK")
}

// 设置回应
func (h *HttpService) response(state http.ConnState, reason string) (rsp []byte) {
	type StdHttpResponse struct {
		Status int    `json:"status"` // 状态码
		Reason string `json:"reason"` //
	}
	rsp, _ = json.MarshalIndent(StdHttpResponse{Status: int(state), Reason: reason}, "", " ")
	return rsp
}

// 登陆回应
func (h *HttpService) loginResponse(state http.ConnState, erason, token string) (rsp []byte) {
	type LoginResponse struct {
		Status int    `json:"status"` // 状态码
		Reason string `json:"reason"` //
		Token  string `json:"token"`
	}
	rsp, _ = json.MarshalIndent(LoginResponse{Status: int(state), Reason: erason, Token: token}, "", " ")
	return rsp
}

func (h *HttpService) checkToken(req *http.Request) (token string, err error) {

	// 如果发现请求来自postman,则放行(方便自己用postman进行调试)
	if agent := req.Header.Get("User-Agent"); agent == "Panda" {
		return "123456", nil
	}

	// 如果发现以help结尾，则不用进行校验
	if strings.HasSuffix(req.URL.Path, "/help") {
		return "987654", nil
	}

	if token := req.Header.Get("Token"); len(token) <= 0 {
		return "", fmt.Errorf("not find token")
	} else {
		if !h.token.IsExist(token) {
			return "", fmt.Errorf("invalid token %s", token)
		}

		return token, nil
	}
}

// 记录请求
func (h *HttpService) recordEvent(w http.ResponseWriter, req *http.Request) {

	// 获取token
	token := req.Header.Get("Token")
	url := req.RequestURI
	ip := req.RemoteAddr

	// 结构化为json
	var record = struct {
		Token string `json:"token"`
		Url   string `json:"url"`
		IP    string `json:"ip"`
	}{Token: token, Url: url, IP: ip}

	data, _ := json.Marshal(record)

	// 发送消息
	h.Network.recordEvent("http", string(data[:]))
}

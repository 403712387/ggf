package network

import (
	"CommonModule"
	"encoding/json"
	"fmt"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
	_ "net/http/pprof"
	"os"
	"strconv"
)

// 时间页面
type timePageInfo struct {
	NTP  common.NTPInfo `json:"ntp"`
	Time string         `json:"time"`
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
	Network     *NetworkManager     //
	httpService *common.ServiceInfo //  http服务
}

// 启动http服务
func (h *HttpService) Startup(info *common.ServiceInfo) error {
	h.httpService = info

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
	http.HandleFunc("/", h.processHttp)
}

// 处理http请求
func (h *HttpService) processHttp(w http.ResponseWriter, req *http.Request) {

	// 记录请求
	h.recordEvent(w, req)

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

	// 记录
	logrus.Infof("http request, url:%s, body:%s", req.URL.Path, body)

	switch req.URL.Path {
	case "/index": // index页面
		response = h.index(&body)
	case "/get/ntp": // 获取ntp配置信息
		response = h.ntpServerInfo(&body)
	case "/update/ntp": // 更新ntp配置信息
		response = h.updateNtpServerInfo(&body)
	case "/get/time": // 获取本机时间
		response = h.getTime(&body)
	case "/update/time": // 修改本机时间
		response = h.updateTime(&body)
	case "/service/info": // 获取ggf服务的信息
		response = h.ggfServiceInfo(&body)
	case "/get/cpu/statistic": // 获取cpu的使用情况
		response = h.getCpuStatistic(&body)
	case "/get/disk/statistic": // 获取磁盘的使用情况
		response = h.getDiskStatistic(&body)
	case "/get/network/statistic": // 获取网络的使用情况
		response = h.getNetworkStatistic(&body)
	case "/get/memory/statistic": // 获取内存络的使用情况
		response = h.getMemoryStatistic(&body)
	case "/set/debug/info": // 修改调试配置
		response = h.setDebugInfo(&body)
	case "/get/debug/info": // 获取调试配置
		response = h.getDebugInfo(&body)
	default:
		response = h.response(404, "not find processor of path "+req.URL.Path)
	}
	// 记录
	logrus.Infof("http response, url:%s, response:%s", req.URL.Path, string(response))
	w.Write(response)
}

// 处理index
func (h *HttpService) index(body *string) []byte {
	return []byte("Welcome to ggf service, I am queen!!!")
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

// 获取ntp服务器信息
func (h *HttpService) ntpServerInfo(body *string) []byte {
	ntp, err := h.Network.getNtpServerInfo()
	if err != nil {
		return h.response(500, err.Error())
	}
	rsp, _ := json.MarshalIndent(ntp, "", " ")
	return rsp
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

// 获取系统时间
func (h *HttpService) getTime(body *string) []byte {
	time, _ := h.Network.getTime()
	data, _ := json.MarshalIndent(struct {
		Time string `json:"time"`
	}{Time: time}, "", " ")
	return data
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

// 获取host服务的信息
func (h *HttpService) ggfServiceInfo(body *string) (data []byte) {

	// 获取主机服务的信息
	hostStartup, systemStartup, gitBranch, gitCommit, err := h.Network.getGgfServiceInfo()
	if err == nil {
		rsp := struct {
			Startup struct {
				HostService string `json:"ggf_service"`
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

// 记录请求
func (h *HttpService) recordEvent(w http.ResponseWriter, req *http.Request) {

	// 获取token
	url := req.RequestURI
	ip := req.RemoteAddr

	// 结构化为json
	var record = struct {
		Url string `json:"url"`
		IP  string `json:"ip"`
	}{Url: url, IP: ip}

	data, _ := json.Marshal(record)

	// 发送消息
	h.Network.recordEvent("http", string(data[:]))
}

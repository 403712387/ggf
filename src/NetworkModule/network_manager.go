package network

import (
	"CommonModule"
	"CommonModule/message"
	"CoreModule"
	"github.com/sirupsen/logrus"
	"time"
)

/*
本模块是网络模块，开启http/protobuf端口供外部调用
*/
type NetworkManager struct {
	core.MessageList                     // 消息列表
	httpService      *common.ServiceInfo //  http服务
}

// 初始化
func (n *NetworkManager) Init() {
	logrus.Infof("begin %s init", n.ModuleName)
	logrus.Infof("end %s init", n.ModuleName)
}

// 反初始化
func (n *NetworkManager) Uninit() {
	logrus.Infof("begin %s uninit", n.ModuleName)
	logrus.Infof("end %s uninit", n.ModuleName)
}

// 开始工作
func (n *NetworkManager) BeginWork() {
	logrus.Infof("begin %s beginwork", n.ModuleName)
	logrus.Infof("end %s beginwork", n.ModuleName)
}

// 停止工作
func (n *NetworkManager) StopWork() {
	logrus.Infof("begin %s stopwork", n.ModuleName)
	logrus.Infof("end %s stopwork", n.ModuleName)
}

// 偷窥消息
func (n *NetworkManager) OnForeseeMessage(msg message.BaseMessage) (done bool) {
	return
}

// 处理消息
func (n *NetworkManager) OnProcessMessage(msg message.BaseMessage) (rsp message.BaseResponse, err error) {
	switch msg.(type) {
	case *message.ConfigureMessage: // 配置消息
		return n.processConfigureMessage(msg)
	}
	return nil, nil
}

// 偷窥消息的回应
func (n *NetworkManager) OnForeseeResponse(rsp message.BaseResponse) (done bool) {
	return
}

// 处理消息的回应
func (n *NetworkManager) OnProcessResponse(rsp message.BaseResponse) {
	return
}

//  处理配置消息
func (n *NetworkManager) processConfigureMessage(msg message.BaseMessage) (rsp message.BaseResponse, err error) {
	conf := msg.(*message.ConfigureMessage)
	http := HttpService{Network: n}

	// 获取http配置信息
	info := conf.GgfService().ServiceInfo

	// 启动http服务
	go http.Startup(&info)
	return nil, nil
}

// 获取系统时间
func (n *NetworkManager) getTime() (time string, err error) {

	// 发送消息
	msg := message.NewSystemTimeMessage(common.Priority_First, message.Trans_Sync)
	rsp, err := n.SendMessage(msg)
	if err == nil {
		timeRsp := rsp.(*message.SystemTimeResponse)
		time = timeRsp.Time.Format("2006-01-02 15:04:05")
	}
	return
}

// 更新系统时间
func (n *NetworkManager) updateTime(time string) (err error) {
	msg := message.NewUpdateSystemTimeMessage(time, common.Priority_First, message.Trans_Sync)
	_, err = n.SendMessage(msg)
	return
}

// 获取ntp服务器信息
func (n *NetworkManager) getNtpServerInfo() (ntp common.NTPInfo, err error) {

	// 生成消息
	msg := message.NewNtpConfigureMessage(common.Priority_First, message.Trans_Sync)

	// 发送消息
	rsp, err := n.SendMessage(msg)

	// 解析回应
	if err == nil {
		ntpRsp := rsp.(*message.NtpConfigureResponse)
		ntp = ntpRsp.NTPInfo
	}
	return
}

// 更新ntp配置
func (n *NetworkManager) updateNtpConfigure(ntp common.NTPInfo, operate common.NtpControlType) (err error) {

	// 发送消息
	msg := message.NewUpdateNtpConfigureMessage(ntp, operate, common.Priority_First, message.Trans_Sync)
	_, err = n.SendMessage(msg)
	return
}

// 停止ggf服务
func (n *NetworkManager) stopGgfService() {
	stopMsg := message.NewStopGgfServiceMessage(common.Priority_Fifth, message.Trans_Async)
	n.SendMessage(stopMsg)
}

// 下载服务器的日志
func (n *NetworkManager) downloadServerLog() (logPath, zipLogPath string, err error) {
	// 发送消息
	msg := message.NewDownloadServerLogMessage(common.Priority_First, message.Trans_Sync)
	rsp, err := n.SendMessage(msg)
	if err != nil {
		return
	}

	// 处理回应
	logRsp := rsp.(*message.DownloadServerLogResponse)
	logPath = logRsp.LogPath
	zipLogPath = logRsp.ZipLogPath
	return
}

// 查看服务组件的日志
func (n *NetworkManager) getServiceLog(service string, begin, end int64) (beginIndex, endIndex int64, log []string, err error) {

	// 发送消息
	msg := message.NewGetServiceLogMessage(service, begin, end, common.Priority_First, message.Trans_Sync)
	rsp, err := n.SendMessage(msg)
	if err != nil {
		return
	}

	// 处理回应
	logRsp := rsp.(*message.GetServiceLogResponse)
	beginIndex = logRsp.Begin
	endIndex = logRsp.End
	log = logRsp.Log
	return
}

// 获取ggf service的信息
func (n *NetworkManager) getGgfServiceInfo() (host, system time.Time, gitBranch, gitCommit string, err error) {

	// 发送消息
	msg := message.NewGgfServiceMessage(common.Priority_First, message.Trans_Sync)
	rsp, err := n.SendMessage(msg)

	// 解析回应
	if err == nil {
		ggfRsp := rsp.(*message.HostServiceResponse)
		host = ggfRsp.HostStartup
		system = ggfRsp.SystemStartup
		gitBranch = ggfRsp.GitBranch
		gitCommit = ggfRsp.GitCommitID
	}

	return
}

// 记录事件
func (n *NetworkManager) recordEvent(eventType, event string) {

	// 生成消息
	msg := message.NewEventMessage(eventType, event, common.Priority_Fifth, message.Trans_Async)

	// 发送消息
	n.SendMessage(msg)
}

// 获取CPU的使用情况
func (n *NetworkManager) getCpuStatistic(begin, end string) (times []time.Time, statistic []common.CapacityInfo, err error) {

	// 发送消息
	msg := message.NewGetCpuStatisticMessage(begin, end, common.Priority_First, message.Trans_Sync)
	rsp, err := n.SendMessage(msg)
	if err != nil {
		return
	}

	// 解析回应
	statisticRsp := rsp.(*message.GetCpuStatisticResponse)
	times = statisticRsp.Time
	statistic = statisticRsp.Cpu
	return
}

// 获取磁盘的使用情况
func (n *NetworkManager) getDiskStatistic(begin, end string) (times []time.Time, statistic [][]common.DiskStatisticInfo, err error) {

	// 发送消息
	msg := message.NewGetDiskStatisticMessage(begin, end, common.Priority_First, message.Trans_Sync)
	rsp, err := n.SendMessage(msg)
	if err != nil {
		return
	}

	// 解析回应
	statisticRsp := rsp.(*message.GetDiskStatisticResponse)
	times = statisticRsp.Time
	statistic = statisticRsp.Statistic
	return
}

// 获取网络的使用情况
func (n *NetworkManager) getNetworkStatistic(begin, end string) (times []time.Time, statistic [][]common.NetworkStatisticInfo, err error) {

	// 发送消息
	msg := message.NewGetNetworkStatisticMessage(begin, end, common.Priority_First, message.Trans_Sync)
	rsp, err := n.SendMessage(msg)
	if err != nil {
		return
	}

	// 解析回应
	statisticRsp := rsp.(*message.GetNetworkStatisticResponse)
	times = statisticRsp.Time
	statistic = statisticRsp.Statistic
	return
}

// 获取内存的使用情况
func (n *NetworkManager) getMemoryStatistic(begin, end string) (times []time.Time, statistic []common.MemoryStatisticInfo, err error) {

	// 发送消息
	msg := message.NewGetMemoryStatisticMessage(begin, end, common.Priority_First, message.Trans_Sync)
	rsp, err := n.SendMessage(msg)
	if err != nil {
		return
	}

	// 解析回应
	statisticRsp := rsp.(*message.GetMemoryStatisticResponse)
	times = statisticRsp.Time
	statistic = statisticRsp.Statistic
	return
}

// 设置调试信息
func (n *NetworkManager) setDebugInfo(debug common.DebugInfo) (err error) {
	level, err := logrus.ParseLevel(debug.Log.LogLevel)
	if err != nil {
		return
	}

	logrus.SetLevel(level)
	return
}

// 获取调试信息
func (n *NetworkManager) getDebugInfo() (debug common.DebugInfo) {
	debug.Log.LogLevel = logrus.GetLevel().String()
	return
}

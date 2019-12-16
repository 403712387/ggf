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
	protobufService  *common.ServiceInfo // protobuf服务
	desKey           string              // des的密钥
}

// 初始化
func (n *NetworkManager) Init() {
	logrus.Infof("begin %s init", n.ModuleName)
	n.desKey = "*u9K_/M8"
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
	info := conf.HostService().ServiceInfo

	// 启动http服务
	go http.Startup(&info)
	return nil, nil
}

// 获取ip信息
func (n *NetworkManager) GetNetworkConfigure() (conf common.NetworkConfigure, err error) {
	// 生成消息
	msg := message.NewNetworkConfigureMessage(common.Priority_First, message.Trans_Sync)

	// 发送消息
	logrus.Tracef("send message:%s", msg.String())
	rsp, err := n.SendMessage(msg)
	if err == nil {
		cfgRsp := rsp.(*message.NetworkConfigureResponse)
		logrus.Tracef("recv response:%s", cfgRsp.String())
		return cfgRsp.GetNetworkConfigure(), err
	}
	return common.NetworkConfigure{}, err
}

// 更新IP信息
func (n *NetworkManager) UpdateNetworkConfigure(info common.NetworkInterface) (err error) {

	// 生成消息
	msg := message.NewUpdateNetworkConfigureMessage(info, common.Priority_Third, message.Trans_Sync)
	logrus.Tracef("send message:%s", msg.String())
	_, err = n.SendMessage(msg)
	return
}

// 升级服务
func (n *NetworkManager) ProcessUpdateService(info message.UpdateInfo, fileName string, fileData []byte) (err error) {

	// 发送升级消息
	logrus.Infof("update service info:%s, file name:%s, file length:%d", info.String(), fileName, len(fileData))
	msg := message.NewUpdateServiceMessage(info, fileName, fileData, common.Priority_Third, message.Trans_Sync)
	_, err = n.SendMessage(msg)
	return
}

// 登陆
func (n *NetworkManager) processLogin(user common.UserInfo) (err error) {

	// 发送消息
	msg := message.NewUserCheckMessage(user, common.Priority_First, message.Trans_Sync)
	_, err = n.SendMessage(msg)

	return
}

// 修改密码
func (n *NetworkManager) processChangePassword(pwd common.ChangePassword) (err error) {

	// 发送消息
	msg := message.NewChangePasswordMessage(pwd, common.Priority_First, message.Trans_Sync)
	_, err = n.SendMessage(msg)

	return
}

// 获取服务组件信息
func (n *NetworkManager) processServiceInfo() (info common.ServiceModules, err error) {

	// 发送消息
	msg := message.NewServiceInfoMessage(common.Priority_First, message.Trans_Sync)
	rsp, err := n.SendMessage(msg)
	if err != nil {
		return
	}

	// 处理回应
	serviceRsp := rsp.(*message.ServiceInfoResponse)
	info = serviceRsp.Info
	return
}

// 控制服务组件
func (n *NetworkManager) processControlService(control common.EntityControl) (err error) {

	// 发送消息
	msg := message.NewControlServiceMessage(control, common.Priority_First, message.Trans_Sync)
	_, err = n.SendMessage(msg)
	return
}

// 下载服务组件的日志
func (n *NetworkManager) processDownloadServiceLog(name string) (log string, err error) {

	// 生成消息
	msg := message.NewDownloadServiceLogMessage(name, common.Priority_First, message.Trans_Sync)
	rsp, err := n.SendMessage(msg)

	// 处理回应
	if err == nil {
		downloadRsp := rsp.(*message.DownloadServiceLogResponse)
		log = downloadRsp.Log
	}

	return
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

// 操作服务器（重启/关机)
func (n *NetworkManager) systemOperate(operate common.OperateType, time string) (err error) {

	// 发送消息
	msg := message.NewSystemControlMessage(operate, time, common.Priority_First, message.Trans_Sync)
	_, err = n.SendMessage(msg)
	return
}

// 获取服务器的标识符
func (n *NetworkManager) getServerId() (id, path, zipPath string, err error) {

	// 发送消息
	msg := message.NewServerIdMessage(common.Priority_First, message.Trans_Sync)
	rsp, err := n.SendMessage(msg)
	if err == nil {
		idRsp := rsp.(*message.ServerIdResponse)
		id = idRsp.ID
		path = idRsp.Path
		zipPath = idRsp.ZipPath
	}
	return
}

// 更新license信息
func (n *NetworkManager) updateLicense(license string) (err error) {

	// 发送消息
	msg := message.NewUpdateLicenseInfoMessage(license, common.Priority_First, message.Trans_Sync)
	_, err = n.SendMessage(msg)
	return
}

// 获取license信息
func (n *NetworkManager) getLicenseInfo() (license common.LicenseInfo, err error) {

	// 发送消息
	msg := message.NewLicenseInfoMessage(common.Priority_First, message.Trans_Sync)
	rsp, err := n.SendMessage(msg)
	if err == nil {
		licenseRsp := rsp.(*message.LicenseInfoResponse)
		license = licenseRsp.License
	}
	return
}

// 获取所有的时区信息
func (n *NetworkManager) getTimeZonesInfo() (info []common.TimeZone, err error) {

	// 生成消息
	msg := message.NewTimeZonesMessage(common.Priority_First, message.Trans_Sync)
	rsp, err := n.SendMessage(msg)
	if err == nil {
		timeZoneRsp := rsp.(*message.TimeZonesResponse)
		info = timeZoneRsp.TimeZones
	}
	return
}

// 获取当前主机的时区信息
func (n *NetworkManager) getTimeZone() (zone common.TimeZone, err error) {

	// 生成消息
	msg := message.NewServerTimeZonesMessage(common.Priority_First, message.Trans_Sync)
	rsp, err := n.SendMessage(msg)
	if err == nil {
		timeZoneRsp := rsp.(*message.ServerTimeZonesResponse)
		zone = timeZoneRsp.TimeZone
	}
	return
}

// 更新当前时区信息
func (n *NetworkManager) updateTimeZone(zone common.TimeZone) (err error) {

	// 生成消息
	msg := message.NewUpdateTimeZonesMessage(zone, common.Priority_First, message.Trans_Sync)
	_, err = n.SendMessage(msg)
	return
}

// 停止主机服务
func (n *NetworkManager) stopHostService() {
	stopMsg := message.NewStopHostServiceMessage(common.Priority_Fifth, message.Trans_Async)
	n.SendMessage(stopMsg)
}

// 查看服务器的日志
func (n *NetworkManager) getServerLog(begin, end int64) (beginIndex, endIndex int64, log []string, err error) {
	// 发送消息
	msg := message.NewGetServerLogMessage(begin, end, common.Priority_First, message.Trans_Sync)
	rsp, err := n.SendMessage(msg)
	if err != nil {
		return
	}

	// 处理回应
	logRsp := rsp.(*message.GetServerLogResponse)
	beginIndex = logRsp.Begin
	endIndex = logRsp.End
	log = logRsp.Log
	return
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

// 获取host service的信息
func (n *NetworkManager) getHostServiceInfo() (host, system time.Time, gitBranch, gitCommit string, err error) {

	// 发送消息
	msg := message.NewHostServiceMessage(common.Priority_First, message.Trans_Sync)
	rsp, err := n.SendMessage(msg)

	// 解析回应
	if err == nil {
		hostRsp := rsp.(*message.HostServiceResponse)
		host = hostRsp.HostStartup
		system = hostRsp.SystemStartup
		gitBranch = hostRsp.GitBranch
		gitCommit = hostRsp.GitCommitID
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

// 获取服务组件的资源使用情况
func (n *NetworkManager) getServiceStatistic(begin, end string) (times []time.Time, statistic [][]common.ProcessInfo, err error) {

	// 发送消息
	msg := message.NewGetServiceStatisticMessage(begin, end, common.Priority_First, message.Trans_Sync)
	rsp, err := n.SendMessage(msg)
	if err != nil {
		return
	}

	// 解析回应
	statisticRsp := rsp.(*message.GetServiceStatisticResponse)
	times = statisticRsp.Time
	statistic = statisticRsp.Statistic
	return
}

// 设置存储配置
func (n *NetworkManager) storageConfigure(storage common.StorageConfigure) {

	// 发送消息
	msg := message.NewStorageConfigureMessage(storage, common.Priority_First, message.Trans_Sync)
	n.SendMessage(msg)
}

// 获取存储配置
func (n *NetworkManager) getStorageConfigure() (removeThreshold int64, remove []string, err error) {

	// 发送消息
	msg := message.NewGetStorageConfigureMessage(common.Priority_First, message.Trans_Sync)
	rsp, err := n.SendMessage(msg)

	if err != nil {
		return
	}

	// 解析回应
	storageRsp := rsp.(*message.GetStorageConfigureResponse)
	removeThreshold = storageRsp.StorageConfigure.RemoveThreshold
	remove = storageRsp.LastRemove
	return
}

//查看磁盘的状态信息
func (n *NetworkManager) getDiskInfo() (partition []common.PartitionInfo, err error) {

	// 生成消息
	msg := message.NewPartitionInfoMessage(common.Priority_First, message.Trans_Sync)
	rsp, err := n.SendMessage(msg)

	// 解析回应
	if err == nil {
		diskRsp := rsp.(*message.DiskInfoResponse)
		partition = diskRsp.Info

	}

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

func (n *NetworkManager) getOperationDisk(disk common.OperationDisk) (err error) {

	// 生成消息
	msg := message.NewOperationDiskMessage(disk, common.Priority_First, message.Trans_Sync)
	_, err = n.SendMessage(msg)

	return
}

// des加密
func (n *NetworkManager) encryptDES(info string) (result string) {
	result = common.DESEncrypt([]byte(info), []byte(n.desKey))
	return
}

// des解密
func (n *NetworkManager) decryptDES(info string) (result string) {
	result = common.DESDecrypt(info, []byte(n.desKey))
	return
}

// 发送kafka消息
func (n *NetworkManager) sendKafkaMessage(topic, body string) {
	msg := message.NewKafkaMessage(topic, body, common.Priority_Third, message.Trans_Async)
	n.SendMessage(msg)
}

// 发送更新kafka配置的消息
func (n *NetworkManager) updateKafkaConfig(conf common.KafkaInfo) (err error) {

	// 生成消息
	msg := message.NewUpdateKafkaConfigMessage(conf, common.Priority_First, message.Trans_Sync)
	_, err = n.SendMessage(msg)

	return
}

// 升级服务
func (n *NetworkManager) ProcessUpdateMapMode(info message.MapInfo, fileName string, fileData []byte) (err error) {

	// 发送升级消息
	logrus.Infof("update service info:%s, file name:%s, file length:%d", info.String(), fileName, len(fileData))
	msg := message.NewUpdateMapModeMessage(info, fileName, fileData, common.Priority_Third, message.Trans_Sync)
	_, err = n.SendMessage(msg)
	return
}

//查看kafka的配置信息
func (n *NetworkManager) getKafkaConfigInfo() (info common.KafkaInfo, err error) {

	// 生成消息
	msg := message.NewGetKafkaConfigInfoMessage(common.Priority_First, message.Trans_Sync)
	rsp, err := n.SendMessage(msg)

	// 解析回应
	if err == nil {
		infoRsp := rsp.(*message.KafkaConfigInfoResponse)
		info = infoRsp.Info
	}
	return
}

//查看地图模式的配置信息
func (n *NetworkManager) getMapModeInfo() (info common.MapModeInfo, err error) {

	// 生成消息
	msg := message.NewGetMapModeInfoMessage(common.Priority_First, message.Trans_Sync)
	rsp, err := n.SendMessage(msg)

	// 解析回应
	if err == nil {
		infoRsp := rsp.(*message.MapModeInfoResponse)
		info = infoRsp.Info
	}
	return
}
package NtpModule

import (
	"CommonModule"
	"CommonModule/message"
	"CoreModule"
	"fmt"
	"github.com/sirupsen/logrus"
	"time"
)

/*
本模块是ntp模块，负责NTP校时
*/

type NtpManager struct {
	core.MessageList // 消息列表
	common.NTPInfo
	chanUpdate chan string // 更新
}

// 初始化
func (n *NtpManager) Init() {
	logrus.Infof("begin %s module uninit", n.ModuleName)
	n.chanUpdate = make(chan string)
	go n.ntpUpdateLoop()

	// 启动ntp服务器
	go StartNtpServer(123)
	logrus.Infof("end %s module uninit", n.ModuleName)
}

// 反初始化
func (n *NtpManager) Uninit() {
	logrus.Infof("begin %s module uninit", n.ModuleName)
	logrus.Infof("end %s module uninit", n.ModuleName)
}

// 开始工作
func (n *NtpManager) BeginWork() {
	logrus.Infof("begin %s module beginwork", n.ModuleName)
	logrus.Infof("end %s module beginwork", n.ModuleName)
}

// 停止工作
func (n *NtpManager) StopWork() {
	logrus.Infof("begin %s module stopwork", n.ModuleName)
	logrus.Infof("end %s module stopwork", n.ModuleName)
}

// 偷窥消息
func (n *NtpManager) OnForeseeMessage(msg message.BaseMessage) (done bool) {
	return
}

// 处理消息
func (n *NtpManager) OnProcessMessage(msg message.BaseMessage) (rsp message.BaseResponse, err error) {
	switch msg.(type) {
	case *message.UpdateNtpConfigureMessage: // 更新ntp
		return n.processUpdateNtpConfigureMessage(msg)
	case *message.NtpConfigureMessage: // 获取ntp
		return n.processNtpConfigureMessage(msg)
	}
	return nil, nil
}

// 偷窥消息的回应
func (n *NtpManager) OnForeseeResponse(rsp message.BaseResponse) (done bool) {
	return
}

// 处理消息的回应
func (n *NtpManager) OnProcessResponse(rsp message.BaseResponse) {
	return
}

// 获取ntp配置
func (n *NtpManager) processNtpConfigureMessage(msg message.BaseMessage) (rsp message.BaseResponse, err error) {
	rsp = message.NewNtpConfigureResponse(n.NTPInfo, msg)
	return
}

// 处理下发下来的ntp配置
func (n *NtpManager) processUpdateNtpConfigureMessage(msg message.BaseMessage) (rsp message.BaseResponse, err error) {
	updateMsg := msg.(*message.UpdateNtpConfigureMessage)
	rsp = message.NewBaseResponse(msg)

	switch updateMsg.Operate {
	case common.Ntp_Control_Set: // 设置ntp
		return n.setNtpServer(updateMsg)
	case common.Ntp_Control_Test: // 测试ntp服务
		return n.testNtpServer(updateMsg)
	default:
		err = fmt.Errorf("not support ntp control %s", updateMsg.Operate)
	}
	return
}

// 测试ntp服务器是否在线
func (n *NtpManager) testNtpServer(msg *message.UpdateNtpConfigureMessage) (rsp message.BaseResponse, err error) {

	// 获取ntp时间
	rsp = message.NewBaseResponse(msg)
	_, err = ntpTime(msg.IP, msg.Port)

	// 判断是否获取成功
	if err != nil {
		logrus.Errorf("get ntp time fail, ip:%s, port:%d, error:%s", msg.IP, msg.Port, err.Error())
	}
	return
}

// 处理ntp
func (n *NtpManager) setNtpServer(msg *message.UpdateNtpConfigureMessage) (rsp message.BaseResponse, err error) {
	rsp = message.NewBaseResponse(msg)

	if msg.Enable && msg.ProofreadInterval <= 0 {
		err = fmt.Errorf("%s", "invalid proofread interval")
		return
	}

	// 更新ntp信息
	n.NTPInfo = msg.NTPInfo
	logrus.Infof("update ntp server info:%s", n.NTPInfo.String())

	// 通知进行校时
	n.chanUpdate <- "hello,ring"
	return
}

func (n *NtpManager) ntpUpdateLoop() {
	logrus.Infof("begin ntp update loop")
	for {
		interval, _ := time.ParseDuration(fmt.Sprintf("%dm", common.MaxInt64(n.ProofreadInterval, 5)))
		select {
		case <-time.After(interval):
		case name := <-n.chanUpdate:
			logrus.Infof("%s", name)
		}

		// ntp不可用
		if !n.NTPInfo.Enable {
			continue
		}

		// 获取ntp时间
		now, err := ntpTime(n.IP, n.Port)
		if err != nil {
			logrus.Errorf("ntp proofread fail, ntp info:%s, error:%s", n.NTPInfo.String(), err.Error())
		} else {

			// 判断时间是否合法
			newTime := now.Format("2006-01-02 15:04:05")
			_, err = time.Parse("2006-01-02 15:04:05", newTime)
			if err != nil {
				logrus.Errorf("get time from ntp, but time not valid, time:%s, error:%s", newTime, err.Error())
				continue
			}

			// 更新时间
			common.CommondResult(fmt.Sprintf(`date -s "%s"`, newTime))
			common.CommondResult("hwclock -w")
			logrus.Infof("ntp proofread successful, now:%s", newTime)
		}
	}
	logrus.Infof("end ntp update loop")
}

package HostServiceModule

import (
	"CommonModule/message"
	"CoreModule"
	"github.com/shirou/gopsutil/host"
	"github.com/sirupsen/logrus"
	"time"
)

// git信息
var (
	GitBranch   string // git版本
	GitCommitID string // git commit id
)

/*
本模块主要负责ggf服务的信息，比如启动时间，git信息等
*/

type HostServiceManager struct {
	core.MessageList             // 消息列表
	HostServiceStartup time.Time // 	主机服务的启动时间
	SystemStartup      time.Time // 系统的启动时间
}

// 初始化
func (h *HostServiceManager) Init() {
	logrus.Infof("begin %s module uninit", h.ModuleName)
	h.HostServiceStartup = time.Now()

	up, _ := host.BootTime()
	h.SystemStartup = time.Unix(int64(up), 0)
	logrus.Infof("end %s module uninit", h.ModuleName)
}

// 反初始化
func (h *HostServiceManager) Uninit() {
	logrus.Infof("begin %s module uninit", h.ModuleName)
	logrus.Infof("end %s module uninit", h.ModuleName)
}

// 开始工作
func (h *HostServiceManager) BeginWork() {
	logrus.Infof("begin %s module beginwork", h.ModuleName)
	logrus.Infof("end %s module beginwork", h.ModuleName)
}

// 停止工作
func (h *HostServiceManager) StopWork() {
	logrus.Infof("begin %s module stopwork", h.ModuleName)
	logrus.Infof("end %s module stopwork", h.ModuleName)
}

// 偷窥消息
func (h *HostServiceManager) OnForeseeMessage(msg message.BaseMessage) (done bool) {
	return
}

// 处理消息
func (h *HostServiceManager) OnProcessMessage(msg message.BaseMessage) (rsp message.BaseResponse, err error) {
	switch msg.(type) {
	case *message.HostServiceMessage:
		return h.processHostServiceMessage(msg)
	}
	return nil, nil
}

// 偷窥消息的回应
func (h *HostServiceManager) OnForeseeResponse(rsp message.BaseResponse) (done bool) {
	return
}

// 处理消息的回应
func (h *HostServiceManager) OnProcessResponse(rsp message.BaseResponse) {
	return
}

// 获取主机服务的信息
func (h *HostServiceManager) processHostServiceMessage(msg message.BaseMessage) (rsp message.BaseResponse, err error) {
	rsp = message.NewHostServiceResponse(h.HostServiceStartup, h.SystemStartup, GitBranch, GitCommitID, msg)
	return
}

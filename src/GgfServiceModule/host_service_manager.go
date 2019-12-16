package GgfServiceModule

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

type GgfServiceManager struct {
	core.MessageList            // 消息列表
	GgfServiceStartup time.Time // 	主机服务的启动时间
	SystemStartup     time.Time // 系统的启动时间
}

// 初始化
func (h *GgfServiceManager) Init() {
	logrus.Infof("begin %s module uninit", h.ModuleName)
	h.GgfServiceStartup = time.Now()

	up, _ := host.BootTime()
	h.SystemStartup = time.Unix(int64(up), 0)
	logrus.Infof("end %s module uninit", h.ModuleName)
}

// 反初始化
func (h *GgfServiceManager) Uninit() {
	logrus.Infof("begin %s module uninit", h.ModuleName)
	logrus.Infof("end %s module uninit", h.ModuleName)
}

// 开始工作
func (h *GgfServiceManager) BeginWork() {
	logrus.Infof("begin %s module beginwork", h.ModuleName)
	logrus.Infof("end %s module beginwork", h.ModuleName)
}

// 停止工作
func (h *GgfServiceManager) StopWork() {
	logrus.Infof("begin %s module stopwork", h.ModuleName)
	logrus.Infof("end %s module stopwork", h.ModuleName)
}

// 偷窥消息
func (h *GgfServiceManager) OnForeseeMessage(msg message.BaseMessage) (done bool) {
	return
}

// 处理消息
func (h *GgfServiceManager) OnProcessMessage(msg message.BaseMessage) (rsp message.BaseResponse, err error) {
	switch msg.(type) {
	case *message.GgfServiceMessage:
		return h.processGgfServiceMessage(msg)
	}
	return nil, nil
}

// 偷窥消息的回应
func (h *GgfServiceManager) OnForeseeResponse(rsp message.BaseResponse) (done bool) {
	return
}

// 处理消息的回应
func (h *GgfServiceManager) OnProcessResponse(rsp message.BaseResponse) {
	return
}

// 获取ggf服务的信息
func (h *GgfServiceManager) processGgfServiceMessage(msg message.BaseMessage) (rsp message.BaseResponse, err error) {
	rsp = message.NewGgfServiceResponse(h.GgfServiceStartup, h.SystemStartup, GitBranch, GitCommitID, msg)
	return
}

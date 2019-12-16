package configure

import (
	"CommonModule"
	"CommonModule/message"
	"CoreModule"
	"encoding/json"
	"fmt"
	"github.com/sirupsen/logrus"
	"io/ioutil"
)

/*
本模块是配置管理模块，从配置文件(./config/config.json)中读取配置，并将配置发送出去，如果配置有更新，则同步到配置文件中
*/
// 配置信息
type ConfigureManager struct {
	core.MessageList               // 消息列表
	ConfigureFile    string        // 配置文件名字
	configure        ConfigureJson // 配置内容
}

// json配置
type ConfigureJson struct {
	HostService GgfServiceJson `json:"ggf_service"` // host name
	NtpServer   common.NTPInfo `json:"ntp_server"`
}

func (c *ConfigureJson) String() string {
	return fmt.Sprintf("host service info:%s", c.HostService.String())
}

type GgfServiceJson struct {
	HttpPort int32 `json:"http_port"` // http port
}

func (h *GgfServiceJson) String() string {
	return fmt.Sprintf("http port:%d", h.HttpPort)
}

// 初始化
func (c *ConfigureManager) Init() {
	logrus.Infof("begin %s init", c.ModuleName)
	logrus.Infof("end %s init", c.ModuleName)
}

// 反初始化
func (c *ConfigureManager) Uninit() {
	logrus.Infof("begin %s uninit", c.ModuleName)
	logrus.Infof("end %s uninit", c.ModuleName)
}

// 开始工作
func (c *ConfigureManager) BeginWork() {
	logrus.Infof("begin %s beginwork", c.ModuleName)

	// 读取配置文件
	configureData, err := ioutil.ReadFile(c.ConfigureFile)
	if err != nil {
		logrus.Fatalf("open configure file fail, configure file:%s, error reason:%s\n", c.ConfigureFile, err)
		return
	}

	// 解析json
	err = json.Unmarshal(configureData, &c.configure)
	if err != nil {
		logrus.Fatalf("parse json info fail, configure info:%s, error reason:%s\n", string(configureData), err)
		return
	}

	logrus.Infof("get configure successful, info:%s", c.configure.String())

	// 发送配置消息
	hostServiceInfo := common.HostServiceInfo{ServiceInfo: common.ServiceInfo{HttpPort: c.configure.HostService.HttpPort}}
	configMessage := message.NewConfigureMessage(hostServiceInfo, common.Priority_First, message.Trans_Sync)
	c.SendMessage(configMessage)

	// 发送ntp消息
	ntpMsg := message.NewUpdateNtpConfigureMessage(c.configure.NtpServer, common.Ntp_Control_Set, common.Priority_Third, message.Trans_Async)
	c.SendMessage(ntpMsg)

	logrus.Infof("end %s beginwork", c.ModuleName)
}

// 停止工作
func (c *ConfigureManager) StopWork() {
	logrus.Infof("begin %s stopwork", c.ModuleName)
	logrus.Infof("end %s stopwork", c.ModuleName)
}

// 偷窥消息
func (c *ConfigureManager) OnForeseeMessage(msg message.BaseMessage) (done bool) {
	switch msg.(type) {
	case *message.UpdateNtpConfigureMessage: // ntp配置有更新
		c.foreseeUpdateNtpConfigureMessage(msg)
	}
	return
}

// 处理消息
func (c *ConfigureManager) OnProcessMessage(msg message.BaseMessage) (rsp message.BaseResponse, err error) {
	switch msg.(type) {
	}
	return nil, nil
}

// 偷窥消息的回应
func (c *ConfigureManager) OnForeseeResponse(rsp message.BaseResponse) (done bool) {
	return
}

// 处理消息的回应
func (c *ConfigureManager) OnProcessResponse(rsp message.BaseResponse) {
	return
}

// 处理ntp配置更新
func (c *ConfigureManager) foreseeUpdateNtpConfigureMessage(msg message.BaseMessage) {
	ntpMsg := msg.(*message.UpdateNtpConfigureMessage)

	// 如果配置有更新，个写入配置文件中
	if ntpMsg.NTPInfo != c.configure.NtpServer {
		c.configure.NtpServer = ntpMsg.NTPInfo
		c.saveConfigure()
	}
}

//  配置写入json文件
func (c *ConfigureManager) saveConfigure() error {
	data, err := json.MarshalIndent(c.configure, "", " ")
	ioutil.WriteFile(c.ConfigureFile, data, 0666)
	return err
}

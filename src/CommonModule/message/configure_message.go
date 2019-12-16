package message

import (
	"CommonModule"
	"fmt"
	"time"
)

// 配置消息
type ConfigureMessage struct {
	BaseMessageInfo
	hostService common.GgfServiceInfo //  ggf服务的信息
}

// 生成配置消息
func NewConfigureMessage(host common.GgfServiceInfo, pri common.Priority, tra TransType) *ConfigureMessage {
	MessageId++
	return &ConfigureMessage{BaseMessageInfo: BaseMessageInfo{id: MessageId, priority: pri, trans: tra, birthday: time.Now()}, hostService: host}
}

// host service信息
func (c *ConfigureMessage) GgfService() common.GgfServiceInfo {
	return c.hostService
}

// 消息转为string
func (c *ConfigureMessage) String() string {
	result := fmt.Sprintf("%s, %s", c.hostService.String(), c.BaseMessageInfo.String())
	return result
}

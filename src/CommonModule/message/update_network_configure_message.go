package message

import (
	"CommonModule"
	"fmt"
	"time"
)

// 更新IP信息
type UpdateNetworkConfigureMessage struct {
	BaseMessageInfo
	ipInfo 			common.NetworkInterface		// ip信息
}

// 生成更新IP的消息
func NewUpdateNetworkConfigureMessage(ip common.NetworkInterface, pri common.Priority, trs TransType) *UpdateNetworkConfigureMessage {
	MessageId++
	return &UpdateNetworkConfigureMessage{BaseMessageInfo: BaseMessageInfo{id: MessageId, priority:pri, trans:trs, birthday:time.Now()}, ipInfo:ip}
}

// 获取ip配置
func (m *UpdateNetworkConfigureMessage)GetIpInfo() common.NetworkInterface {
	return m.ipInfo
}

func (m *UpdateNetworkConfigureMessage)String() string {
	return fmt.Sprintf("network configure info:%s, %s", m.ipInfo.String(), m.BaseMessageInfo.String())
}
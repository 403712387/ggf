package message

import (
	"CommonModule"
	"fmt"
	"time"
)

// 获取IP信息
type NetworkConfigureMessage struct {
	BaseMessageInfo
}

// 生成消息
func NewNetworkConfigureMessage(pri common.Priority, tra TransType) (msg *NetworkConfigureMessage) {
	MessageId++
	return &NetworkConfigureMessage{BaseMessageInfo: BaseMessageInfo{id: MessageId, priority: pri, trans: tra, birthday: time.Now()}}
}

func (g *NetworkConfigureMessage) String() string {
	return g.BaseMessageInfo.String()
}

// 获取网络配置的响应
type NetworkConfigureResponse struct {
	BaseResponseInfo
	networkConfigure common.NetworkConfigure // ip信息
}

// 生成回应
func NewNetworkConfigureResponse(config common.NetworkConfigure, msg BaseMessage) (response *NetworkConfigureResponse) {
	resp := NetworkConfigureResponse{BaseResponseInfo: BaseResponseInfo{message: msg}, networkConfigure: config}
	return &resp
}

// 获取网络配置
func (g *NetworkConfigureResponse) GetNetworkConfigure() (config common.NetworkConfigure) {
	return g.networkConfigure
}

func (g *NetworkConfigureResponse) String() (result string) {
	return fmt.Sprintf("%s, %s", g.networkConfigure.String(), g.BaseResponseInfo.String())
}

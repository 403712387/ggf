package message

import (
	"CommonModule"
	"fmt"
	"time"
)

// 获取ntp配置
type NtpConfigureMessage struct {
	BaseMessageInfo
}

// 生成消息
func NewNtpConfigureMessage(pri common.Priority, tra TransType) (msg *NtpConfigureMessage) {
	MessageId++
	return &NtpConfigureMessage{BaseMessageInfo: BaseMessageInfo{id: MessageId, priority: pri, trans: tra, birthday: time.Now()}}
}

func (n *NtpConfigureMessage) String() string {
	return fmt.Sprintf("%s", n.BaseMessageInfo.String())
}

// 获取ntp配置的响应
type NtpConfigureResponse struct {
	BaseResponseInfo
	common.NTPInfo
}

// 生成回应
func NewNtpConfigureResponse(ntp common.NTPInfo, msg BaseMessage) (response *NtpConfigureResponse) {
	resp := NtpConfigureResponse{BaseResponseInfo: BaseResponseInfo{message: msg}, NTPInfo: ntp}
	return &resp
}

func (n *NtpConfigureResponse) String() (result string) {
	result = fmt.Sprintf("%s, %s", n.NTPInfo.String(), n.BaseResponseInfo.String())
	return
}

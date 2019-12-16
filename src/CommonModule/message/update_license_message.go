package message

import (
	"CommonModule"
	"fmt"
	"time"
)

// 获取服务器标识符的消息
type UpdateLicenseInfoMessage struct {
	BaseMessageInfo
	License string // 服务器标识
}

// 生成消息
func NewUpdateLicenseInfoMessage(license string, pri common.Priority, trs TransType) *UpdateLicenseInfoMessage {
	MessageId++
	return &UpdateLicenseInfoMessage{BaseMessageInfo: BaseMessageInfo{id: MessageId, priority: pri, trans: trs, birthday: time.Now()}, License: license}
}

func (s *UpdateLicenseInfoMessage) String() string {
	return fmt.Sprintf("license:%s, %s", s.License, s.BaseMessageInfo.String())
}

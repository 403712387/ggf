package message

import (
	"CommonModule"
	"fmt"
	"time"
)

// 更新ntp配置
type UpdateNtpConfigureMessage struct {
	BaseMessageInfo
	common.NTPInfo                       // NTP配置
	Operate        common.NtpControlType // 测试还是设置
}

// 生成消息
func NewUpdateNtpConfigureMessage(ntp common.NTPInfo, operate common.NtpControlType, pri common.Priority, tra TransType) (msg *UpdateNtpConfigureMessage) {
	MessageId++
	return &UpdateNtpConfigureMessage{BaseMessageInfo: BaseMessageInfo{id: MessageId, priority: pri, trans: tra, birthday: time.Now()}, NTPInfo: ntp, Operate: operate}
}

func (n *UpdateNtpConfigureMessage) String() string {
	return fmt.Sprintf("%s, operate:%s,  %s", n.NTPInfo.String(), n.Operate, n.BaseMessageInfo.String())
}

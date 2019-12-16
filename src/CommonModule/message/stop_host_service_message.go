package message

import (
	"CommonModule"
	"time"
)

// 停止主机服务
type StopHostServiceMessage struct {
	BaseMessageInfo
}

// 生成消息
func NewStopGgfServiceMessage(pri common.Priority, tra TransType) (msg *StopHostServiceMessage) {
	MessageId++
	return &StopHostServiceMessage{BaseMessageInfo: BaseMessageInfo{id: MessageId, priority: pri, trans: tra, birthday: time.Now()}}
}

func (s *StopHostServiceMessage) String() string {
	return s.BaseMessageInfo.String()
}

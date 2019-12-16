package message

import (
	"CommonModule"
	"fmt"
	"time"
)

// 更新系统时间
type UpdateSystemTimeMessage struct {
	BaseMessageInfo
	Time string // 设置的系统时间
}

// 生成消息
func NewUpdateSystemTimeMessage(systemTime string, pri common.Priority, tra TransType) (msg *UpdateSystemTimeMessage) {
	MessageId++
	return &UpdateSystemTimeMessage{BaseMessageInfo: BaseMessageInfo{id: MessageId, priority: pri, trans: tra, birthday: time.Now()}, Time: systemTime}
}

func (s *UpdateSystemTimeMessage) String() string {
	return fmt.Sprintf("time:%s, %s", s.Time, s.BaseMessageInfo.String())
}

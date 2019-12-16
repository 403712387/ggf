package message

import (
	"CommonModule"
	"fmt"
	"time"
)

// 记录事件的消息
type EventMessage struct {
	BaseMessageInfo
	Type    string // 事件类型
	Explain string // 事件的描述
}

// 生成消息
func NewEventMessage(eventType, explain string, pri common.Priority, tra TransType) (msg *EventMessage) {
	MessageId++
	return &EventMessage{BaseMessageInfo: BaseMessageInfo{id: MessageId, priority: pri, trans: tra, birthday: time.Now()}, Type: eventType, Explain: explain}
}

func (d *EventMessage) String() string {
	return fmt.Sprintf("type:%s, explain:%s, %s", d.Type, d.Explain, d.BaseMessageInfo.String())
}

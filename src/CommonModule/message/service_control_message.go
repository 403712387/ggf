package message

import (
	"CommonModule"
	"fmt"
	"time"
)

// 控制服务的重启/启动/停止
type ServiceControlMessage struct {
	BaseMessageInfo
	common.EntityControl
}

// 生成消息
func NewControlServiceMessage(control common.EntityControl, pri common.Priority, tra TransType) (msg *ServiceControlMessage) {
	MessageId++
	return &ServiceControlMessage{BaseMessageInfo: BaseMessageInfo{id: MessageId, priority: pri, trans: tra, birthday: time.Now()}, EntityControl: control}
}

func (c *ServiceControlMessage) String() string {
	return fmt.Sprintf("control type:%s, time:%s, %s", c.Control, c.Time, c.BaseMessageInfo.String())
}

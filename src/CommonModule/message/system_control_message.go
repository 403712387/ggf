package message

import (
	"CommonModule"
	"fmt"
	"time"
)

// 系统控制，用来控制系统的停止/重启
type SystemControlMessage struct {
	BaseMessageInfo
	common.OperateType
	Time string // 是否立即执行
}

// 生成消息
func NewSystemControlMessage(operate common.OperateType, operateTime string, pri common.Priority, tra TransType) (msg *SystemControlMessage) {
	MessageId++
	return &SystemControlMessage{BaseMessageInfo: BaseMessageInfo{id: MessageId, priority: pri, trans: tra, birthday: time.Now()}, OperateType: operate, Time: operateTime}
}

func (d *SystemControlMessage) String() string {
	return fmt.Sprintf("operate:%s, %s", d.OperateType, d.BaseMessageInfo.String())
}

package message

import (
	"CommonModule"
	"time"
)

type Times struct {
	Time []string `json:"time"`
}

// 删除存储的消息
type RemoveStorageMessage struct {
	BaseMessageInfo
	Times
}

// 生成消息
func NewRemoveStorageMessage(times []string, pri common.Priority, trs TransType) *RemoveStorageMessage {
	MessageId++
	return &RemoveStorageMessage{BaseMessageInfo: BaseMessageInfo{id: MessageId, priority: pri, trans: trs, birthday: time.Now()}, Times: Times{Time: times}}
}

// 获取删除的时间
func (r *RemoveStorageMessage) GetRemoveTime() Times {
	return r.Times
}

func (r *RemoveStorageMessage) String() string {
	var result string
	for _, date := range r.Time {
		result += date + ", "
	}
	result += r.BaseMessageInfo.String()
	return result
}


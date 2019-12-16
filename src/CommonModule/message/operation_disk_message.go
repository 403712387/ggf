package message

import (
	"CommonModule"
	"time"
)

type OperationDiskMessage struct {
	BaseMessageInfo
	Partitions common.OperationDisk
}

//生成需要挂载磁盘的信息
func NewOperationDiskMessage(partition common.OperationDisk, pri common.Priority, tra TransType) (msg *OperationDiskMessage) {
	MessageId++
	return &OperationDiskMessage{BaseMessageInfo: BaseMessageInfo{id: MessageId, priority: pri, trans: tra, birthday: time.Now()}, Partitions: partition}
}

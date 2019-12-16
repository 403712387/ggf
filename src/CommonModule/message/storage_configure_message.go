package message

import (
	"CommonModule"
	"fmt"
	"time"
)

// 更新存储配置
type StorageConfigureMessage struct {
	BaseMessageInfo
	StorageConfigure common.StorageConfigure // 存储的配置
}

// 生成新的消息
func NewStorageConfigureMessage(conf common.StorageConfigure, pri common.Priority, trs TransType) *StorageConfigureMessage {
	MessageId++
	return &StorageConfigureMessage{BaseMessageInfo: BaseMessageInfo{id: MessageId, priority: pri, trans: trs, birthday: time.Now()}, StorageConfigure: conf}
}

func (s *StorageConfigureMessage) String() string {
	return fmt.Sprintf("%s, %s", s.StorageConfigure.String(), s.BaseMessageInfo.String())
}

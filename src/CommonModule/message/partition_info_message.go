package message

import (
	"CommonModule"
	"time"
)

// 磁盘的分区信息
type PartitionInfoMessage struct {
	BaseMessageInfo
}

// 生成消息
func NewPartitionInfoMessage(pri common.Priority, tra TransType) (msg *PartitionInfoMessage) {
	MessageId++
	return &PartitionInfoMessage{BaseMessageInfo: BaseMessageInfo{id: MessageId, priority: pri, trans: tra, birthday: time.Now()}}
}

func (p *PartitionInfoMessage) String() (result string) {

	result = p.BaseMessageInfo.String()
	return
}

type DiskInfoResponse struct {
	BaseResponseInfo
	Info []common.PartitionInfo // 磁盘的信息
}

// 生成回应
func NewDiskInfoResponse(info []common.PartitionInfo, msg BaseMessage) (response *DiskInfoResponse) {
	resp := DiskInfoResponse{BaseResponseInfo: BaseResponseInfo{message: msg}, Info: info}
	return &resp
}

func (s *DiskInfoResponse) String() (result string) {

	for _, info := range s.Info {
		result += info.String() + ", "
	}
	result += s.BaseResponseInfo.String()
	return
}






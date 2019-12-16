package message

import (
	"CommonModule"
	"time"
)

type UpdateKafkaConfigMessage struct {
	BaseMessageInfo
	Info common.KafkaInfo
}

//生成需要修改kafka的信息
func NewUpdateKafkaConfigMessage(kafkainfo common.KafkaInfo, pri common.Priority, tra TransType) (msg *UpdateKafkaConfigMessage) {
	MessageId++
	return &UpdateKafkaConfigMessage{BaseMessageInfo: BaseMessageInfo{id: MessageId, priority: pri, trans: tra, birthday: time.Now()}, Info: kafkainfo}
}

// kafka的配置信息
type GetKafkaConfigInfoMessage struct {
	BaseMessageInfo
}

// 生成消息
func NewGetKafkaConfigInfoMessage(pri common.Priority, tra TransType) (msg *GetKafkaConfigInfoMessage) {
	MessageId++
	return &GetKafkaConfigInfoMessage{BaseMessageInfo: BaseMessageInfo{id: MessageId, priority: pri, trans: tra, birthday: time.Now()}}
}

func (p *GetKafkaConfigInfoMessage) String() (result string) {

	result = p.BaseMessageInfo.String()
	return
}

type KafkaConfigInfoResponse struct {
	BaseResponseInfo
	Info common.KafkaInfo // kafka的配置信息
}

// 生成回应
func NewGetKafkaConfigInfoResponse(info common.KafkaInfo, msg BaseMessage) (response *KafkaConfigInfoResponse) {
	resp := KafkaConfigInfoResponse{BaseResponseInfo: BaseResponseInfo{message: msg}, Info: info}
	return &resp
}

func (s *KafkaConfigInfoResponse) String() (result string) {

	result = s.BaseResponseInfo.String()
	return
}
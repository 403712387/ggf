package message

import (
	"CommonModule"
	"time"
)
// 发送kafka消息
type KafkaMessage struct {
	BaseMessageInfo
	Topic 		string 		// kafka的topic
	Body 		string 		// kafka的body
}

// 生成消息
func NewKafkaMessage(topic, body string, pri common.Priority, tra TransType) (msg *KafkaMessage) {
	MessageId++
	return &KafkaMessage{BaseMessageInfo: BaseMessageInfo{id: MessageId, priority:pri, trans:tra, birthday:time.Now()}, Topic:topic, Body:body}
}

func (k *KafkaMessage)String() string {
	return "topic:" + k.Topic + ", body:" + k.Body + k.BaseMessageInfo.String()
}

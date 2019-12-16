package message

import (
	"CommonModule"
	"fmt"
	"time"
)

// 消息发送类型（同步消息还是异步消息
type TransType int32

const (
	Trans_Async TransType = iota // 异步发送
	Trans_Sync                   // 同步发送
)

func (trans TransType) String() string {
	result := "trans async"
	switch trans {
	case Trans_Async:
		result = "trans async"
	case Trans_Sync:
		result = "trans sync"
	}

	return result
}

// message的接口
type BaseMessage interface {
	Id() int64                        // 消息ID
	MessagePriority() common.Priority // 消息优先级
	TransType() TransType             // 发送类型，同步消息还是异步消息
	Birthday() time.Time              // 消息的生成时间
	String() string
}

/*
基础消息类
*/
var MessageId int64 // 消息的ID
type BaseMessageInfo struct {
	id       int64           // 消息的ID
	priority common.Priority // 消息优先级
	trans    TransType       // 消息传送类型
	birthday time.Time       // 消息的产生时间
}

// 消息的id
func (m *BaseMessageInfo) Id() int64 {
	return m.id
}

// 消息的优先级
func (m *BaseMessageInfo) MessagePriority() common.Priority {
	return m.priority
}

// 消息的传输类型
func (m *BaseMessageInfo) TransType() TransType {
	return m.trans
}

// 消息的生成时间
func (m *BaseMessageInfo) Birthday() time.Time {
	return m.birthday
}

// 基类消息转为string
func (m *BaseMessageInfo) String() string {
	result := fmt.Sprintf("id:%d, priority:%s, trans type:%s, birthday:%s", m.id, m.priority.String(), m.trans.String(), m.birthday.Format("2006-01-02 15:04:05.000"))
	return result
}

// response的接口
type BaseResponse interface {
	Message() BaseMessage // 消息
	String() string
}

// 消息的回应
type BaseResponseInfo struct {
	message BaseMessage // 消息
}

// 生成回应
func NewBaseResponse(msg BaseMessage) *BaseResponseInfo {
	return &BaseResponseInfo{message: msg}
}

// 消息
func (r *BaseResponseInfo) Message() BaseMessage {
	return r.message
}

// 基类消息的回应转为string
func (r *BaseResponseInfo) String() string {
	result := fmt.Sprintf("message:%s", r.message.String())
	return result
}

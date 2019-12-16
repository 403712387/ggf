package message

import (
	"CommonModule"
	"fmt"
	"time"
)

// 获取系统时间
type SystemTimeMessage struct {
	BaseMessageInfo
}

// 生成消息
func NewSystemTimeMessage(pri common.Priority, tra TransType) (msg *SystemTimeMessage) {
	MessageId++
	return &SystemTimeMessage{BaseMessageInfo: BaseMessageInfo{id: MessageId, priority: pri, trans: tra, birthday: time.Now()}}
}

func (s *SystemTimeMessage) String() string {
	return s.BaseMessageInfo.String()
}

// 获取系统时间的回应
type SystemTimeResponse struct {
	BaseResponseInfo
	Time time.Time
}

// 生成回应
func NewSystemTimeResponse(current time.Time, msg BaseMessage) (response *SystemTimeResponse) {
	resp := SystemTimeResponse{BaseResponseInfo: BaseResponseInfo{message: msg}, Time: current}
	return &resp
}

func (s *SystemTimeResponse) String() (result string) {
	return fmt.Sprintf("system time:%s, %s", s.Time.Format("2006-01-02 15:04:05.000"), s.BaseResponseInfo.String())
}

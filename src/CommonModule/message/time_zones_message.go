package message

import (
	"CommonModule"
	"fmt"
	"time"
)

// 获取所有时区信息
type TimeZonesMessage struct {
	BaseMessageInfo
}

// 生成消息
func NewTimeZonesMessage(pri common.Priority, tra TransType) (msg *TimeZonesMessage) {
	MessageId++
	return &TimeZonesMessage{BaseMessageInfo: BaseMessageInfo{id: MessageId, priority: pri, trans: tra, birthday: time.Now()}}
}

func (t *TimeZonesMessage) String() string {
	return t.BaseMessageInfo.String()
}

// 获取所有时区信息的回应
type TimeZonesResponse struct {
	BaseResponseInfo
	TimeZones []common.TimeZone
}

// 生成回应
func NewTimeZonesResponse(zones []common.TimeZone, msg BaseMessage) (response *TimeZonesResponse) {
	resp := TimeZonesResponse{BaseResponseInfo: BaseResponseInfo{message: msg}, TimeZones: zones}
	return &resp
}

func (t *TimeZonesResponse) String() (result string) {
	return fmt.Sprintf("time zones:%s, %s", t.TimeZones, t.BaseResponseInfo.String())
}

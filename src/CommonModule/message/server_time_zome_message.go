package message

import (
	"CommonModule"
	"fmt"
	"time"
)

// 获取时区信息
type ServerTimeZoneMessage struct {
	BaseMessageInfo
}

// 生成消息
func NewServerTimeZonesMessage(pri common.Priority, tra TransType) (msg *ServerTimeZoneMessage) {
	MessageId++
	return &ServerTimeZoneMessage{BaseMessageInfo: BaseMessageInfo{id: MessageId, priority: pri, trans: tra, birthday: time.Now()}}
}

func (s *ServerTimeZoneMessage) String() string {
	return s.BaseMessageInfo.String()
}

// 获取时区信息的回应
type ServerTimeZonesResponse struct {
	BaseResponseInfo
	TimeZone common.TimeZone
}

// 生成回应
func NewServerTimeZonesResponse(zone common.TimeZone, msg BaseMessage) (response *ServerTimeZonesResponse) {
	resp := ServerTimeZonesResponse{BaseResponseInfo: BaseResponseInfo{message: msg}, TimeZone: zone}
	return &resp
}

func (s *ServerTimeZonesResponse) String() (result string) {
	return fmt.Sprintf("time zone:%s, %s", s.TimeZone.String(), s.BaseResponseInfo.String())
}

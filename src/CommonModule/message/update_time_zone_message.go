package message

import (
	"CommonModule"
	"fmt"
	"time"
)

// 更新时区信息
type UpdateTimeZoneMessage struct {
	BaseMessageInfo
	TimeZone 	common.TimeZone		// 时区信息
}

// 生成消息
func NewUpdateTimeZonesMessage(zone common.TimeZone, pri common.Priority, tra TransType) (msg *UpdateTimeZoneMessage) {
	MessageId++
	return &UpdateTimeZoneMessage{BaseMessageInfo: BaseMessageInfo{id: MessageId, priority: pri, trans: tra, birthday: time.Now()}, TimeZone:zone}
}

func (u *UpdateTimeZoneMessage) String() string {
	return fmt.Sprintf("time zone:%s, %s", u.TimeZone.String(), u.BaseMessageInfo.String())
}

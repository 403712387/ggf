package message

import (
	"CommonModule"
	"time"
)

// 检验用户是否合法的消息
type UserVerifyMessage struct {
	BaseMessageInfo
	common.UserInfo
}

// 生成消息
func NewUserCheckMessage(user common.UserInfo, pri common.Priority, tra TransType) (msg *UserVerifyMessage) {
	MessageId++
	return &UserVerifyMessage{BaseMessageInfo: BaseMessageInfo{id: MessageId, priority: pri, trans: tra, birthday: time.Now()}, UserInfo: user}
}

func (d *UserVerifyMessage) String() string {
	return d.UserInfo.String() + "," + d.BaseMessageInfo.String()
}

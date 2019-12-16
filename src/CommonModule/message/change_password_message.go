package message

import (
	"CommonModule"
	"time"
)

// 修改用户密码的消息
type ChangePasswordMessage struct {
	BaseMessageInfo
	common.ChangePassword
}

// 生成消息
func NewChangePasswordMessage(user common.ChangePassword, pri common.Priority, tra TransType) (msg *ChangePasswordMessage) {
	MessageId++
	return &ChangePasswordMessage{BaseMessageInfo: BaseMessageInfo{id: MessageId, priority: pri, trans: tra, birthday: time.Now()}, ChangePassword: user}
}

func (d *ChangePasswordMessage) String() string {
	return d.ChangePassword.String() + "," + d.BaseMessageInfo.String()
}

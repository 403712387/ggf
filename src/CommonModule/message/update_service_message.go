package message

import (
	"CommonModule"
	"fmt"
	"time"
)

type UpdateInfo struct {
	UpdateTime string `json:"update_time"` // 升级时间，立即升级的话，为空或者"now",否则的话就是 "2010-01-02 10:20:23"这种格式
}

func (u *UpdateInfo) String() string {
	return fmt.Sprintf("time:%s", u.UpdateTime)
}

// 升级服务的消息
type UpdateServiceMessage struct {
	BaseMessageInfo
	FileData []byte     // 升级包/安装包的内容
	FileName string     // 升级包/安装包的名称
	Info     UpdateInfo // 升级的信息
}

// 生成消息
func NewUpdateServiceMessage(info UpdateInfo, fileName string, fileData []byte, pri common.Priority, tra TransType) (msg *UpdateServiceMessage) {
	MessageId++
	return &UpdateServiceMessage{BaseMessageInfo: BaseMessageInfo{id: MessageId, priority: pri, trans: tra, birthday: time.Now()}, Info: info, FileName: fileName, FileData: fileData}
}

func (u *UpdateServiceMessage) String() string {
	return fmt.Sprintf("info:%s, file name:%s, file data length:%d, %s", u.Info.String(), u.FileName, len(u.FileData), u.BaseMessageInfo.String())
}

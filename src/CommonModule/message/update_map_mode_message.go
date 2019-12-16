package message

import (
	common "CommonModule"
	"fmt"
	"time"
)

//地图配置
type MapInfo struct {
	Mode	int `json:"mode"`    // 地图模式
}

type UpdateMapModeMessage struct {
	BaseMessageInfo
	FileData []byte     // 地图压缩包的内容
	FileName string     // 地图压缩包的名称
	Info     MapInfo // 地图模式
}

func (u *MapInfo) String() string {
	return fmt.Sprintf("mode:%d", u.Mode)
}

// 生成消息
func NewUpdateMapModeMessage(info MapInfo, fileName string, fileData []byte, pri common.Priority, tra TransType) (msg *UpdateMapModeMessage) {
	MessageId++
	return &UpdateMapModeMessage{BaseMessageInfo: BaseMessageInfo{id: MessageId, priority: pri, trans: tra, birthday: time.Now()}, Info: info, FileName: fileName, FileData: fileData}
}

func (u *UpdateMapModeMessage) String() string {
	return fmt.Sprintf("info:%s, file name:%s, file data length:%d, %s", u.Info.String(), u.FileName, len(u.FileData), u.BaseMessageInfo.String())
}

// 地图的配置信息
type GetMapModeInfoMessage struct {
	BaseMessageInfo
}

// 生成消息
func NewGetMapModeInfoMessage(pri common.Priority, tra TransType) (msg *GetMapModeInfoMessage) {
	MessageId++
	return &GetMapModeInfoMessage{BaseMessageInfo: BaseMessageInfo{id: MessageId, priority: pri, trans: tra, birthday: time.Now()}}
}

func (p *GetMapModeInfoMessage) String() (result string) {

	result = p.BaseMessageInfo.String()
	return
}

type MapModeInfoResponse struct {
	BaseResponseInfo
	Info common.MapModeInfo // 地图模式的信息
}

// 生成回应
func NewMapModeInfoResponse(info common.MapModeInfo, msg BaseMessage) (response *MapModeInfoResponse) {
	resp := MapModeInfoResponse{BaseResponseInfo: BaseResponseInfo{message: msg}, Info: info}
	return &resp
}

func (s *MapModeInfoResponse) String() (result string) {

	result = s.BaseResponseInfo.String()
	return
}



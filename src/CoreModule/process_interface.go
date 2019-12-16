package core

import "CommonModule/message"

// 处理消息的接口
type ProcessInterface interface {
	// 初始化
	Init()

	// 反初始化
	Uninit()

	// 开始工作
	BeginWork()

	// 停止工作
	StopWork()

	// 开启消息队列
	BeginProcessLoop()

	// 退出消息队列
	StopProcessLoop()

	// 获取模块名称
	GetModuleName() string

	// 把消息放入消息队列
	PushMessage(message message.BaseMessage)

	// 把消息的回应放入回应队列
	PushResponse(response message.BaseResponse)

	// 偷窥消息
	OnForeseeMessage(message message.BaseMessage)(done bool)

	// 处理消息
	OnProcessMessage(message message.BaseMessage)(response message.BaseResponse, err error)

	// 偷窥消息的回应
	OnForeseeResponse(response message.BaseResponse)(done bool)

	// 处理消息的回应
	OnProcessResponse(response message.BaseResponse)

	// 发送消息
	SendMessage(message message.BaseMessage)(response message.BaseResponse, err error)
}
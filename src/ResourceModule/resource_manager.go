package ResourceModule

import (
	"CommonModule"
	"CommonModule/message"
	"CoreModule"
	"github.com/sirupsen/logrus"
	"sync"
	"time"
)

/*
本模块获取资源的使用情况（CPU,内存，网络，硬盘）
*/
const loopCount int = 5

type ResourceManager struct {
	core.MessageList                 // 消息列表
	Service          map[string]bool // 本机的所有服务组件的名称
	ServiceLock      sync.RWMutex
	SamplingRate     int                    // 采样频率
	notify           [loopCount]chan string //
	LastAccess       time.Time              // 用户最后一次查看资源使用的时间
}

// 初始化
func (r *ResourceManager) Init() {
	logrus.Infof("begin %s module init", r.ModuleName)
	r.Service = make(map[string]bool)
	r.SamplingRate = 30
	for i := 0; i < loopCount; i++ {
		r.notify[i] = make(chan string)
	}
	r.LastAccess = time.Now().AddDate(0, 0, -1)
	logrus.Infof("end %s module init", r.ModuleName)
}

// 反初始化
func (r *ResourceManager) Uninit() {
	logrus.Infof("begin %s module uninit", r.ModuleName)
	logrus.Infof("end %s module uninit", r.ModuleName)
}

// 开始工作
func (r *ResourceManager) BeginWork() {
	logrus.Infof("begin %s module beginwork", r.ModuleName)

	// cpu使用情况采集
	go r.cpuStatisticLoop(0)

	// 内存使用情况采集
	go r.memoryStatisticLoop(1)

	//  网络使用情况采集
	go r.networkStatisticLoop(2)

	// 硬盘使用情况采集
	go r.DiskStatisticLoop(3)

	// 修改采样频率的线程
	go r.changeSamplingRateLoop()
	logrus.Infof("end %s module beginwork", r.ModuleName)
}

// 停止工作
func (r *ResourceManager) StopWork() {
	logrus.Infof("begin %s module stopwork", r.ModuleName)
	logrus.Infof("end %s module stopwork", r.ModuleName)
}

// 偷窥消息
func (r *ResourceManager) OnForeseeMessage(msg message.BaseMessage) (done bool) {
	switch msg.(type) {
	case *message.GetCpuStatisticMessage:
		r.processSomeoneIsLooking()
	case *message.GetMemoryStatisticMessage:
		r.processSomeoneIsLooking()
	case *message.GetDiskStatisticMessage:
		r.processSomeoneIsLooking()
	case *message.GetNetworkStatisticMessage:
		r.processSomeoneIsLooking()
	}
	return
}

// 处理消息
func (r *ResourceManager) OnProcessMessage(msg message.BaseMessage) (rsp message.BaseResponse, err error) {
	switch msg.(type) {
	}
	return nil, nil
}

// 偷窥消息的回应
func (r *ResourceManager) OnForeseeResponse(rsp message.BaseResponse) (done bool) {
	return
}

// 处理消息的回应
func (r *ResourceManager) OnProcessResponse(rsp message.BaseResponse) {
	return
}

// 如果有人正在页面上看资源使用情况，那么就两秒钟采样一次，如果没有人看的时候，就30秒采样一次
func (r *ResourceManager) processSomeoneIsLooking() {
	r.LastAccess = time.Now()

	// 采样频率改为2秒一次
	if r.SamplingRate > 2 {
		r.SamplingRate = 2
		logrus.Infof("change sampling rate time %d", r.SamplingRate)

		// 通知所有的资源采集的loop
		for i := 0; i < loopCount; i++ {
			r.notify[i] <- "hello,panda"
		}
	}
}

// cpu使用情况采集
func (r *ResourceManager) cpuStatisticLoop(index int) {
	logrus.Infof("begin cpu statistic loop")
	for {

		// 获取CPU使用
		begin := time.Now()
		cpu, err := statisticCPU()
		if err == nil {
			msg := message.NewCPUStatisticMessage(time.Now(), cpu, common.Priority_Third, message.Trans_Async)
			r.SendMessage(msg)
		} else {
			logrus.Errorf("statistic cpu fail, err %s", err.Error())
		}

		// 休眠
		select {
		case <-time.After(time.Duration(r.SamplingRate)*time.Second - (time.Since(begin))):
			continue
		case data := <-r.notify[index]:
			logrus.Infof("%d change sampling rate, %s", index, data)
			continue
		}
	}
	logrus.Infof("begin cpu statistic loop")
}

// 内存使用情况采集
func (r *ResourceManager) memoryStatisticLoop(index int) {
	logrus.Infof("begin memory statistic loop")
	for {
		// 获取内存使用
		begin := time.Now()
		memory, err := statisticMemory()
		if err == nil {
			msg := message.NewMemoryInfoMessage(time.Now(), memory, common.Priority_Third, message.Trans_Async)
			r.SendMessage(msg)
		} else {
			logrus.Errorf("statistic memory fail, err %s", err.Error())
		}

		// 休眠
		select {
		case <-time.After(time.Duration(r.SamplingRate)*time.Second - (time.Since(begin))):
			continue
		case data := <-r.notify[index]:
			logrus.Infof("%d change sampling rate, %s", index, data)
			continue
		}
	}
	logrus.Infof("begin memory statistic loop")
}

//  网络使用情况采集
func (r *ResourceManager) networkStatisticLoop(index int) {
	logrus.Infof("begin network statistic loop")
	for {
		// 统计网络使用
		begin := time.Now()
		networks, err := statisticNetwork()
		if err == nil {
			msg := message.NewNetworkStatisticMessage(time.Now(), networks, common.Priority_Third, message.Trans_Async)
			r.SendMessage(msg)
		} else {
			logrus.Errorf("statistic network fail, err %s", err.Error())
		}

		// 休眠
		select {
		case <-time.After(time.Duration(r.SamplingRate)*time.Second - (time.Since(begin))):
			continue
		case data := <-r.notify[index]:
			logrus.Infof("%d change sampling rate, %s", index, data)
			continue
		}
	}
	logrus.Infof("begin network statistic loop")
}

// 硬盘使用情况采集
func (r *ResourceManager) DiskStatisticLoop(index int) {
	logrus.Infof("begin disk statistic loop")
	for {
		// 统计磁盘使用
		begin := time.Now()
		disk, err := statisticDisk()
		if err == nil {
			msg := message.NewDiskStatisticMessage(time.Now(), disk, common.Priority_Third, message.Trans_Async)
			r.SendMessage(msg)
		} else {
			logrus.Errorf("statistic disk fail, error:%s", err.Error())
		}

		// 休眠
		select {
		case <-time.After(time.Duration(r.SamplingRate)*time.Second - (time.Since(begin))):
			continue
		case data := <-r.notify[index]:
			logrus.Infof("%d change sampling rate, %s", index, data)
			continue
		}
	}
	logrus.Infof("end disk statistic loop")
}

// 修改采样频率的线程
func (r *ResourceManager) changeSamplingRateLoop() {
	logrus.Infof("begin change sampling rate loop")
	for {
		time.Sleep(5 * time.Second)

		// 判断时间是不是被修改过
		now := time.Now()
		if now.Sub(r.LastAccess).Seconds() < 0 {
			r.LastAccess = now
		}

		// 采样频率为30秒一次，无需修改
		if r.SamplingRate > 2 {
			continue
		}

		// 如果有30秒没有人访问，则采样频率改为30秒
		if now.Sub(r.LastAccess).Seconds() >= 30 {
			r.SamplingRate = 30
			logrus.Infof("change sampling rate time %d", r.SamplingRate)
		}
	}
	logrus.Infof("end change sampling rate loop")
}

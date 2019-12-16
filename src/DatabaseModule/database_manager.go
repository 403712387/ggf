package DatabaseModule

import (
	"CommonModule"
	"CommonModule/message"
	"CoreModule"
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"github.com/sirupsen/logrus"
	"os"
	"strings"
	"sync"
	"time"
)

/*
本模块是和数据库交互的模块,使用SQLite数据库
*/
const msgCapacity int = 20

type DatabaseManager struct {
	core.MessageList        // 消息列表
	databaseDir      string // 数据库存放的目录
	databaseName     string // 数据库的名称

	db     *sql.DB         // 操作数据库的句柄
	dbLock sync.Mutex      // 不能并行操作

	cpuStatistic []message.BaseMessage // cpu使用统计
	cpuLock      sync.Mutex

	diskStatistic []message.BaseMessage // 磁盘的使用统计
	diskLock      sync.Mutex

	networkStatistic []message.BaseMessage // 网络的使用统计
	networkLock      sync.Mutex

	memoryStatistic []message.BaseMessage // 内存的使用统计
	memoryLock      sync.Mutex

	serviceStatistic []message.BaseMessage // 服务组件的资源使用情况
	serviceLock      sync.Mutex
}

// 初始化
func (d *DatabaseManager) Init() {
	logrus.Infof("begin %s module uninit", d.ModuleName)
	d.databaseDir = "./database/"
	d.databaseName = "host.db"
	d.initDatabase()
	go d.checkDatabaseLoop()
	logrus.Infof("end %s module uninit", d.ModuleName)
}

// 反初始化
func (d *DatabaseManager) Uninit() {
	logrus.Infof("begin %s module uninit", d.ModuleName)
	logrus.Infof("end %s module uninit", d.ModuleName)
}

// 开始工作
func (d *DatabaseManager) BeginWork() {
	logrus.Infof("begin %s module beginwork", d.ModuleName)
	logrus.Infof("end %s module beginwork", d.ModuleName)
}

// 停止工作
func (d *DatabaseManager) StopWork() {
	logrus.Infof("begin %s module stopwork", d.ModuleName)
	logrus.Infof("end %s module stopwork", d.ModuleName)
}

// 偷窥消息
func (d *DatabaseManager) OnForeseeMessage(msg message.BaseMessage) (done bool) {
	return
}

// 处理消息
func (d *DatabaseManager) OnProcessMessage(msg message.BaseMessage) (rsp message.BaseResponse, err error) {
	switch msg.(type) {
	case *message.CPUStatisticMessage: // 统计CPU的使用情况
		return d.processCPUStatisticMessage(msg)
	case *message.DiskStatisticMessage: // 统计磁盘的使用情况
		return d.processDiskStatisticMessage(msg)
	case *message.NetworkStatisticMessage: // 统计网络使用情况
		return d.processNetworkStatisticMessage(msg)
	case *message.MemoryStatisticMessage: // 统计内存使用情况
		return d.processMemoryStatisticMessage(msg)
	case *message.ServiceStatisticMessage: // 统计服务组件的资源使用情况
		return d.processServiceStatisticMessage(msg)
	case *message.EventMessage: // 事件
		return d.processEventMessage(msg)
	case *message.GetCpuStatisticMessage: // 获取CPU使用情况
		return d.processGetCpuStatisticMessage(msg)
	case *message.GetDiskStatisticMessage: // 获取磁盘使用情况
		return d.processGetDiskStatisticMessage(msg)
	case *message.GetNetworkStatisticMessage: // 获取网络使用情况
		return d.processGetNetworkStatisticMessage(msg)
	case *message.GetMemoryStatisticMessage: // 获取内存使用情况
		return d.processGetMemoryStatisticMessage(msg)
	}
	return nil, nil
}

// 偷窥消息的回应
func (d *DatabaseManager) OnForeseeResponse(rsp message.BaseResponse) (done bool) {
	return
}

// 处理消息的回应
func (d *DatabaseManager) OnProcessResponse(rsp message.BaseResponse) {
	return
}

// 初始化数据库（如果数据库不存在，则创建）
func (d *DatabaseManager) initDatabase() (err error) {
	//  打开/创建数据库
	os.Mkdir(d.databaseDir, os.ModeDir)
	db, err := sql.Open("sqlite3", d.databaseDir+d.databaseName)
	if err != nil {
		logrus.Errorf("create database %s fail, error reason:%s", d.databaseDir, err.Error())
		return
	}

	// 创建表
	d.db = db
	return d.createTables()
}

// 检查数据库（定期删除两个月前的数据表）
func (d *DatabaseManager) checkDatabaseLoop() {
	logrus.Infof("begin check database loop %s", d.ModuleName)
	for {

		// 删除十五天前的表
		expire := time.Now().AddDate(0, 0, -15)
		var expireTables []string

		// 获取数据库中所有的表
		tables := getAllTables(d)

		// 获取十五天前的表
		for _, table := range tables {
			if index := strings.Index(table, "_"); index > 0 {
				strTime := table[index+1:]
				_, err := time.Parse("20060102", strTime)
				if err != nil {
					continue
				}

				// 判断是否超过一定时间
				if expire.Format("20060102") > strTime {
					expireTables = append(expireTables, table)
				}
			}
		}

		// 删除过期的表
		removeTables(d, expireTables)

		// 休眠
		time.Sleep(time.Hour * 24)
	}
	logrus.Infof("end check database loop %s", d.ModuleName)
}

// 创建数据表
func (d *DatabaseManager) createTables() (err error) {

	now := time.Now().Format("20060102")
	err = d.createEventTable(now)
	if err != nil {
		logrus.Errorf("create event table fail, error reason:%s", err.Error())
	}

	err = d.createSystemResourceTable(now)
	if err != nil {
		logrus.Errorf("create system resource table fail, error reason:%s", err.Error())
	}

	return
}

// 创建事件表
func (d *DatabaseManager) createEventTable(now string) (err error) {
	d.dbLock.Lock()
	defer d.dbLock.Unlock()

	statement := `CREATE TABLE IF NOT EXISTS Event_%s(
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			type VARCHAR(64) NULL,
			birthday VARCHAR(64) NULL,
			recordTime VARCHAR(64) NULL,
			explain VARCHAR(1024) NULL
			);`
	_, err = d.db.Exec(fmt.Sprintf(statement, now))
	if err != nil {
		logrus.Errorf("create table Event_%s fail, error:%s", now, err.Error())
	}
	return
}

// 创建系统资源使用情况表
func (d *DatabaseManager) createSystemResourceTable(now string) (err error) {
	d.dbLock.Lock()
	defer d.dbLock.Unlock()

	statement := `CREATE TABLE IF NOT EXISTS SystemResource_%s(
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			time VARCHAR(64) NULL,
			type VARCHAR(64) NULL,
			info VARCHAR(1048576) NULL,
			explain VARCHAR(512) NULL
			);`
	_, err = d.db.Exec(fmt.Sprintf(statement, now))
	if err != nil {
		logrus.Errorf("create table SystemResource_%s fail, error:%s", now, err.Error())
	}
	return
}

// cpu的使用情况
func (d *DatabaseManager) processCPUStatisticMessage(msg message.BaseMessage) (rsp message.BaseResponse, err error) {
	statisticMsg := msg.(*message.CPUStatisticMessage)

	// 插入队列中
	d.cpuStatistic = d.insertMessage(msg, d.cpuStatistic, &d.cpuLock)

	// 插入到数据库中
	data, err := json.Marshal(statisticMsg.CapacityInfo)
	err = insertResourceUsedInfo(d, statisticMsg.Time, "cpu", data)

	// 如果插入失败，重新插入一次
	if err != nil {
		tableName := statisticMsg.Time.Format("20060102")
		d.createSystemResourceTable(tableName)
		err = insertResourceUsedInfo(d, statisticMsg.Time, "cpu", data)
	}

	return
}

// 磁盘的使用情况
func (d *DatabaseManager) processDiskStatisticMessage(msg message.BaseMessage) (rsp message.BaseResponse, err error) {
	statisticMsg := msg.(*message.DiskStatisticMessage)

	// 插入队列中
	d.diskStatistic = d.insertMessage(msg, d.diskStatistic, &d.diskLock)

	// 插入到数据库中
	data, err := json.Marshal(statisticMsg.Statistic)
	err = insertResourceUsedInfo(d, statisticMsg.Time, "disk", data)

	// 如果插入失败，重新插入一次
	if err != nil {
		tableName := statisticMsg.Time.Format("20060102")
		d.createSystemResourceTable(tableName)
		err = insertResourceUsedInfo(d, statisticMsg.Time, "disk", data)
	}

	return
}

// 网络的使用情况
func (d *DatabaseManager) processNetworkStatisticMessage(msg message.BaseMessage) (rsp message.BaseResponse, err error) {
	statisticMsg := msg.(*message.NetworkStatisticMessage)

	// 插入队列中
	d.networkStatistic = d.insertMessage(msg, d.networkStatistic, &d.networkLock)

	// 插入到数据库中
	data, err := json.Marshal(statisticMsg.Statistic)
	err = insertResourceUsedInfo(d, statisticMsg.Time, "network", data)

	// 如果插入失败，重新插入一次
	if err != nil {
		tableName := statisticMsg.Time.Format("20060102")
		d.createSystemResourceTable(tableName)
		err = insertResourceUsedInfo(d, statisticMsg.Time, "network", data)
	}

	return
}

// 内存的使用情况
func (d *DatabaseManager) processMemoryStatisticMessage(msg message.BaseMessage) (rsp message.BaseResponse, err error) {
	statisticMsg := msg.(*message.MemoryStatisticMessage)

	// 插入队列中
	d.memoryStatistic = d.insertMessage(msg, d.memoryStatistic, &d.memoryLock)

	// 插入到数据库中
	data, err := json.Marshal(statisticMsg.MemoryStatisticInfo)
	err = insertResourceUsedInfo(d, statisticMsg.Time, "memory", data)

	// 如果插入失败，重新插入一次
	if err != nil {
		tableName := statisticMsg.Time.Format("20060102")
		d.createSystemResourceTable(tableName)
		err = insertResourceUsedInfo(d, statisticMsg.Time, "memory", data)
	}
	return
}

// 服务组件的资源使用情况
func (d *DatabaseManager) processServiceStatisticMessage(msg message.BaseMessage) (rsp message.BaseResponse, err error) {
	statisticMsg := msg.(*message.ServiceStatisticMessage)

	// 插入队列中
	d.serviceStatistic = d.insertMessage(msg, d.serviceStatistic, &d.serviceLock)

	// 插入到数据库中
	data, err := json.Marshal(statisticMsg.Statistic)
	err = insertResourceUsedInfo(d, statisticMsg.Time, "service", data)

	// 如果插入失败，重新插入一次
	if err != nil {
		tableName := statisticMsg.Time.Format("20060102")
		d.createSystemResourceTable(tableName)
		err = insertResourceUsedInfo(d, statisticMsg.Time, "service", data)
	}
	return
}

// 记录事件的信息
func (d *DatabaseManager) processEventMessage(msg message.BaseMessage) (rsp message.BaseResponse, err error) {
	eventMsg := msg.(*message.EventMessage)

	// 插入数据库中
	err = insertEvent(d, eventMsg.Type, eventMsg.Explain, eventMsg.Birthday())

	// 如果插入失败，则再插入一次
	if err != nil {
		tableName := eventMsg.Birthday().Format("20060102")
		d.createEventTable(tableName)
		err = insertEvent(d, eventMsg.Type, eventMsg.Explain, eventMsg.Birthday())
	}
	return
}

// 获取CPU使用情况
func (d *DatabaseManager) processGetCpuStatisticMessage(msg message.BaseMessage) (rsp message.BaseResponse, err error) {
	getMsg := msg.(*message.GetCpuStatisticMessage)
	begin := strings.Trim(getMsg.BeginTime, " ")
	end := strings.Trim(getMsg.EndTime, " ")

	var times []time.Time
	var statistic []common.CapacityInfo
	if len(begin) <= 0 && len(end) <= 0 { // 如果begin和end均为空，则返回内存中统计信息
		times, statistic = d.getCpuStatisticInBuffer()
	} else if len(begin) > 0 && len(end) > 0 { // 获取指定范围的统计信息
		beginTime, e := time.Parse("2006-01-02 15:04:05", begin)
		if e == nil {
			endTime, e := time.Parse("2006-01-02 15:04:05", end)
			if e == nil {
				times, statistic = d.getCpuStatisticByTime(beginTime, endTime)
			} else {
				err = e
			}
		} else {
			err = e
		}
	} else if len(begin) > 0 && len(end) <= 0 { // 获取最新的统计信息
		beginTime, e := time.Parse("2006-01-02 15:04:05", begin)
		if e == nil {
			times, statistic = d.getCpuStatisticByTime(beginTime, time.Now())
		} else {
			err = e
		}
	} else {
		err = fmt.Errorf("invalid request, begin:%s, end:%s", begin, end)
	}
	rsp = message.NewGetCpuStatisticResponse(times, statistic, msg)
	return
}

// 获取内存中的cpu统计信息
func (d *DatabaseManager) getCpuStatisticInBuffer() (samplingTime []time.Time, stat []common.CapacityInfo) {
	d.cpuLock.Lock()
	defer d.cpuLock.Unlock()

	// 遍历slice,获取所有的统计信息
	for _, msg := range d.cpuStatistic {
		statisticMsg := msg.(*message.CPUStatisticMessage)
		samplingTime = append(samplingTime, statisticMsg.Time)
		stat = append(stat, statisticMsg.CapacityInfo)
	}
	return
}

// 根据时间获取cpu的统计信息
func (d *DatabaseManager) getCpuStatisticByTime(begin, end time.Time) (samplingTime []time.Time, stat []common.CapacityInfo) {
	d.cpuLock.Lock()
	defer d.cpuLock.Unlock()

	beginTime := begin.Format("2006-01-02 15:04:05")
	endTime := end.Format("2006-01-02 15:04:05")

	// 先检测看内存中的数据是否满足查询条件
	if len(d.cpuStatistic) > 0 {
		firstMsg := d.cpuStatistic[0].(*message.CPUStatisticMessage)
		firstTime := firstMsg.Time.Format("2006-01-02 15:04:05")
		if beginTime >= firstTime {
			for _, msg := range d.cpuStatistic {
				statisticMsg := msg.(*message.CPUStatisticMessage)
				statisticTime := statisticMsg.Time.Format("2006-01-02 15:04:05")
				if statisticTime >= beginTime && statisticTime <= endTime {
					samplingTime = append(samplingTime, statisticMsg.Time)
					stat = append(stat, statisticMsg.CapacityInfo)
				}
			}
			return
		}
	}

	// 如果内存中的数据不满足查询条件，则向数据库中查询
	times, datas, _ := queryResourceStatistic(d, "cpu", begin, end)
	for index, data := range datas {
		tmpTime, _ := time.Parse("2006-01-02 15:04:05", times[index])
		samplingTime = append(samplingTime, tmpTime)

		var tmpStat common.CapacityInfo
		json.Unmarshal([]byte(data), &tmpStat)
		stat = append(stat, tmpStat)
	}
	return
}

// 获取磁盘使用情况
func (d *DatabaseManager) processGetDiskStatisticMessage(msg message.BaseMessage) (rsp message.BaseResponse, err error) {
	getMsg := msg.(*message.GetDiskStatisticMessage)
	begin := strings.Trim(getMsg.BeginTime, " ")
	end := strings.Trim(getMsg.EndTime, " ")

	var times []time.Time
	var statistic [][]common.DiskStatisticInfo
	if len(begin) <= 0 && len(end) <= 0 { // 如果begin和end均为空，则返回内存中统计信息
		times, statistic = d.getDiskStatisticInBuffer()
	} else if len(begin) > 0 && len(end) > 0 { // 获取指定范围的统计信息
		beginTime, e := time.Parse("2006-01-02 15:04:05", begin)
		if e == nil {
			endTime, e := time.Parse("2006-01-02 15:04:05", end)
			if e == nil {
				times, statistic = d.getDiskStatisticByTime(beginTime, endTime)
			} else {
				err = e
			}
		} else {
			err = e
		}
	} else if len(begin) > 0 && len(end) <= 0 { // 获取最新的统计信息
		beginTime, e := time.Parse("2006-01-02 15:04:05", begin)
		if e == nil {
			times, statistic = d.getDiskStatisticByTime(beginTime, time.Now())
		} else {
			err = e
		}
	} else {
		err = fmt.Errorf("invalid request, begin:%s, end:%s", begin, end)
	}
	rsp = message.NewGetDiskStatisticResponse(times, statistic, msg)
	return
}

// 获取内存中的磁盘统计信息
func (d *DatabaseManager) getDiskStatisticInBuffer() (samplingTime []time.Time, stat [][]common.DiskStatisticInfo) {
	d.diskLock.Lock()
	defer d.diskLock.Unlock()

	// 遍历slice,获取所有的统计信息
	for _, msg := range d.diskStatistic {
		statisticMsg := msg.(*message.DiskStatisticMessage)
		samplingTime = append(samplingTime, statisticMsg.Time)
		stat = append(stat, statisticMsg.Statistic)
	}
	return
}

// 根据时间获取磁盘的统计信息
func (d *DatabaseManager) getDiskStatisticByTime(begin, end time.Time) (samplingTime []time.Time, stat [][]common.DiskStatisticInfo) {
	d.diskLock.Lock()
	defer d.diskLock.Unlock()

	beginTime := begin.Format("2006-01-02 15:04:05")
	endTime := end.Format("2006-01-02 15:04:05")

	// 先检测看内存中的数据是否满足查询条件
	if len(d.diskStatistic) > 0 {
		firstMsg := d.diskStatistic[0].(*message.DiskStatisticMessage)
		firstTime := firstMsg.Time.Format("2006-01-02 15:04:05")
		if beginTime >= firstTime {
			for _, msg := range d.diskStatistic {
				statisticMsg := msg.(*message.DiskStatisticMessage)
				statisticTime := statisticMsg.Time.Format("2006-01-02 15:04:05")
				if statisticTime >= beginTime && statisticTime <= endTime {
					samplingTime = append(samplingTime, statisticMsg.Time)
					stat = append(stat, statisticMsg.Statistic)
				}
			}
			return
		}
	}

	// 如果内存中的数据不满足查询条件，则向数据库中查询
	times, datas, _ := queryResourceStatistic(d, "disk", begin, end)
	for index, data := range datas {
		tmpTime, _ := time.Parse("2006-01-02 15:04:05", times[index])
		samplingTime = append(samplingTime, tmpTime)

		var tmpStat []common.DiskStatisticInfo
		json.Unmarshal([]byte(data), &tmpStat)
		stat = append(stat, tmpStat)
	}
	return
}

// 获取网络使用情况
func (d *DatabaseManager) processGetNetworkStatisticMessage(msg message.BaseMessage) (rsp message.BaseResponse, err error) {
	getMsg := msg.(*message.GetNetworkStatisticMessage)
	begin := strings.Trim(getMsg.BeginTime, " ")
	end := strings.Trim(getMsg.EndTime, " ")

	var times []time.Time
	var statistic [][]common.NetworkStatisticInfo
	if len(begin) <= 0 && len(end) <= 0 { // 如果begin和end均为空，则返回内存中统计信息
		times, statistic = d.getNetworkStatisticInBuffer()
	} else if len(begin) > 0 && len(end) > 0 { // 获取指定范围的统计信息
		beginTime, e := time.Parse("2006-01-02 15:04:05", begin)
		if e == nil {
			endTime, e := time.Parse("2006-01-02 15:04:05", end)
			if e == nil {
				times, statistic = d.getNetworkStatisticByTime(beginTime, endTime)
			} else {
				err = e
			}
		} else {
			err = e
		}
	} else if len(begin) > 0 && len(end) <= 0 { // 获取最新的统计信息
		beginTime, e := time.Parse("2006-01-02 15:04:05", begin)
		if e == nil {
			times, statistic = d.getNetworkStatisticByTime(beginTime, time.Now())
		} else {
			err = e
		}
	} else {
		err = fmt.Errorf("invalid request, begin:%s, end:%s", begin, end)
	}
	rsp = message.NewGetNetworkStatisticResponse(times, statistic, msg)
	return
}

// 获取内存中的网络统计信息
func (d *DatabaseManager) getNetworkStatisticInBuffer() (samplingTime []time.Time, stat [][]common.NetworkStatisticInfo) {
	d.networkLock.Lock()
	defer d.networkLock.Unlock()

	// 遍历slice,获取所有的统计信息
	for _, msg := range d.networkStatistic {
		statisticMsg := msg.(*message.NetworkStatisticMessage)
		samplingTime = append(samplingTime, statisticMsg.Time)
		stat = append(stat, statisticMsg.Statistic)
	}
	return
}

// 根据时间获取网络的统计信息
func (d *DatabaseManager) getNetworkStatisticByTime(begin, end time.Time) (samplingTime []time.Time, stat [][]common.NetworkStatisticInfo) {
	d.networkLock.Lock()
	defer d.networkLock.Unlock()

	beginTime := begin.Format("2006-01-02 15:04:05")
	endTime := end.Format("2006-01-02 15:04:05")

	// 先检测看内存中的数据是否满足查询条件
	if len(d.networkStatistic) > 0 {
		firstMsg := d.networkStatistic[0].(*message.NetworkStatisticMessage)
		firstTime := firstMsg.Time.Format("2006-01-02 15:04:05")
		if beginTime >= firstTime {
			for _, msg := range d.networkStatistic {
				statisticMsg := msg.(*message.NetworkStatisticMessage)
				statisticTime := statisticMsg.Time.Format("2006-01-02 15:04:05")
				if statisticTime >= beginTime && statisticTime <= endTime {
					samplingTime = append(samplingTime, statisticMsg.Time)
					stat = append(stat, statisticMsg.Statistic)
				}
			}
			return
		}
	}

	// 如果内存中的数据不满足查询条件，则向数据库中查询
	times, datas, _ := queryResourceStatistic(d, "network", begin, end)
	for index, data := range datas {
		tmpTime, _ := time.Parse("2006-01-02 15:04:05", times[index])
		samplingTime = append(samplingTime, tmpTime)

		var tmpStat []common.NetworkStatisticInfo
		json.Unmarshal([]byte(data), &tmpStat)
		stat = append(stat, tmpStat)
	}
	return
}

// 获取内存使用情况
func (d *DatabaseManager) processGetMemoryStatisticMessage(msg message.BaseMessage) (rsp message.BaseResponse, err error) {
	getMsg := msg.(*message.GetMemoryStatisticMessage)
	begin := strings.Trim(getMsg.BeginTime, " ")
	end := strings.Trim(getMsg.EndTime, " ")

	var times []time.Time
	var statistic []common.MemoryStatisticInfo
	if len(begin) <= 0 && len(end) <= 0 { // 如果begin和end均为空，则返回内存中统计信息
		times, statistic = d.getMemoryStatisticInBuffer()
	} else if len(begin) > 0 && len(end) > 0 { // 获取指定范围的统计信息
		beginTime, e := time.Parse("2006-01-02 15:04:05", begin)
		if e == nil {
			endTime, e := time.Parse("2006-01-02 15:04:05", end)
			if e == nil {
				times, statistic = d.getMemoryStatisticByTime(beginTime, endTime)
			} else {
				err = e
			}
		} else {
			err = e
		}
	} else if len(begin) > 0 && len(end) <= 0 { // 获取最新的统计信息
		beginTime, e := time.Parse("2006-01-02 15:04:05", begin)
		if e == nil {
			times, statistic = d.getMemoryStatisticByTime(beginTime, time.Now())
		} else {
			err = e
		}
	} else {
		err = fmt.Errorf("invalid request, begin:%s, end:%s", begin, end)
	}
	rsp = message.NewGetMemoryStatisticResponse(times, statistic, msg)
	return
}

// 获取内存中的memory统计信息
func (d *DatabaseManager) getMemoryStatisticInBuffer() (samplingTime []time.Time, stat []common.MemoryStatisticInfo) {
	d.memoryLock.Lock()
	defer d.memoryLock.Unlock()

	// 遍历slice,获取所有的统计信息
	for _, msg := range d.memoryStatistic {
		statisticMsg := msg.(*message.MemoryStatisticMessage)
		samplingTime = append(samplingTime, statisticMsg.Time)
		stat = append(stat, statisticMsg.MemoryStatisticInfo)
	}
	return
}

// 根据时间获取memory的统计信息
func (d *DatabaseManager) getMemoryStatisticByTime(begin, end time.Time) (samplingTime []time.Time, stat []common.MemoryStatisticInfo) {
	d.memoryLock.Lock()
	defer d.memoryLock.Unlock()

	beginTime := begin.Format("2006-01-02 15:04:05")
	endTime := end.Format("2006-01-02 15:04:05")

	// 先检测看内存中的数据是否满足查询条件
	if len(d.memoryStatistic) > 0 {
		firstMsg := d.memoryStatistic[0].(*message.MemoryStatisticMessage)
		firstTime := firstMsg.Time.Format("2006-01-02 15:04:05")
		if beginTime >= firstTime {
			for _, msg := range d.memoryStatistic {
				statisticMsg := msg.(*message.MemoryStatisticMessage)
				statisticTime := statisticMsg.Time.Format("2006-01-02 15:04:05")
				if statisticTime >= beginTime && statisticTime <= endTime {
					samplingTime = append(samplingTime, statisticMsg.Time)
					stat = append(stat, statisticMsg.MemoryStatisticInfo)
				}
			}
			return
		}
	}

	// 如果内存中的数据不满足查询条件，则向数据库中查询
	times, datas, _ := queryResourceStatistic(d, "memory", begin, end)
	for index, data := range datas {
		tmpTime, _ := time.Parse("2006-01-02 15:04:05", times[index])
		samplingTime = append(samplingTime, tmpTime)

		var tmpStat common.MemoryStatisticInfo
		json.Unmarshal([]byte(data), &tmpStat)
		stat = append(stat, tmpStat)
	}
	return
}

// 获取内存中的服务组件统计信息
func (d *DatabaseManager) getServiceStatisticInBuffer() (samplingTime []time.Time, stat [][]common.ProcessInfo) {
	d.serviceLock.Lock()
	defer d.serviceLock.Unlock()

	// 遍历slice,获取所有的统计信息
	for _, msg := range d.serviceStatistic {
		statisticMsg := msg.(*message.ServiceStatisticMessage)
		samplingTime = append(samplingTime, statisticMsg.Time)
		stat = append(stat, statisticMsg.Statistic)
	}
	return
}

// 根据时间获取服务组件的统计信息
func (d *DatabaseManager) getServiceStatisticByTime(begin, end time.Time) (samplingTime []time.Time, stat [][]common.ProcessInfo) {
	d.serviceLock.Lock()
	defer d.serviceLock.Unlock()

	beginTime := begin.Format("2006-01-02 15:04:05")
	endTime := end.Format("2006-01-02 15:04:05")

	// 先检测看内存中的数据是否满足查询条件
	if len(d.serviceStatistic) > 0 {
		firstMsg := d.serviceStatistic[0].(*message.ServiceStatisticMessage)
		firstTime := firstMsg.Time.Format("2006-01-02 15:04:05")
		if beginTime >= firstTime {
			for _, msg := range d.serviceStatistic {
				statisticMsg := msg.(*message.ServiceStatisticMessage)
				statisticTime := statisticMsg.Time.Format("2006-01-02 15:04:05")
				if statisticTime >= beginTime && statisticTime <= endTime {
					samplingTime = append(samplingTime, statisticMsg.Time)
					stat = append(stat, statisticMsg.Statistic)
				}
			}
			return
		}
	}

	// 如果内存中的数据不满足查询条件，则向数据库中查询
	times, datas, _ := queryResourceStatistic(d, "service", begin, end)
	for index, data := range datas {
		tmpTime, _ := time.Parse("2006-01-02 15:04:05", times[index])
		samplingTime = append(samplingTime, tmpTime)

		var tmpStat []common.ProcessInfo
		json.Unmarshal([]byte(data), &tmpStat)
		stat = append(stat, tmpStat)
	}
	return
}

// 消息插入切片中
func (d *DatabaseManager) insertMessage(msg message.BaseMessage, sliceMsg []message.BaseMessage, lock *sync.Mutex) (result []message.BaseMessage) {
	lock.Lock()
	defer lock.Unlock()

	if len(sliceMsg) < msgCapacity {
		result = append(sliceMsg, msg)
	} else {
		result = append(sliceMsg[1:], msg)
	}
	return
}

// 打印
func (d *DatabaseManager) debug() {
	logrus.Infof("cpu message:%d, disk message:%d, network message:%d, memory message:%d, service message:%d", len(d.cpuStatistic), len(d.diskStatistic), len(d.networkStatistic), len(d.memoryStatistic), len(d.serviceStatistic))
}

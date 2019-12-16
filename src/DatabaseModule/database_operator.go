package DatabaseModule

import (
	"CommonModule"
	"fmt"
	"github.com/sirupsen/logrus"
	"time"
)

// 加载用户信息
func loadUserInfo(d *DatabaseManager) (user common.UserInfo, err error) {
	d.dbLock.Lock()
	defer d.dbLock.Unlock()

	sql := fmt.Sprintf("select name, password from %s", "User")
	rows, err := d.db.Query(sql)
	if err == nil {
		defer rows.Close()
		for rows.Next() {
			rows.Scan(&user.User, &user.Password)
			return
		}
		err = fmt.Errorf("not query user info")
		return
	} else {
		logrus.Errorf("load user info fail, reason:%s", err.Error())
	}
	return
}

// 获取所有的表
func getAllTables(d *DatabaseManager) (tables []string) {
	d.dbLock.Lock()
	defer d.dbLock.Unlock()

	// 查询数据
	sql := fmt.Sprintf(`select name from sqlite_master where type='table' order by name`)
	rows, err := d.db.Query(sql)
	if err == nil {
		for rows.Next() {
			var name string
			rows.Scan(&name)
			tables = append(tables, name)
		}
		rows.Close()
	}
	return
}

// 删除表
func removeTables(d *DatabaseManager, tables []string) (err error) {
	d.dbLock.Lock()
	defer d.dbLock.Unlock()

	// 循环删除表
	for _, table := range tables {
		sql := fmt.Sprintf("drop table %s", table)
		d.db.Exec(sql)
		logrus.Infof("remove table:%s, sql:%s", table, sql)
	}
	return
}

// 修改密码
func changePassword(d *DatabaseManager, user common.ChangePassword) (err error) {
	d.dbLock.Lock()
	defer d.dbLock.Unlock()

	// 根据用户名更新密码
	sql := fmt.Sprintf("update %s set password = '%s' where name = '%s'", "User", user.NewPassword, user.User)
	result, err := d.db.Exec(sql)

	// 更新内存中的用户名和密码
	if err != nil {
		logrus.Fatalf("change password fail, error reason:%s", err.Error())
	}
	count, err := result.RowsAffected()
	if count <= 0 {
		err = fmt.Errorf("update password fail")
	}

	return
}

// 插入系统资源使用
func insertResourceUsedInfo(d *DatabaseManager, samplingTime time.Time, dataType string, data []byte) (err error) {
	d.dbLock.Lock()
	defer d.dbLock.Unlock()

	// 查找到对应的表
	tableName := fmt.Sprintf("SystemResource_%s", samplingTime.Format("20060102"))

	// 构造sql语句
	sql := fmt.Sprintf(`insert into %s (%s, %s, %s) values ('%s', '%s', '%s')`, tableName, "time", "type", "info", samplingTime.Format("2006-01-02 15:04:05"), dataType, string(data[:]))
	_, err = d.db.Exec(sql)
	if err != nil {
		logrus.Errorf("insert data fail, sql:%s, error:%s", sql, err.Error())
	}
	return
}

// 记录事件
func insertEvent(d *DatabaseManager, eventType, explain string, birthday time.Time) (err error) {
	d.dbLock.Lock()
	defer d.dbLock.Unlock()

	// 查找到对应的表
	tableName := fmt.Sprintf("Event_%s", birthday.Format("20060102"))
	now := time.Now()

	// 构造sql语句
	sql := fmt.Sprintf(`insert into %s (%s, %s, %s, %s) values ('%s', '%s', '%s', '%s')`, tableName, "type", "birthday", "recordTime", "explain", eventType, birthday.Format("2006-01-02 15:04:05"), now.Format("2006-01-02 15:04:05"), explain)
	_, err = d.db.Exec(sql)
	if err != nil {
		logrus.Errorf("insert data fail, sql:%s, error:%s", sql, err.Error())
	}
	return
}

// 根据时间，查询资源使用情况
func queryResourceStatistic(d *DatabaseManager, dataType string, begin, end time.Time) (times []string, data []string, err error) {
	d.dbLock.Lock()
	defer d.dbLock.Unlock()

	index := begin
	endTime := end.Format("20060102")
	fullBeginTime := begin.Format("2006-01-02 15:04:05")
	fullEndTime := end.Format("2006-01-02 15:04:05")
	for {
		indexTime := index.Format("20060102")

		// 判断查询是否结束
		if indexTime > endTime {
			break
		}

		// 查询数据
		sql := fmt.Sprintf(`select time, info from %s_%s where type = "%s" and time >= "%s" and time <= "%s order by id"`, "SystemResource", indexTime, dataType, fullBeginTime, fullEndTime)
		rows, err := d.db.Query(sql)
		if err == nil {
			for rows.Next() {
				var time, info string
				rows.Scan(&time, &info)
				times = append(times, time)
				data = append(data, info)
			}
			rows.Close()
		}

		// 查询下一个表
		index = index.AddDate(0, 1, 0)
	}

	return
}

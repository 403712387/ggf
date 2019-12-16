package main

import (
	"ConfigureModule"
	"CoreModule"
	"DatabaseModule"
	"HostServiceModule"
	"KafkaModule"
	"LicenseModule"
	"NetworkModule"
	"NtpModule"
	"ResourceModule"
	"ServerModule"
	"ServiceModule"
	"StorageManager"
	"TimeZoneModule"
	"UpdateService"
	"flag"
	"github.com/sirupsen/logrus"
	"gopkg.in/natefinch/lumberjack.v2"
)

// 设置日志配置
func init() {
	customFormatter := new(logrus.TextFormatter)
	customFormatter.TimestampFormat = "2006-01-02 15:04:05.000000000"
	logrus.SetFormatter(customFormatter)
	logrus.SetOutput(&lumberjack.Logger{
		Filename:   "./logs/host.log",
		MaxSize:    10, // megabytes
		MaxBackups: 10,
		MaxAge:     10,    //days
		Compress:   false, // disabled by default
	})
}

func main() {

	// 解析日志级别
	level := flag.String("l", "info", "log level")
	flag.Parse()
	logLevel, err := logrus.ParseLevel(*level)
	if err != nil {
		logLevel = logrus.InfoLevel
	}
	logrus.SetLevel(logLevel)

	logrus.Info("----------------------Welcome-------------------------")
	var msgRoute core.MessageRoute
	confFile := "config/config.json"

	// 初始化模块
	conf := configure.ConfigureManager{MessageList: core.MessageList{MessageRoute: &msgRoute, ModuleName: "ConfigureModule"}, ConfigureFile: confFile} // 配置模块
	net := network.NetworkManager{MessageList: core.MessageList{MessageRoute: &msgRoute, ModuleName: "NetworkModule"}}                                 // 网络模块
	server := ServerModule.ServerManager{MessageList: core.MessageList{MessageRoute: &msgRoute, ModuleName: "ServerModule"}}                           // 服务器模块
	storage := StorageManager.StorageManager{MessageList: core.MessageList{MessageRoute: &msgRoute, ModuleName: "StorageManager"}}                     // 存储管理模块
	kafka := KafkaModule.KafkaManager{MessageList: core.MessageList{MessageRoute: &msgRoute, ModuleName: "KafkaManager"}}                              // kafka模块
	update := UpdateService.UpdateServiceManager{MessageList: core.MessageList{MessageRoute: &msgRoute, ModuleName: "UpdateServiceManager"}}           // 服务的安装升级模块
	database := DatabaseModule.DatabaseManager{MessageList: core.MessageList{MessageRoute: &msgRoute, ModuleName: "DatabaseManager"}}                  // 数据库交互的模块
	service := ServiceModule.ServiceManager{MessageList: core.MessageList{MessageRoute: &msgRoute, ModuleName: "ServiceManager"}}                      // 此模块是负责服务组件信息的获取
	ntp := NtpModule.NtpManager{MessageList: core.MessageList{MessageRoute: &msgRoute, ModuleName: "NtpManager"}}                                      // 此模块是负责服务组件信息的获取
	license := LicenseModule.LicenseManager{MessageList: core.MessageList{MessageRoute: &msgRoute, ModuleName: "LicenseManager"}}                      // license相关的模块
	timeZone := TimeZoneModule.TimeZoneManager{MessageList: core.MessageList{MessageRoute: &msgRoute, ModuleName: "TimeZoneManager"}}                  // 时区的模块
	resource := ResourceModule.ResourceManager{MessageList: core.MessageList{MessageRoute: &msgRoute, ModuleName: "ResourceManager"}}                  //  获取资源使用情况
	host := HostServiceModule.HostServiceManager{MessageList: core.MessageList{MessageRoute: &msgRoute, ModuleName: "HostServiceManager"}}

	// 初始化处理模块（注意，这个地方如果不初始化的话，会收不到异步消息）
	conf.Process = &conf
	net.Process = &net
	server.Process = &server
	storage.Process = &storage
	kafka.Process = &kafka
	update.Process = &update
	database.Process = &database
	service.Process = &service
	ntp.Process = &ntp
	license.Process = &license
	timeZone.Process = &timeZone
	resource.Process = &resource
	host.Process = &host

	// 模块添加到消息路由中
	msgRoute.AddProcess(&conf)
	msgRoute.AddProcess(&net)
	msgRoute.AddProcess(&server)
	msgRoute.AddProcess(&storage)
	msgRoute.AddProcess(&kafka)
	msgRoute.AddProcess(&update)
	msgRoute.AddProcess(&database)
	msgRoute.AddProcess(&service)
	msgRoute.AddProcess(&ntp)
	msgRoute.AddProcess(&license)
	msgRoute.AddProcess(&timeZone)
	msgRoute.AddProcess(&resource)
	msgRoute.AddProcess(&host)

	// 启动程序
	msgRoute.Beginwork()
	logrus.Info("----------------------Bye-------------------------")
}

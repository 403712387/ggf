package CommonUtil

import (
	"fmt"
	"io/ioutil"
	"strings"
)

// 解析hosts配置的方法（即解析/etc/hosts）

// 一个配置项
type HostsText struct {
	Key   string // 键
	Value string // 值
}

// 配置
type HostsSetting struct {
	FileName string // 配置文件的名称
	Conf     []HostsText
}

// 解析配置文件
func (c *HostsSetting) ParseConfigFile() (err error) {
	bytes, err := ioutil.ReadFile(c.FileName)
	if err != nil {
		return err
	}

	data := string(bytes)
	return c.Parse(&data)
}

// 解析配置
func (c *HostsSetting) Parse(data *string) (err error) {
	// 内容为空
	if data == nil || len(*data) <= 0 {
		return nil
	}

	// 获取每一行的配置
	lines := strings.Split(*data, "\n")
	for _, line := range lines {
		line = strings.Trim(line, " ")
		index := strings.Index(line, " ")

		// 不合法的配置项，丢弃
		if index <= 0 {
			continue
		}

		conf := HostsText{Key: strings.Trim(line[:index], " "), Value: strings.Trim(line[index+1:], " ")}
		c.Conf = append(c.Conf, conf)

	}
	return nil
}

// 配置项转成string
func (c *HostsSetting) String() (result string) {
	breakline := "\n"
	for _, conf := range c.Conf {

		// 如果value为空，则不保存该项
		conf.Value = strings.Trim(conf.Value, " ")
		if len(conf.Value) <= 0 {
			continue
		}
		result += conf.Key + " " + conf.Value + breakline
	}
	return
}

// 保存到文件
func (c *HostsSetting) Save() error {
	data := c.String()
	return ioutil.WriteFile(c.FileName, []byte(data), 0644)
}

// 是否存在对应的key
func (c *HostsSetting) Exist(key, value string) (result bool) {
	for _, conf := range c.Conf {
		if key == conf.Key && value == conf.Value {
			return true
		}
	}
	return false
}

// 是否存在对应的配置
func (c *HostsSetting) ExistConf(key, value string) (result bool) {
	for _, conf := range c.Conf {
		if key == conf.Key && value == conf.Value {
			return true
		}
	}
	return false
}

// 根据key获取Value
func (c *HostsSetting) Value(key string) (value string, err error) {
	for _, conf := range c.Conf {
		if key == conf.Key {
			return conf.Value, nil
		}
	}
	return key, fmt.Errorf("not find value %s", key)
}

// 根据value获取key
func (c *HostsSetting) Key(value string) (key string, err error) {
	for _, conf := range c.Conf {
		if value == conf.Value {
			return conf.Key, nil
		}
	}
	return key, fmt.Errorf("not find value %s", value)
}

// 添加配置
func (c *HostsSetting) Add(key, value string) (err error) {
	conf := HostsText{Key: key, Value: value}
	c.Conf = append(c.Conf, conf)
	return nil
}

// 更新配置
func (c *HostsSetting) Update(key, value string) (err error) {
	for _, conf := range c.Conf {

		// 判断是否存在
		if value == conf.Value && key == conf.Key {
			return nil
		}
	}

	// 没有找到，则添加配置
	return c.Add(key, value)
}

// 删除配置
func (c *HostsSetting) Remove(key, value string) (err error) {
	for i := len(c.Conf) - 1; i >= 0; i-- {
		if c.Conf[i].Key == key && c.Conf[i].Value == value {
			c.Conf = append(c.Conf[:i], c.Conf[i+1:]...)
		}
	}
	return
}

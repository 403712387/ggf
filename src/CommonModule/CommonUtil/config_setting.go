package CommonUtil

import (
	"fmt"
	"io/ioutil"
	"strings"
)

// 解析配置文件的方法（主要是IP的配置文件）
// 类型（是注释还是配置)
type TextType int32

const (
	Type_ConfSetting = iota // 配置
	Type_Annotation         // 注释
)

// 一个配置项
type confText struct {
	TextType TextType //	类型
	Key      string   // 键
	Value    string   // 值
}

// 配置
type ConfSetting struct {
	FileName string // 配置文件的名称
	Conf     []confText
}

// 解析配置文件
func (c *ConfSetting) ParseConfigFile() (err error) {
	bytes, err := ioutil.ReadFile(c.FileName)
	if err != nil {
		return err
	}

	data := string(bytes)
	return c.Parse(&data)
}

// 解析配置
func (c *ConfSetting) Parse(data *string) (err error) {
	// 内容为空
	if data == nil || len(*data) <= 0 {
		return nil
	}

	// 获取每一行的配置
	lines := strings.Split(*data, "\n")
	for _, line := range lines {
		line = strings.Trim(line, " ")

		if strings.HasPrefix(line, "#") { // 本行为注释
			conf := confText{TextType: Type_Annotation, Key: line}
			c.Conf = append(c.Conf, conf)
		} else { // 本行为配置

			index := strings.Index(line, "=")

			// 不合法的配置项，丢弃
			if index <= 0 {
				continue
			}

			conf := confText{TextType: Type_ConfSetting, Key: strings.Trim(line[:index], " "), Value: strings.Trim(line[index+1:], " ")}
			c.Conf = append(c.Conf, conf)
		}
	}
	return nil
}

// 配置项转成string
func (c *ConfSetting) String() (result string) {
	breakline := "\n"
	for _, conf := range c.Conf {

		if conf.TextType == Type_Annotation { // 本行为注释行
			result += conf.Key + breakline
		} else { // 本行为配置行
			conf.Value = strings.Trim(conf.Value, " ")

			// 如果value为空，则删除本行
			if len(conf.Value) <= 0 {
				continue
			}

			result += conf.Key + "=" + conf.Value + breakline
		}
	}
	return
}

// 保存到文件
func (c *ConfSetting) Save() error {
	data := c.String()
	return ioutil.WriteFile(c.FileName, []byte(data), 0644)
}

// 是否存在对应的key
func (c *ConfSetting) ExistKey(key string) (result bool) {
	for _, conf := range c.Conf {

		// 忽略注释项
		if conf.TextType == Type_Annotation {
			continue
		}

		if key == conf.Key {
			return true
		}
	}
	return false
}

// 是否存在对应的配置
func (c *ConfSetting) ExistConf(key, value string) (result bool) {
	for _, conf := range c.Conf {

		// 忽略注释项
		if conf.TextType == Type_Annotation {
			continue
		}

		if key == conf.Key && value == conf.Value {
			return true
		}
	}
	return false
}

// 根据key获取value
func (c *ConfSetting) Value(key string) (value string, err error) {
	for _, conf := range c.Conf {
		if conf.TextType == Type_Annotation {
			continue
		}

		if key == conf.Key {
			return conf.Value, nil
		}
	}
	return value, fmt.Errorf("not find key %s", value)
}

// 添加配置
func (c *ConfSetting) Add(key, value string) (err error) {
	conf := confText{TextType: Type_ConfSetting, Key: key, Value: value}
	c.Conf = append(c.Conf, conf)
	return nil
}

// 更新配置
func (c *ConfSetting) Update(key, value string) (err error) {
	for i, conf := range c.Conf {
		if conf.TextType == Type_Annotation {
			continue
		}

		// 更新配置
		if key == conf.Key {

			// 因为go语言是值传递，所以必须要这样更新
			c.Conf[i] = confText{Key: key, Value: value}
			return nil
		}
	}

	// 没有找到，则添加配置
	return c.Add(key, value)
}

// 删除配置
func (c *ConfSetting) Remove(key string) (err error) {
	for i := len(c.Conf) - 1; i >= 0; i-- {
		if c.Conf[i].Key == key {
			c.Conf = append(c.Conf[:i], c.Conf[i+1:]...)
		}
	}
	return
}

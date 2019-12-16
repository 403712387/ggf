package HttpHelper

import (
	"CommonModule"
	"crypto/md5"
	"encoding/json"
	"fmt"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"os"
	"path"
	"strconv"
	"sync"
	"time"
)

type timeInfo struct {
	Ip     string `json:"ip"`
	Create string `json:"create_time"` // 创建时间
	Update string `json:"update_time"` // 更新时间
}

func (t timeInfo) String() (result string) {
	return fmt.Sprintf("create time;%s, update time:%s", t.Create, t.Update)
}

// 管理http的token
type TokenHelper struct {
	Token       map[string]timeInfo `json:"token"`
	tokenLock   sync.RWMutex
	SessionFile string // 保存token的文件
}

// 创建token
func InitTokenHelper() *TokenHelper {
	result := &TokenHelper{Token: make(map[string]timeInfo)}
	result.SessionFile = "./database/session"

	// 创建目录
	if !common.IsExist(path.Dir(result.SessionFile)) {
		os.MkdirAll(path.Dir(result.SessionFile), os.ModeDir)
	}

	// 加载tokan
	result.loadToken()

	// 定期检查token
	go result.checkTokenLoop()

	return result
}

// 从文件中加载token
func (t *TokenHelper) loadToken() (err error) {
	if !common.IsExist(t.SessionFile) {
		err = fmt.Errorf("not find session file %s", t.SessionFile)
		return
	}

	t.tokenLock.Lock()
	defer t.tokenLock.Unlock()
	data, err := ioutil.ReadFile(t.SessionFile)
	err = json.Unmarshal(data, &t.Token)
	return
}

// 把token保存到文件
func (t *TokenHelper) saveToken() (err error) {
	data, err := json.MarshalIndent(t.Token, "", " ")
	if err != nil {
		return
	}

	err = ioutil.WriteFile(t.SessionFile, data, 0666)
	return
}

// 创建一个新的token
func (t *TokenHelper) CreateToken(ip string) (token string) {

	// 生成token
	now := strconv.FormatInt(time.Now().UnixNano(), 10)
	token = fmt.Sprintf("%x", md5.Sum([]byte(now)))

	// 保存到map中
	t.tokenLock.Lock()
	defer t.tokenLock.Unlock()
	t.Token[token] = timeInfo{Create: time.Now().Format("2006-01-02 15:04:05"), Update: time.Now().Format("2006-01-02 15:04:05"), Ip: ip}

	// 保存到文件
	t.saveToken()
	return
}

// 是否存在token
func (t *TokenHelper) IsExist(token string) bool {
	t.tokenLock.Lock()
	defer t.tokenLock.Unlock()

	// 更新token中的update时间
	tm, ok := t.Token[token]
	if ok {
		tm.Update = time.Now().Format("2006-01-02 15:04:05")
		t.Token[token] = tm
	}
	return ok
}

// 删除token
func (t *TokenHelper) RemoveToken(token string) (err error) {
	if exist := t.IsExist(token); !exist {
		return fmt.Errorf("not find token %s", token)
	}

	t.tokenLock.Lock()
	defer t.tokenLock.Unlock()
	delete(t.Token, token)

	// 保存到文件
	t.saveToken()
	return
}

// 定期检查token是否有效
func (t *TokenHelper) checkTokenLoop() {
	logrus.Infof("begin check token loop")

	for {

		// 休眠五分钟
		time.Sleep(5 * time.Minute)

		// 检查token是否超过一天
		tokens := t.expireToken()
		for _, token := range tokens {
			logrus.Infof("token expire, remove token %s", token)
			t.RemoveToken(token)
		}

		// 保存到文件
		t.tokenLock.Lock()
		t.tokenLock.Unlock()
		t.saveToken()
	}

	logrus.Infof("end check token loop")
}

// 获取所有的过期tokan
func (t *TokenHelper) expireToken() (tokens []string) {
	t.tokenLock.Lock()
	defer t.tokenLock.Unlock()

	now := time.Now()
	for k, v := range t.Token {
		update, err := time.Parse("2006-01-02 15:04:05", v.Update)
		if err != nil {
			logrus.Error("parse update time of token fail, error:" + err.Error())
			continue
		}

		// 如果有一天没有使用token,则认为token失效
		if now.Sub(update) > 24*time.Hour {
			tokens = append(tokens, k)
		}
	}
	return
}

func (t *TokenHelper) String() (result string) {
	for k, v := range t.Token {
		result += fmt.Sprintf("token:%s, time:%s,", k, v.String())
	}
	return
}

func (t *TokenHelper) Json() (result []byte) {

	result, _ = json.MarshalIndent(t, "", " ")
	return
}

package ResourceModule

import (
	common "CommonModule"
	"errors"
	"fmt"
	"os"
	"strings"
)

func updateMapMode(mode int)(err error){
	_, err = common.CommondResult("ls /mars/web")
	if err != nil{
		err = errors.New("web not exist")
		return
	}
	if mode == 2 {
		resultMapMode,_ := common.CommondResult("cat /mars/web/webconfig.js |grep 'window.IS_OFFLINE_MAP' |awk -F'=' '{print $NF}' |awk -F';' '{print $1}'")
		resultMapModeSlice := strings.Split(strings.Trim(resultMapMode, "\n"), "\n")
		_,err = common.CommondResult(fmt.Sprintf(`sed -i "s/window.IS_OFFLINE_MAP =%s/window.IS_OFFLINE_MAP = %s/g" /mars/web/webconfig.js`, resultMapModeSlice[0], "true"))
		if err != nil {
			return
		}
		//把压缩包解压到对应的目录下
		if common.IsExist("/mars/web/amap") {
			os.RemoveAll("/mars/web/amap")
		}
		_, err = common.CommondResult("unzip -d /mars/web/ amap.zip")
		if err != nil {
			return
		}
		os.Remove("amap.zip")
	}
	return
}

func getMapMode()(mode common.MapModeInfo, err error){
	_, err1 := common.CommondResult("ls /mars/web")
	if err1 != nil{
		err = errors.New("web not exist")
		return
	}
	resultMapMode,_ := common.CommondResult("cat /mars/web/webconfig.js |grep 'window.IS_OFFLINE_MAP' |awk -F'=' '{print $NF}' |awk -F';' '{print $1}'")
	resultMapModeSlice := strings.Split(strings.Trim(resultMapMode, "\n"), "\n")
	if resultMapModeSlice[0] == " true" {
		mode.Status = 2
	}
	if resultMapModeSlice[0] ==  " false" {
		mode.Status = 1
	}
	return
}

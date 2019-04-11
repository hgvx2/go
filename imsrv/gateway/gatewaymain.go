package main

import (
	"fmt"
	"tools/configFile/iniConfig"
)

func main()  {

	/*
	args := os.Args
	if nil == args || len(args) != 3 {
		fmt.Println("参数错误", args)
		return
	}
	srvIndex := args[1]
	configPath := args[2]
	*/
	configPath := "E:\\GoWork\\im\\src\\imsrv\\config.ini"

	pConfig := iniConfig.NewIniConfig()
	if err := pConfig.AnalyzeConfigFile(configPath); nil != err {
		fmt.Println(err, "配置文件解析错误 path=", configPath)
		return
	}


}

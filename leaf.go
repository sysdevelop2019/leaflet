package leaf

import (
	"fmt"
	"gitee.com/aarlin/leaflet/cluster"
	"gitee.com/aarlin/leaflet/conf"
	"gitee.com/aarlin/leaflet/console"
	"gitee.com/aarlin/leaflet/log"
	"gitee.com/aarlin/leaflet/module"
	"os"
	"os/signal"
	"runtime/debug"
	"time"
)

func Run(mods ...module.Module) {
	defer TryE()
	// logger
	if conf.LogLevel != "" {
		logger, err := log.New(conf.LogLevel, conf.LogPath,conf.LogNamePrefix,conf.LogKeepHour, conf.LogFlag)
		if err != nil {
			panic(err)
		}
		log.Export(logger)
		defer logger.Close()
	}

	log.Release("Leaf starting up v", version)

	// module
	for i := 0; i < len(mods); i++ {
		module.Register(mods[i])
	}
	module.Init()

	// cluster
	cluster.Init()

	// console
	console.Init()

	// close
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, os.Kill)
	sig := <-c
	log.Release("Leaf closing down (signal: %v)", sig)
	console.Destroy()
	cluster.Destroy()
	module.Destroy()
}

func TryE() {
	errs := recover()
	if errs == nil {
		return
	}
	exeName := os.Args[0] //获取程序名称

	now := time.Now()  //获取当前时间
	pid := os.Getpid() //获取进程ID

	time_str := now.Format("20060102150405")                          //设定时间格式
	fname := fmt.Sprintf("%s-%d-%s-dump.log", exeName, pid, time_str) //保存错误信息文件名:程序名-进程ID-当前时间（年月日时分秒）
	fmt.Println("dump to file ", fname)

	f, err := os.Create(fname)
	if err != nil {
		return
	}
	defer f.Close()

	f.WriteString(fmt.Sprintf("%v\r\n", errs)) //输出panic信息
	f.WriteString("========\r\n")

	f.WriteString(string(debug.Stack())) //输出堆栈信息
}

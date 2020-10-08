package log

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
	"strings"
	"time"
)

// levels
const (
	debugLevel   = 0
	releaseLevel = 1
	errorLevel   = 2
	fatalLevel   = 3
)

const (
	printDebugLevel   = "[debug  ] "
	printReleaseLevel = "[release] "
	printErrorLevel   = "[error  ] "
	printFatalLevel   = "[fatal  ] "
)

const fileMaxSize  = 1024 * 1024 * 500 //日记最大的大小 默认500M
const checkNewFileDur  = time.Minute * 1  //每分钟检查一次

type Logger struct {
	level      int
	baseLogger *log.Logger
	baseFile   *os.File
}

var filePath string
var fileBaseName string
var fileIndex int
var curFile *os.File
var fileKeepHour time.Duration //日记保留的时间（小时）
var gLogger, _ = New("debug", "","",120, log.LstdFlags)

func getFileName() string  {
	fileIndex++
	filename := fileBaseName + fmt.Sprintf("_%02d.log",fileIndex)
	return path.Join(filePath,filename)
}

func New(strLevel string, pathname string,fileNamePrefix string,keepHour int, flag int) (*Logger, error) {
	// level
	var level int
	switch strings.ToLower(strLevel) {
	case "debug":
		level = debugLevel
	case "release":
		level = releaseLevel
	case "error":
		level = errorLevel
	case "fatal":
		level = fatalLevel
	default:
		return nil, errors.New("unknown level: " + strLevel)
	}

	if keepHour == 0 {
		keepHour = 120
	}

	// logger
	var baseLogger *log.Logger
	if pathname != "" {
		now := time.Now()
		fileKeepHour = time.Hour * time.Duration(keepHour)
		filePath = pathname
		if fileNamePrefix == "" {
			fileBaseName = fmt.Sprintf("%d%02d%02d_%02d%02d%02d",
				now.Year(),
				now.Month(),
				now.Day(),
				now.Hour(),
				now.Minute(),
				now.Second())
		}else{
			fileBaseName = fileNamePrefix + "_" + fmt.Sprintf("%d%02d%02d_%02d%02d%02d",
				now.Year(),
				now.Month(),
				now.Day(),
				now.Hour(),
				now.Minute(),
				now.Second())
		}


		os.MkdirAll(pathname,0777)

		file, err := os.Create(getFileName())
		if err != nil {
			return nil, err
		}

		baseLogger = log.New(file, "", flag)
		curFile = file
	} else {
		baseLogger = log.New(os.Stdout, "", flag)
	}

	// new
	logger := new(Logger)
	logger.level = level
	logger.baseLogger = baseLogger
	logger.baseFile = curFile
	return logger, nil
}

func checkNewFile()  {
	tick := time.Tick(checkNewFileDur)
	for{
		select {
		case <-tick:
			if curFile != nil {
				if info,err := os.Stat(curFile.Name()); err == nil {
					size := info.Size()
					if size > fileMaxSize {
						if file, err := os.Create(getFileName()); err == nil {
							gLogger.baseLogger.SetOutput(file)
							gLogger.baseFile = file
							curFile = file

							go clearOldFile()
						}
					}
				}
			}
		}
	}
}

//清理旧的文件
func clearOldFile()  {
	if files ,err :=ioutil.ReadDir(filePath); err == nil {
		now := time.Now()
		for _,info := range files {
			if now.Sub(info.ModTime()) > fileKeepHour {
				err := os.Remove(path.Join(filePath,info.Name()))
				if err != nil {
					Error("log","removeOldFile","err",err.Error())
				}
			}
		}
	}
}

// It's dangerous to call the method on logging
func (logger *Logger) Close() {
	if logger.baseFile != nil {
		logger.baseFile.Close()
	}

	logger.baseLogger = nil
	logger.baseFile = nil
}

func (logger *Logger) doPrint(level int, printLevel string, keyvals ...interface{}) {
	if level < logger.level {
		return
	}
	if logger.baseLogger == nil {
		panic("logger closed")
	}

	logs := " "
	alen := len(keyvals)
	for i := 0; i < alen;i ++ {
		if i % 2 == 0 {
			logs = logs + fmt.Sprint(keyvals[i]) + "="
		}else{
			logs = logs + fmt.Sprint(keyvals[i]) + ", "
		}
	}

	logger.baseLogger.Output(3, printLevel + logs)

	if level == fatalLevel {
		os.Exit(1)
	}
}

func (logger *Logger) doPrintf(level int, printLevel string, format string, a ...interface{}) {
	if level < logger.level {
		return
	}
	if logger.baseLogger == nil {
		panic("logger closed")
	}

	format = printLevel + format
	logger.baseLogger.Output(3, fmt.Sprintf(format, a...))

	if level == fatalLevel {
		os.Exit(1)
	}
}


func (logger *Logger) Debug( keyvals ...interface{}) {
	logger.doPrint(debugLevel, printDebugLevel, keyvals...)
}

func (logger *Logger) Release( keyvals ...interface{}) {
	logger.doPrint(releaseLevel, printReleaseLevel,  keyvals...)
}

func (logger *Logger) Error(keyvals ...interface{}) {
	logger.doPrint(errorLevel, printErrorLevel,  keyvals...)
}

func (logger *Logger) Fatal(keyvals ...interface{}) {
	logger.doPrint(fatalLevel, printFatalLevel,  keyvals...)
}



// It's dangerous to call the method on logging
func Export(logger *Logger) {
	if logger != nil {
		gLogger = logger
		go clearOldFile()
		go checkNewFile()
	}
}

func Debug( keyvals ...interface{}) {
	gLogger.doPrint(debugLevel, printDebugLevel,  keyvals...)
}

func Release( keyvals ...interface{}) {
	gLogger.doPrint(releaseLevel, printReleaseLevel,  keyvals...)
}

func Error( keyvals ...interface{}) {
	gLogger.doPrint(errorLevel, printErrorLevel, keyvals...)
}

func Fatal( keyvals ...interface{}) {
	gLogger.doPrint(fatalLevel, printFatalLevel,  keyvals...)
}

func DebugF(format string, a ...interface{}) {
	gLogger.doPrintf(debugLevel, printDebugLevel, format, a...)
}

func ReleaseF(format string, a ...interface{}) {
	gLogger.doPrintf(releaseLevel, printReleaseLevel, format, a...)
}

func ErrorF(format string, a ...interface{}) {
	gLogger.doPrintf(errorLevel, printErrorLevel, format, a...)
}

func FatalF(format string, a ...interface{}) {
	gLogger.doPrintf(fatalLevel, printFatalLevel, format, a...)
}

func Close() {
	gLogger.Close()
}

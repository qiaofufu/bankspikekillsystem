package global

import (
	"github.com/spf13/viper"
	"log"
	"os"
	"path"
	"time"
)

type Log struct {
	Info    *log.Logger
	Error   *log.Logger
	Warning *log.Logger
	Request *log.Logger
}

var initialized = false

func NewLog() *Log {
	receiver := Log{}
	if initialized {
		return nil
	}
	filepath := path.Join(viper.GetString("dir.log_dir"), time.Now().Format("20060102"))
	fileRelativePath := path.Join(filepath, "spike.log")
	// 文件目录是否存在， 不存在创建
	if _, err := os.Stat(filepath); err != nil {
		if !os.IsExist(err) {
			os.MkdirAll(filepath, os.ModePerm)
		}
	}
	file, err := os.OpenFile(fileRelativePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal(err)
	}

	receiver.Info = log.New(file, "[INFO]", log.Ldate|log.Ltime|log.Lshortfile)
	receiver.Warning = log.New(file, "[WARN]", log.Ldate|log.Ltime|log.Lshortfile)
	receiver.Error = log.New(file, "[ERROR]", log.Ldate|log.Ltime|log.Lshortfile)

	fileRelativePath = path.Join(filepath, "request.log")
	file, err = os.OpenFile(fileRelativePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal(err)
	}

	receiver.Request = log.New(file, "[REQUEST]", log.Ldate|log.Ltime|log.Lshortfile)

	initialized = true
	return &receiver
}

func (receiver Log) RequestInfo(info string) {
	receiver.Request.Println(info)
}

func (receiver Log) Println(v ...interface{}) {
	receiver.Info.Println(v)
}

func (receiver Log) Printf(format string, v ...interface{}) {
	receiver.Info.Printf(format, v)
}

func (receiver Log) Fatal(v ...interface{}) {
	receiver.Warning.Fatal(v)
}

func (receiver Log) ErrorPrintln(v ...interface{}) {
	receiver.Error.Println(v)
}

func (receiver Log) ErrorPrintf(format string, v ...interface{}) {
	receiver.Error.Printf(format, v)
}

// SpecialWarning 特别警告， 重大问题
func (receiver Log) SpecialWarning(v ...interface{}) {
	receiver.Warning.Println("[WARN]  ", v, "  [WARN]")
}

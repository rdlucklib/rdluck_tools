package log

import (
	"log"
	"os"
	"sync"
	"time"
)

const (
	defaultDateFormat = "20060102"
	defaultFileDir    = "./rdlucklog/"
)

type Log struct {
	date       string
	file       *os.File
	logger     *log.Logger
	dateFormat string
	fileDir    string
	mu         sync.Mutex
}

//第一个参数 文件名字（日期格式） 第二个参数 文件路径
func Init(v ...string) *Log {
	dateFormat := defaultDateFormat
	if len(v) > 0 {
		dateFormat = v[0]
	}
	fileDir := defaultFileDir
	if len(v) > 1 {
		fileDir = v[1]
	}
	date := time.Now().Format(dateFormat)
	file := createOrOpenFile(fileDir, date+".log")
	logger := log.New(file, "", log.LstdFlags|log.Lshortfile)
	return &Log{dateFormat: dateFormat, logger: logger, date: date, file: file, fileDir: fileDir}
}

func (l *Log) Println(v ...interface{}) {
	if l.date != time.Now().Format(l.dateFormat) {
		l.checkDate()
	}
	//输出到文件，多添加一个换行符
	l.logger.Println(v, "\n")
}

func (l *Log) checkDate() {
	l.mu.Lock()
	defer l.mu.Unlock()
	//如果相等 说明 日期已经处理过
	if l.date == time.Now().Format(l.dateFormat) {
		return
	}
	l.file.Close()
	l.date = time.Now().Format(l.dateFormat)
	l.file = createOrOpenFile(l.fileDir, l.date+".log")
	l.logger = log.New(l.file, "", log.LstdFlags|log.Lshortfile)
}

func createOrOpenFile(fileDir, path string) *os.File {
	os.MkdirAll(fileDir+time.Now().Format("/200601/"), os.ModePerm)
	path = fileDir + time.Now().Format("/200601/") + path
	fi, err := os.OpenFile(path, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0666)
	if err != nil {
		if !os.IsExist(err) {
			fi, err = os.Create(path)
		}
	}
	if err != nil {
		log.Println("logerror:" + err.Error())
	}
	return fi
}

package ezap

import (
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/weblazy/easy/filex"
)

func NewFileEzap(path ...string) *Ezap {
	var logPath string
	if len(path) == 1 {
		logPath = path[0]
	}
	if logPath == "" {
		logPath = filex.GetPath() + "/Runtime"
	}

	if !filex.CheckDir(logPath) {
		if err := filex.MkdirDir(logPath); err != nil {
			log.Printf("l.initZap(),err:%+v.\n", err)
		}
	}

	config := DefaultConfig()

	now := time.Now()
	filename := logPath + "/" + now.Format("2006-01-02") + ".log"
	logfile, err := os.OpenFile(filename, os.O_RDWR|os.O_APPEND, 0666)
	if err != nil {
		logfile, err = os.Create(filename)
		if err != nil {
			log.Println(err)
		}
	}

	config.ZapConfig.ErrorOutputPaths = []string{filename, "stderr"}
	config.ZapConfig.OutputPaths = []string{filename, "stdout"}
	l, err := config.ZapConfig.Build()
	if err != nil {
		log.Printf("l.initZap(),err:%+v.\n", err)
		return nil
	}
	e := &Ezap{
		Logfile: logfile,
		Logger:  l,
		Config:  config,
	}
	go e.updateLogFile(logPath)
	return e
}

// updateLogFile 检测是否跨天了,把记录记录到新的文件目录中
func (e *Ezap) updateLogFile(logPath string) {
	for {
		// 创建文件
		now := time.Now()
		next := now.Add(time.Hour * 24)
		next = time.Date(next.Year(), next.Month(), next.Day(), 0, 0, 0, 0, next.Location())
		t := time.NewTimer(next.Sub(now))
		<-t.C
		//以下为定时执行的操作
		if !filex.CheckDir(logPath) {
			if err := filex.MkdirDir(logPath); err != nil {
				log.Printf("l.initZap(),err:%+v.\n", err)
			}
		}
		filename := logPath + "/" + time.Now().Format("2006-01-02") + ".log"

		logfile, err := os.OpenFile(filename, os.O_RDWR|os.O_APPEND, 0666)
		if err != nil {
			logfile, err = os.Create(filename)
			if err != nil {
				log.Println(err)
			}
		}

		e.Config.ZapConfig.ErrorOutputPaths = []string{filename, "stderr"}
		e.Config.ZapConfig.OutputPaths = []string{filename, "stdout"}
		l, err := e.Config.ZapConfig.Build()
		if err != nil {
			log.Println(err)
			continue
		}
		lastLogfile := e.Logfile
		e.Logfile = logfile
		e.Logger = l
		lastLogfile.Close()
		go deleteLog(logPath, float64(e.Config.MaxAge))
		// 计算下一个零点

	}
}

// deleteLog 删除修改时间在saveDays天前的文件
func deleteLog(source string, saveDays float64) {
	err := filepath.Walk(source, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() && !strings.HasSuffix(info.Name(), ".log") {
			return nil
		}
		t := time.Since(info.ModTime()).Hours()
		if t >= (saveDays-1)*24 {
			e := os.Remove(path)
			if e != nil {
				log.Println(e)
			}
		}
		return err
	})
	if err != nil {
		return
	}
}

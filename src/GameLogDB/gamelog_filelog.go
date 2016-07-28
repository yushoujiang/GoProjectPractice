package GameLogDB

import (
	"errors"
	"fmt"
	"io"
	"log"
	"strings"
	"time"

	seelog "github.com/cihub/seelog"
)

var debugLogger seelog.LoggerInterface
var debugFileTime string = ""
var debugWrite io.WriteCloser = nil
var debugGenFileChan chan string

// var errorFileTime string = ""

// var errorWrite *seelog.rollingFileWriterSize

func init() {
	DisableLog()
	debugGenFileChan = make(chan string, 1)

	// 保证在主线程里处理
	go func() {
		for {
			select {
			case path := <-debugGenFileChan:
				{
					// log.Println("path:", path)
					if debugWrite != nil {
						err := debugWrite.Close()
						log.Println("writer_close_err:", err)
					}

					write, err := seelog.NewRollingFileWriterSize(path+"/debug.log", 2, path+"/ArchiveFile", 500*1024*1024, 1, 0, true) //500k

					if err != nil {
						panic(err)
					}

					SetLogWriter(write)
				}
			}
		}
	}()
}

// DisableLog disables all library log output.
func DisableLog() {
	debugLogger = seelog.Disabled
}

func RecordGameLog(value string) {
	n
	defer FlushLog()

	//根据时间判断是否要换一个文件夹写

	today := time.Now().Format("2006-01-02")
	needGenNewWrite := false
	if len(debugFileTime) <= 0 {
		needGenNewWrite = true
	} else if strings.EqualFold(debugFileTime, today) == false {
		needGenNewWrite = true
		log.Println("day change~")
	}
	if needGenNewWrite {
		debugFileTime = today
		//close old　write
		go func() {
			filePath := fmt.Sprintf("%s", today)
			debugGenFileChan <- filePath
			log.Println("5_value:", value)
			debugLogger.Debug(value)
		}()
	} else {
		log.Println("3_value:", value)
		debugLogger.Debug(value)
	}

}

// UseLogger uses a specified seelog.LoggerInterface to output library log.
// Use this func if you are using Seelog logging system in your app.
func UseLogger(newLogger seelog.LoggerInterface) {
	debugLogger = newLogger
}

// SetLogWriter uses a specified io.Writer to output library log.
// Use this func if you are not using Seelog logging system in your app.
func SetLogWriter(writer io.Writer) error {
	if writer == nil {
		return errors.New("Nil writer")
	}

	newLogger, err := seelog.LoggerFromWriterWithMinLevelAndFormat(writer, seelog.TraceLvl, "%Msg%n")
	if err != nil {
		return err
	}

	UseLogger(newLogger)
	return nil
}

// Call this before app shutdown
func FlushLog() {
	debugLogger.Flush()
}

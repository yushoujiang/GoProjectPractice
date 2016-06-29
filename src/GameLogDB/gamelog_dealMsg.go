package GameLogDB

import (
	"fmt"
	"log"
	"math/rand"
	"strings"
	"time"
)

const maxChanMsg int = 1024
const maxChanCount int = 3

var chanMsg []chan Message

type Message struct {
	chanid  byte      //编号id我看看的
	logData string    //数据消息
	channel chan byte //消息队列
}

func init() {
	chanMsg = make([]chan Message, maxChanCount)
	for i := 0; i < maxChanCount; i++ {
		chanMsg[i] = make(chan Message, maxChanMsg)
		go Run(chanMsg[i])
	}
}

const spliteString string = "\r\n\r\n"

//分割解析字符串
func spliteData(data []byte, connIndex uint16) []byte {

	str := string(data)

	log.Println("before spliteData:", str)

	index := strings.Index(str, spliteString)
	fmt.Println("find index:", index)
	if index == -1 { //不全等下次一起
		return data
	} else {
		message := string(data[:index])
		log.Println("after spliteData:", message)
		go Put(message, connIndex)

		//剩余内容返回

		leftStartIndex := index + len(spliteString)
		trueLen := len(data)
		log.Println("leftStartIndex:", leftStartIndex, ",", "trueLen:", trueLen)
		if leftStartIndex == trueLen {

		} else if leftStartIndex < trueLen {
			return data[leftStartIndex:trueLen]
		}
	}
	return make([]byte, 0)
}

//从socket里读取到的数据入队列
func Put(msg string, connIndex uint16) {

	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	var id byte = byte(r.Intn(maxChanCount))
	dataStruct := Message{chanid: id, logData: msg, channel: make(chan byte)}

	chanMsg[id] <- dataStruct

}

func Run(msg chan Message) {

	timeout := make(chan bool, 1)
	for {

		//做个超时机制
		go func() {
			time.Sleep(1e9 * 5) //等待1秒
			timeout <- true
		}()

		select {
		case message := <-msg:
			dispatchLog(message)
			// fmt.Println("入库", message.logData)
		case <-timeout:
			// fmt.Println("呵呵超时了")
		}
	}

}

func DataFromSocket(msg string) {

	dataStruct := Message{logData: msg, channel: make(chan byte)}

	go dispatchLog(dataStruct)
}

func dispatchLog(data Message) {
	message := data.logData
	finalStr := strings.TrimSpace(message)
	fmt.Printf("inDispatchLog[%d]=%s\r\n", data.chanid, finalStr)
	// if strings.HasPrefix(finalStr, "INFO") {

	// 	//先在这里处理消息吧
	// 	records := strings.Split(message, "\t")
	// 	if records == nil || len(records) < 6 {
	// 		return
	// 	}
	// 	//消息内容具体处理

	// 	if strings.Contains(records[5], "remain_record"); len(records) >= 7 {
	// 		// self.login_msg(record)
	// 	} else if strings.Contains(records[5], "recharge_successful"); len(records) >= 7 {
	// 		// self.recharge_msg(record)
	// 		// self.log_save(record)
	// 	} else if strings.Contains(records[5], "poll"); len(records) >= 7 {

	// 	} else {
	// 		log_save(records)
	// 	}

	// } else if strings.HasPrefix(finalStr, "ERROR") {

	// }

}

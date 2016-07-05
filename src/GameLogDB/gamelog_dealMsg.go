package GameLogDB

import (
	// "fmt"
	"log"
	"math/rand"
	"strings"
	"time"
)

const maxChanMsg int = 100
const maxChanCount int = 3

var chanMsg []chan Message

type Message struct {
	chanid   byte //编号id
	socketid uint
	logData  string    //数据消息
	channel  chan byte //消息队列
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
func spliteData(data []byte, connIndex uint) []byte {

	str := string(data)

	// log.Println("before spliteData:", str)

	finalStr := strings.SplitAfter(str, spliteString)
	//分割
	strLen := len(finalStr) - 1
	// log.Println("----------------------------strLen:", strLen)
	if strLen < 1 {
		return data
	}
	for i := 0; i < strLen-1; i++ {
		// log.Println("----------------------------put:", finalStr[i])
		go Put(finalStr[i], connIndex)
	}
	//最后一个判断下是否完整
	// log.Println("----------------------------put last:", finalStr[strLen-1])
	index := strings.Index(finalStr[strLen-1], spliteString)
	// log.Println("find index:", index)
	if index == -1 { //不全等下次一起
		return data
	} else {
		lastData := []byte(finalStr[strLen-1])
		message := string(lastData[:index])
		// log.Println("after spliteData:", message)
		go Put(message, connIndex)

		//剩余内容返回

		leftStartIndex := index + len(spliteString)
		trueLen := len(lastData)
		// log.Println("leftStartIndex:", leftStartIndex, ",", "trueLen:", trueLen)
		if leftStartIndex == trueLen {

		} else if leftStartIndex < trueLen {
			return lastData[leftStartIndex:trueLen]
		}
	}
	return make([]byte, 0)
}

// func spliteData(data []byte, connIndex uint) []byte {

// 	str := string(data)

// 	// log.Println("before spliteData:", str)

// 	index := strings.Index(str, spliteString)
// 	// log.Println("find index:", index)
// 	if index == -1 { //不全等下次一起
// 		return data
// 	} else {
// 		message := string(data[:index])
// 		// log.Println("after spliteData:", message)
// 		go Put(message, connIndex)

// 		//剩余内容返回

// 		leftStartIndex := index + len(spliteString)
// 		trueLen := len(data)
// 		log.Println("leftStartIndex:", leftStartIndex, ",", "trueLen:", trueLen)
// 		if leftStartIndex == trueLen {

// 		} else if leftStartIndex < trueLen {
// 			return data[leftStartIndex:trueLen]
// 		}
// 	}
// 	return make([]byte, 0)
// }

//从socket里读取到的数据入队列
func Put(msg string, connIndex uint) {

	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	var id byte = byte(r.Intn(maxChanCount))
	dataStruct := Message{chanid: id, logData: msg, channel: make(chan byte), socketid: connIndex}

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
			go dispatchLog(message)
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

//按照规则切割流里的数据内容
func dispatchLog(data Message) {
	message := data.logData
	finalStr := strings.TrimSpace(message)
	log.Printf("socket[%d],dispatchLog[%d]=%s\r\n", data.socketid, data.chanid, finalStr)

	if strings.HasPrefix(finalStr, "INFO") {

		//先在这里处理消息吧
		records := strings.Split(message, "\t")
		log.Printf("len(records) :%d\r\n", len(records))
		if records == nil || len(records) < 6 {
			return
		}
		log.Println("run here?   1")
		if strings.Contains(records[5], "remain_record"); len(records) >= 7 {
			// self.login_msg(record)
		} else if strings.Contains(records[5], "recharge_successful"); len(records) >= 7 {
			// self.recharge_msg(record)
			// self.log_save(record)
		} else if strings.Contains(records[5], "poll"); len(records) >= 7 {

		} else {
			log.Println("run here?   2")
			log_save(records)
		}

	} else if strings.HasPrefix(finalStr, "ERROR") {

	}

}

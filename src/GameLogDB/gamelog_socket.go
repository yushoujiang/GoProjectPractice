package GameLogDB

import (
	"fmt"
	"log"
	"net"
	// "os"
)

const socketPort uint16 = 4000

var gameSocketIndex uint

type GameConnect struct {
	index   uint     //索引
	connect net.Conn //实际socket连接句柄
}

func init() {
	gameSocketIndex = 0
}

//处理每一个连接
func handleConnection(con GameConnect) {
	log.Println(con.connect.RemoteAddr().String(), " connect successful,index:", con.index)

	defer func() {
		con.connect.Close()
		log.Println(con.connect.RemoteAddr().String(), " socket close")
	}()

	readerBuffer := make([]byte, 1024) //默认单条4K数据
	leftBuffer := make([]byte, 0)

	for {
		n, err := con.connect.Read(readerBuffer)
		// if err != nil {
		// 	log.Fatalln(con.connect.RemoteAddr().String(), " read error: ", err)

		if err != nil {
			//会丢弃该数据流里不完整的数据
			return
		}
		// log.Println("阻塞式的?")
		//截断字符内容
		leftBuffer = spliteData(append(leftBuffer, readerBuffer[:n]...), con.index)
	}

	//以下方式在需要拼接包的时候感觉会有问题，暂缓使用
	// readChannel := make(chan []byte)
	// writeChannel := make(chan []byte)

	// go func() {
	// 	for {
	// 		con.connect.Write(<-writeChannel)
	// 	}
	// }()
	// go func() {
	// 	for {
	// 		readChannel <- con.connect.Read()
	// 	}
	// }()
}

//启动日志服务器
func StartServer() {

	netListen, err := net.Listen("tcp", fmt.Sprintf(":%d", socketPort))

	checkError(err)

	defer netListen.Close()

	// fmt.Println("GameLog StartOver,WaitFor Client")
	log.Println("GameLog StartOver,WaitFor Client")

	for {

		connect, err := netListen.Accept()

		if err != nil {
			log.Fatalln(connect.RemoteAddr().String(), " connect error:", err)
			continue
		}

		con := GameConnect{index: gameSocketIndex, connect: connect}
		gameSocketIndex++

		go handleConnection(con)
	}
}

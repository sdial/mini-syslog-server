package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"strings"
	"time"
)

// main 是程序的入口点。
func main() {
	// 定义监听的地址和端口
	addr := "0.0.0.0:514"
	// 解析UDP地址
	udpAddr, err := net.ResolveUDPAddr("udp", addr)
	if err != nil {
		log.Fatal(err)
	}
	// 创建UDP连接，开始监听
	conn, err := net.ListenUDP("udp", udpAddr)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("正在监听", udpAddr)

	// 定义日志目录，并尝试创建
	logDir := "./logs"
	err = os.MkdirAll(logDir, os.ModePerm)
	if err != nil {
		log.Fatal(err)
	}

	// 循环监听，处理接收到的数据
	for {
		// 创建一个缓冲区来接收数据,当出现\n换行符时，表示一条日志消息结束，进入下一步处理
		buffer := make([]byte, 1024)
		// 从UDP连接中读取数据
		n, remoteAddr, err := conn.ReadFromUDP(buffer)
		if err != nil {
			log.Println(err)
			continue
		}

		// 并行处理接收到的数据
		go handleLog(buffer[:n], remoteAddr, logDir)
	}
}

// 定义一个函数，用于处理接收到的日志消息
func handleLog(buffer []byte, remoteAddr *net.UDPAddr, logDir string) {
	// 解析接收到的消息
	message := string(buffer)
	lines := strings.Split(message, "\n")

	// 处理每一行日志消息
	for _, line := range lines {
		// 获取当前日期
		now := time.Now()
		date := now.Format("2006-01-02")

		// 定义IP子目录，并尝试创建
		ipDir := fmt.Sprintf("%s/%s", logDir, remoteAddr.IP.String())
		err := os.MkdirAll(ipDir, os.ModePerm)
		if err != nil {
			log.Println(err)
			continue
		}

		// 定义日志文件名，并尝试打开文件
		fileName := fmt.Sprintf("%s/%s.log", ipDir, date)
		file, err := os.OpenFile(fileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			log.Println(err)
			continue
		}
		defer file.Close()

		// 将line前面加上当前时间，写入日志文件
		line = fmt.Sprintf("[%s] %s\r\n", now.Format("2006-01-02 15:04:05"), line)
		_, err = file.WriteString(line)
		if err != nil {
			log.Println(err)
			continue
		}

		// 将接收到的日志内容同时显示在屏幕上,并加上来源IP地址
		line = fmt.Sprintf("[%s]%s", remoteAddr.IP.String(), line)
		fmt.Println(line)
	}
}

package common

import (
	"fmt"
	"log"
	"os"
	"time"
)

func NewLog(logPath string) *os.File {
	//暂时没有使用， LogWriterLn自动会创建
	logFile, err := os.Create(logPath)
	if err != nil {
		panic(err)
	}
	defer logFile.Close()

	logs := log.New(logFile, "[Info]", log.Ldate|log.Lmicroseconds|log.Lshortfile)
	logs.Println("newLog file created!")

	return logFile
}

func LogWriterLn(logPath string, logLevel string, logDetail interface{}) {
	// 暂时实现Println
	logFile, err := os.OpenFile(logPath, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0644)
	if err != nil {
		fmt.Println(err)
		//logFile = NewLog(logPath)
	}
	logs := log.New(logFile, logLevel, log.Ldate|log.Lmicroseconds|log.Lshortfile)
	logs.Println(logDetail)
	fmt.Println(logLevel, time.Now(), logDetail)
}

func LogWriterF(logPath string, logLevel string, logDetail string) {
	// TODO
}

func LogWriterPanic(logPath string, logLevel string, logDetail string) {
	// TODO
}

func LogWriterFatal(logPath string, logLevel string, logDetail string) {
	// TODO
}

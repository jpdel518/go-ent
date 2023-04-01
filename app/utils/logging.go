package utils

import (
	"gopkg.in/natefinch/lumberjack.v2"
	"io"
	"log"
	"os"
)

func LoggingSettings(logFile string) {
	// logfile, err := os.OpenFile(logFile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666) // なかったら作成、Read, Write, Appendも許可
	// if err != nil {
	// 	log.Fatalln(err)
	// }
	lumberjackLogger := &lumberjack.Logger{
		Filename:   logFile,
		MaxSize:    1, // megabytes
		MaxBackups: 3,
		MaxAge:     28, // days
		Compress:   true,
		LocalTime:  true,
	}
	multiLogFile := io.MultiWriter(os.Stdout, lumberjackLogger) // logの出力先を標準出力とログファイルに指定
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)        // format指定
	log.SetOutput(multiLogFile)
}

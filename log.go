package main

import (
	"fmt"
	"log"
	"os"
)

// ANSI 颜色代码
const (
	ColorReset  = "\033[0m"
	ColorRed    = "\033[31m"
	ColorGreen  = "\033[32m"
	ColorYellow = "\033[33m"
	ColorBlue   = "\033[34m"
)

func logDebugF(format string, v ...interface{}) {
	if config.DebugLog {
		log.Printf(format, v...)
	}
}

func logDebug(v ...interface{}) {
	if config.DebugLog {
		log.Print(v...)
	}
}

func logInfoF(format string, v ...interface{}) {
	log.Printf(format+"\n", v...)
}

func logSuccessF(format string, v ...interface{}) {
	log.Printf(ColorGreen+format+ColorReset, v...)
}

func logWarningF(format string, v ...interface{}) {
	log.Printf(ColorYellow+format+ColorReset, v...)
}

func logErrorF(format string, v ...interface{}) {
	log.Printf(ColorRed+format+ColorReset, v...)
}

func logFatal(v ...interface{}) {
	log.Println(ColorRed, fmt.Sprint(v...), ColorReset)
	os.Exit(1)
}

func logFatalF(format string, v ...interface{}) {
	log.Printf(ColorRed+format+ColorReset, v...)
	os.Exit(1)
}

package log

import (
	"fmt"
	"log"
	"os"

	config "adb-backup/internal/config"
)

// ANSI 颜色代码
const (
	ColorReset  = "\033[0m"
	ColorRed    = "\033[31m"
	ColorGreen  = "\033[32m"
	ColorYellow = "\033[33m"
	ColorBlue   = "\033[34m"
)

func DebugF(format string, v ...interface{}) {
	if config.App.DebugLog {
		log.Printf(format, v...)
	}
}

func Debug(v ...interface{}) {
	if config.App.DebugLog {
		log.Print(v...)
	}
}

func InfoF(format string, v ...interface{}) {
	log.Printf(format+"\n", v...)
}

func SuccessF(format string, v ...interface{}) {
	log.Printf(ColorGreen+format+ColorReset, v...)
}

func WarningF(format string, v ...interface{}) {
	log.Printf(ColorYellow+format+ColorReset, v...)
}

func ErrorF(format string, v ...interface{}) {
	log.Printf(ColorRed+format+ColorReset, v...)
}

func Fatal(v ...interface{}) {
	log.Println(ColorRed, fmt.Sprint(v...), ColorReset)
	os.Exit(1)
}

func FatalF(format string, v ...interface{}) {
	log.Printf(ColorRed+format+ColorReset, v...)
	os.Exit(1)
}

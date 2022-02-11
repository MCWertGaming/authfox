package logHelper

import (
	"log"
	"os"
	"runtime"
)

func redBoldColor(message string) string {
	if _, exists := os.LookupEnv("DISABLE_COLOR"); exists {
		return message
	} else {
		return "\033[31m\033[1m" + message + "\033[m"
	}
}
func boldColor(message string) string {
	if _, exists := os.LookupEnv("DISABLE_COLOR"); exists {
		return message
	} else {
		return "\033[1m" + message + "\033[m"
	}
}
func ErrorFatal(name string, err error) {
	if err != nil {
		_, filename, line, _ := runtime.Caller(1)
		log.Fatalf("[%v] [%v] [%v:%v] %v", name, redBoldColor("FATAL"), filename, line, redBoldColor(err.Error()))
	}
}
func ErrorPanic(err error) {
	if err != nil {
		panic(err)
	}
}
func LogEvent(name string, message string) {
	_, filename, line, _ := runtime.Caller(1)
	log.Printf("[%v] [%v] [%v:%v] %v", name, boldColor("LOG"), filename, line, boldColor(message))
}
func LogError(name string, err error) {
	_, filename, line, _ := runtime.Caller(1)
	log.Printf("[%v] [%v] [%v:%v] %v", name, redBoldColor("ERROR"), filename, line, redBoldColor(err.Error()))
}

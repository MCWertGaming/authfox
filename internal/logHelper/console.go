package logHelper

import (
	"fmt"
	"log"
	"time"
)

// TODO: make pretty [name] [FATAL] 12.01.2022 MESSAGE
func ErrorFatal(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

// TODO: make pretty
func ErrorPanic(err error) {
	if err != nil {
		panic(err)
	}
}
func LogEvent(name string, message string) {
	fmt.Println("[" + name + "] " + time.Now().Format("2006/01/02-15:04:05: ") + message)
}
func LogError(name string, err error) {
	fmt.Println("[" + name + "] [ERROR] " + time.Now().Format("2006/01/02-15:04:05: ") + err.Error())
}

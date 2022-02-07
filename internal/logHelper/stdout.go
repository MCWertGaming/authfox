package loghelper

import (
	"log"
	"runtime"

	"github.com/fatih/color"
)

func ErrorFatal(name string, err error) {
	if err != nil {
		clr := color.New(color.FgRed, color.Bold).SprintFunc()
		_, filename, line, _ := runtime.Caller(1)
		log.Fatalf("[%v] [%v] [%v:%v] %v", name, clr("FATAL"), filename, line, clr(err.Error()))
	}
}

// TODO: make pretty
func ErrorPanic(err error) {
	if err != nil {
		panic(err)
	}
}
func LogEvent(name string, message string) {
	clr := color.New(color.Bold).SprintFunc()
	_, filename, line, _ := runtime.Caller(1)
	log.Printf("[%v] [LOG] [%v:%v] %v", name, filename, line, clr(message))
}
func LogError(name string, err error) {
	clr := color.New(color.FgRed).SprintFunc()
	_, filename, line, _ := runtime.Caller(1)
	log.Printf("[%v] [%v] [%v:%v] %v", name, clr("ERROR"), filename, line, clr(err.Error()))
}

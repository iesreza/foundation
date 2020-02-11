package cmd

import (
	"fmt"
	"reflect"
)

const (
	infoColor    = "\033[1;34m%s\033[0m"
	noticeColor  = "\033[1;36m%s\033[0m"
	warningColor = "\033[1;33m%s\033[0m"
	errorColor   = "\033[1;31m%s\033[0m"
	debugColor   = "\033[0;36m%s\033[0m"
)

func Sprint(a ...interface{}) string {
	if len(a) == 0 {
		return ""
	}

	if reflect.TypeOf(a[0]).String() == "string" {
		if len(a) > 1 {
			return fmt.Sprintf(a[0].(string), a[1:]...)
		} else {
			return a[0].(string)
		}

	}

	return fmt.Sprint(a...)

}

func Info(a ...interface{}) {
	fmt.Printf(infoColor, Sprint(a...)+"\n")
}

func Error(a ...interface{}) {
	fmt.Printf(errorColor, Sprint(a...)+"\n")
}

func Notice(a ...interface{}) {
	fmt.Printf(noticeColor, Sprint(a...)+"\n")
}

func Warning(a ...interface{}) {
	fmt.Printf(warningColor, Sprint(a...)+"\n")
}

func Debug(a ...interface{}) {
	fmt.Printf(debugColor, Sprint(a...)+"\n")
}

func Println(a ...interface{}) {
	fmt.Printf(debugColor, Sprint(a...)+"\n")
}

func Print(a ...interface{}) {
	fmt.Printf(debugColor, Sprint(a...))
}

package system

import "fmt"

func Critical(err error,args ...interface{})  {
	fmt.Println(err.Error())
}

func Error(err error,args ...interface{})  {
	fmt.Println(err.Error())
}

func Log()  {
	
}
package utils

import (
	"fmt"
	"strings"
)

func PrintCallErr(functionName string, callName string, err error) {
	fmt.Println(functionName, callName+"返回失败", err)
}

func PrintErr(functionName string, msg ...interface{}) {
	fmt.Println(functionName, msg)
}

func IsStringEmpty(str string) bool {
	return strings.Trim(str, " ") == ""
}

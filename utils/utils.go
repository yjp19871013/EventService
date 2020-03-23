package utils

import (
	"fmt"
	"io/ioutil"
	"os"
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

func PathExists(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		return true
	}

	if os.IsNotExist(err) {
		return false
	}

	return false
}

func GetDirFiles(dir string) ([]string, error) {
	dirList, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	filesRet := make([]string, 0)

	for _, file := range dirList {
		if file.IsDir() {
			files, err := GetDirFiles(dir + string(os.PathSeparator) + file.Name())
			if err != nil {
				return nil, err
			}

			filesRet = append(filesRet, files...)
		} else {
			filesRet = append(filesRet, dir+string(os.PathSeparator)+file.Name())
		}
	}

	return filesRet, nil
}

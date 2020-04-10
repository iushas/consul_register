package common

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
)

func FileRead(name string) []string {
	//bufio, read each line ,and return slice

	var str []string
	if fileObj, err := os.Open(name); err == nil {
		defer fileObj.Close()
		reader := bufio.NewReader(fileObj)
		for {
			if strLine, _, err := reader.ReadLine(); err == nil {
				str = append(str, string(strLine))
			} else {
				break
			}
		}
	} else {
		fmt.Println(err)
	}

	return str
}

func FileReadAll(filePath string) string {
	//ioutil, read all, and return all

	str, err := ioutil.ReadFile(filePath)
	if err != nil {
		LogWriterLn(LogPath, ERROR, "Read file failed: "+err.Error())
		return ""
		panic(err)
	}
	return string(str)
}

package kcli

import "strings"

const (
	sep string = "/"
)

func getCurrDir(path string) string {
	pathList := strings.Split(path, sep)
	len := len(pathList)
	if len == 0 {
		return ""
	}
	return pathList[len-1]
}

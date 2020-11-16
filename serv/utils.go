package serv

import (
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

//IsBlankStr ...
func IsBlankStr(value string) bool {
	value = strings.TrimSpace(value)
	con := []byte(value)
	if len(con) > 0 && con[0] != 0 {
		return false
	}
	return true
}

//IntegerToString ...
func IntegerToString(value int) string {
	return strconv.Itoa(value)
}

//StringToInteger ...
func StringToInteger(value string) int {
	res, _ := strconv.Atoi(value)
	return res
}

//ListFiles ...
func ListFiles(basePath string) []string {
	paths := make([]string, 0)
	filepath.Walk(basePath, func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			paths = append(paths, path)
		}
		return nil
	})
	return paths
}

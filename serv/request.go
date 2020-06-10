package serv

import (
	"io"
	"net"
	"strings"
)

//Request ...
type Request struct {
	conn       net.Conn
	protocal   string
	uri        string
	method     string
	headers    map[string]string
	body       io.Reader
	remoteAddr string
	attributes map[string]string
}

//判断uri指向的路径 是否是静态文件
func (req *Request) isStaticFile() (bool, string) {
	flag := false
	suffix := ""
	conf := GetConfig()
	for k := range conf.contentTypeMap {
		if strings.HasSuffix(req.uri, k) {
			flag = true
			suffix = k
			break
		}
	}
	return flag, suffix
}

//用于判断是否连接是keep-alive
func (req *Request) isKeepAlive() bool {
	connection := strings.TrimSpace(req.headers["Connection"])
	if connection == "keep-alive" {
		return true
	}

	return false

}

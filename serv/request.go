package serv

import (
	"io"
	"net"
	"strings"
)

//Request ...
type Request struct {
	Conn       net.Conn
	Protocal   string
	URI        string
	Method     string
	Headers    map[string]string
	Body       io.Reader
	RemoteAddr string
	Attributes map[string]string
	serv       *HTTPServer
}

//IsStaticFile 判断uri指向的路径 是否是静态文件
func (req *Request) IsStaticFile() (bool, string) {
	flag := false
	suffix := ""
	for k := range req.serv.ContentTypeMap {
		if strings.HasSuffix(req.URI, k) {
			flag = true
			suffix = k
			break
		}
	}
	return flag, suffix
}

//IsKeepAlive 用于判断是否连接是keep-alive
func (req *Request) IsKeepAlive() bool {
	connection := strings.TrimSpace(req.Headers["Connection"])
	if connection == "keep-alive" {
		return true
	}
	return false
}

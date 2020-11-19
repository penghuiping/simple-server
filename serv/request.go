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
//return 是否是静态文件，uri地址，静态文件后缀
func (req *Request) IsStaticFile() (bool, string, string) {
	flag := false
	suffix := ""
	questionMark := "?"
	uri := req.URI
	//先判断是否包含?,如果包含则先移除？后面的字符串信息
	if strings.Contains(uri, questionMark) {
		//包含?
		uri = uri[0:strings.Index(uri, questionMark)]
	}
	for k := range req.serv.ContentTypeMap {
		if strings.HasSuffix(uri, k) {
			flag = true
			suffix = k
			break
		}
	}
	return flag, uri, suffix
}

//IsKeepAlive 用于判断是否连接是keep-alive
func (req *Request) IsKeepAlive() bool {
	connection := strings.TrimSpace(req.Headers["Connection"])
	if connection == "keep-alive" {
		return true
	}
	return false
}

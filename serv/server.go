package serv

import (
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"strconv"
	"strings"
)

//StartServer ...
func StartServer(conf *Config) {
	listener, err := net.Listen("tcp", ":"+strconv.Itoa(conf.Port))
	if err != nil {
		log.Println(err)
		return
	}
	defer listener.Close()
	dispatcher := newDispatcher(config.GoroutineNum)
	dispatcher.run()
	for {
		conn, err1 := listener.Accept()
		if err1 != nil {
			log.Println(err1)
			continue
		}
		dispatcher.addJob(&conn)
	}
}

//AddRoute ...
func AddRoute(path string, handler func(*Request, *Response)) {
	router[path] = handler
}

func handle(req *Request, resp *Response) {
	//处理静态html文件
	result, suffix := isStaticFile(req)
	if result {
		handleStaticFile(req, resp, suffix)
		return
	}

	//处理自定义router
	handle1 := router[req.uri]
	if handle1 != nil {
		handle1(req, resp)
		resp.write()
		return
	}

	//404处理
	defaultHandle(req, resp)
	return
}

var router map[string]func(*Request, *Response) = make(map[string]func(*Request, *Response))

func defaultHandle(req *Request, resp *Response) {
	body := []byte("404 您访问的页面不存在\r\n")
	resp.code = StatusNotFound
	resp.codeMsg = "Not Found"
	resp.headers["Content-Type"] = "text/html;charset=utf-8"
	resp.headers["Connection"] = "keep-alive"
	resp.body = body
	resp.write()
}

var contentTypeMap map[string]string = initContentTypeMap()

func initContentTypeMap() map[string]string {
	map1 := make(map[string]string, 0)
	map1[".html"] = "text/html;charset=utf-8"
	map1[".css"] = "text/css;charset=utf-8"
	map1[".js"] = "application/x-javascript"
	map1[".gif"] = "image/gif"
	map1[".png"] = "image/png"
	map1[".woff"] = "application/x-font-woff"
	map1[".woff2"] = "application/x-font-woff"
	return map1
}

func isStaticFile(req *Request) (bool, string) {
	flag := false
	suffix := ""
	for k, _ := range contentTypeMap {
		if strings.HasSuffix(req.uri, k) {
			flag = true
			suffix = k
			break
		}
	}
	return flag, suffix
}

func handleStaticFile(req *Request, resp *Response, suffix string) {
	//列出html文件夹下所有的静态文件
	conf := GetConfig()
	paths := ListFiles(conf.HTMLPath)

	for _, path := range paths {
		if strings.HasSuffix(path, suffix) {
			path1 := strings.SplitN(path, "/", 2)[1]
			uri := req.uri
			if strings.HasPrefix(uri, "/") {
				uri = strings.Replace(uri, "/", "", 1)
			}
			if path1 == uri {
				content, err := ioutil.ReadFile(conf.HTMLPath + "/" + uri)
				if err != nil {
					log.Println(err)
					return
				}

				resp.code = StatusOK
				resp.codeMsg = "OK"
				resp.headers["Content-Type"] = contentTypeMap[suffix]
				resp.headers["Connection"] = "keep-alive"
				resp.headers["Content-Length"] = fmt.Sprintf("%d", len(content))
				resp.body = content
				resp.write()
				return
			}
		}
	}
}

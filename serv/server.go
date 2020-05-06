package serv

import (
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
	if strings.HasSuffix(req.uri, ".html") {
		conf := GetConfig()
		paths := ListFiles(conf.HTMLPath)
		for _, path := range paths {
			if strings.HasSuffix(path, ".html") {
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
					handleHTML(req, resp, content)
					return
				}
			}
		}
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

func handleHTML(req *Request, resp *Response, htmlContent []byte) {
	resp.code = StatusOK
	resp.codeMsg = "OK"
	resp.headers["Content-Type"] = "text/html;charset=utf-8"
	resp.headers["Connection"] = "keep-alive"
	resp.body = htmlContent
	resp.write()
}

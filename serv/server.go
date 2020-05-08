package serv

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"strconv"
	"strings"
)

//HTTPServer http服务器
type HTTPServer struct {
	router map[string]func(*Request, *Response)
}

//Start 启动服务器
func (serv *HTTPServer) Start() {
	serv.router = make(map[string]func(*Request, *Response))
	conf := GetConfig()
	listener, err := net.Listen("tcp", ":"+strconv.Itoa(conf.Port))
	if err != nil {
		log.Println(err)
		return
	}
	defer listener.Close()

	boss := &Boss{}
	boss.Start(config.GoroutineNum)

	boss.AddJobHandler("net.conn", func(job *Job) {
		conn := *(job.Content.(*net.Conn))
		defer func() {
			if err := recover(); err != nil {
				defer conn.Close()
				log.Println("panic异常是:" + fmt.Sprint(err))
			}
		}()

		total := make([]byte, 0)
		for {
			buf := make([]byte, 512)
			len, err3 := conn.Read(buf)
			if err3 != nil {
				log.Println(err3)
				conn.Close()
				break
			}
			if len > 0 {
				total = bytes.Join([][]byte{total, buf}, []byte{})
				if len < 512 {
					//一个请求完结
					req, _ := parseRequest(total)
					req.remoteAddr = conn.RemoteAddr().String()
					total = make([]byte, 0)

					resp := &Response{}
					resp.init(conn)
					serv.handle(req, resp)
				}
			}
		}

	})

	for {
		conn, err1 := listener.Accept()
		if err1 != nil {
			log.Println(err1)
			continue
		}

		job := &Job{}
		job.JobType = "net.conn"
		job.Content = &conn
		boss.AddJob(job)
	}
}

//AddRoute 加入路径路由
func (serv *HTTPServer) AddRoute(path string, handler func(*Request, *Response)) {
	serv.router[path] = handler
}

func (serv *HTTPServer) handle(req *Request, resp *Response) {
	//处理静态html文件
	result, suffix := req.isStaticFile()
	if result {
		handleStaticFile(req, resp, suffix)
		return
	}
	//处理自定义router
	handle1 := serv.router[req.uri]
	if handle1 != nil {
		handle1(req, resp)
		resp.write()
		return
	}

	//404处理
	defaultHandle(req, resp)
	return
}

//默认404处理
func defaultHandle(req *Request, resp *Response) {
	body := []byte("404 您访问的页面不存在\r\n")
	resp.code = StatusNotFound
	resp.codeMsg = "Not Found"
	resp.headers["Content-Type"] = "text/html;charset=utf-8"
	resp.headers["Connection"] = "keep-alive"
	resp.body = body
	resp.write()
}

//处理静态文件
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
				config := GetConfig()
				resp.headers["Content-Type"] = config.contentTypeMap[suffix]
				resp.headers["Connection"] = "keep-alive"
				resp.headers["Content-Length"] = fmt.Sprintf("%d", len(content))
				if suffix == ".woff2" {
					resp.headers["cache-control"] = "max-age=2592000"
				}
				resp.body = content
				resp.write()
				return
			}
		}
	}
}

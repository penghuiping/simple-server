package serv

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"log"
	"net"
	"os"
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

	addr := &net.TCPAddr{}
	addr.Port = conf.Port

	listener, err := net.ListenTCP("tcp", addr)
	if err != nil {
		log.Println("tcp listen出错:", err)
		return
	}
	defer listener.Close()

	boss := &Boss{}
	boss.Start(config.GoroutineNum)

	boss.AddJobHandler("net.conn", func(job *Job) {
		conn := *(job.Content.(*net.Conn))
		conn1 := conn.(*net.TCPConn)
		conn1.SetLinger(0)
		conn1.SetKeepAlive(true)
		conn1.SetNoDelay(true)
		defer func() {
			if err := recover(); err != nil {
				conn.Close()
				log.Println("job异常是:" + fmt.Sprint(err))
			}
		}()

		total := make([]byte, 0)
		for {
			buf := make([]byte, 512)
			len, err3 := conn.Read(buf)
			if err3 != nil {
				if err3 != io.EOF {
					log.Println("connection read出错:", err3)
				}
				conn.Close()
				break
			}
			if len > 0 {
				total = bytes.Join([][]byte{total, buf}, []byte{})
				if len < 512 {
					//一个请求完结
					req, _ := parseRequest(total)
					log.Println(req.uri)
					req.remoteAddr = conn.RemoteAddr().String()
					total = make([]byte, 0)
					resp := &Response{}
					resp.init(conn)
					serv.handle(req, resp)
					resp.writer.Flush()
				}
			}
		}

	})

	for {
		conn, err1 := listener.Accept()
		if err1 != nil {
			log.Println("server listen accept:", err1)
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
	resp.code = StatusNotFound
	resp.codeMsg = "Not Found"
	resp.headers["Content-Type"] = "text/html;charset=utf-8"
	if req.headers["Connection"] == "" {
		resp.headers["Connection"] = "close"
	} else {
		resp.headers["Connection"] = req.headers["Connection"]
	}
	resp.headers["Accept-Ranges"] = "bytes"
	bodyContent := "404 您访问的页面不存在\r\n"
	resp.bodySize = int64(len([]byte(bodyContent)))
	resp.body = bufio.NewReader(strings.NewReader("404 您访问的页面不存在\r\n"))
	resp.write()
}

//处理静态文件
func handleStaticFile(req *Request, resp *Response, suffix string) {
	//列出html文件夹下所有的静态文件
	conf := GetConfig()
	file, err := os.OpenFile(conf.HTMLPath+req.uri, os.O_RDONLY, 0666)
	defer file.Close()
	if err != nil {
		log.Println("打开静态文件出错:", err)
		defaultHandle(req, resp)
		return
	}

	resp.code = StatusOK
	resp.codeMsg = "OK"
	config := GetConfig()
	resp.headers["Content-Type"] = config.contentTypeMap[suffix]
	resp.headers["Server"] = "simple-server"
	if req.headers["Connection"] == "" {
		resp.headers["Connection"] = "close"
	} else {
		resp.headers["Connection"] = req.headers["Connection"]
	}

	resp.headers["Accept-Ranges"] = "bytes"
	// resp.headers["Date"] = time.Now().Format("Mon, 02 Jan 2006 15:04:05 GMT")
	fileInfo, err1 := os.Lstat(conf.HTMLPath + req.uri)
	resp.bodySize = fileInfo.Size()
	if err1 != nil {
		log.Println(err)
		return
	}
	if suffix == ".woff2" {
		resp.headers["cache-control"] = "max-age=2592000"
	}
	resp.body = bufio.NewReader(file)
	resp.write()
	return
}

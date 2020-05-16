package serv

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

//Interceptor 拦截器接口
type Interceptor interface {
	preHandle(req *Request)

	//返回值用于判断是否继续执行链路 true:继续执行
	handle(req *Request, resp *Response) bool

	postHandle(req *Request, resp *Response)
}

//FirstInterceptor ...
type FirstInterceptor struct {
}

func (h *FirstInterceptor) preHandle(req *Request) {
	conn := req.conn
	reader := bufio.NewReader(conn)
	//处理http request 第一行
	line, err1 := reader.ReadString('\n')
	if err1 != nil {
		panic(err1)
	}
	line = strings.TrimRight(line, "\r")
	firstLine := strings.Split(line, " ")
	req.method = firstLine[0]
	req.uri = firstLine[1]
	req.protocal = firstLine[2]

	//处理http request headers
	req.headers = make(map[string]string, 0)
	for {
		line2, err2 := reader.ReadString('\n')
		if err2 != nil {
			panic(err2)
		}
		line2 = strings.TrimRight(line2, "\r")
		if IsBlankStr(line2) {
			break
		}

		header := strings.Split(line2, ":")
		req.headers[header[0]] = header[1]
	}

	//处理http request body
	req.body = conn
}

func (h *FirstInterceptor) handle(req *Request, resp *Response) bool {
	return true
}

func (h *FirstInterceptor) postHandle(req *Request, res *Response) {
	file, ok := res.body.(*os.File)
	if ok {
		defer file.Close()
	}
	res.headers["Content-Length"] = fmt.Sprintf("%d", res.bodySize)
	if res.code == 0 {
		res.code = StatusOK
		res.codeMsg = "OK"
	}
	if req.isKeepAlive() {
		res.headers["Connection"] = "keep-alive"
	} else {
		res.headers["Connection"] = "close"
	}
	res.headers["Server"] = "simple-server"
	res.headers["Accept-Ranges"] = "bytes"
	writer := bufio.NewWriter(req.conn)
	writer.Write([]byte(fmt.Sprintf("HTTP/1.1 %d %s\r\n", res.code, res.codeMsg)))
	for k, v := range res.headers {
		writer.Write([]byte(fmt.Sprintf("%s: %s\r\n", k, v)))
	}
	writer.Write([]byte("\r\n"))
	buf := make([]byte, 4096)
	for {
		len, err := res.body.Read(buf)
		if err != nil {
			if err != io.EOF {
				log.Println("response流，写出出错:", err)
			}
			break
		}
		writer.Write(buf[0:len])
	}
	writer.Flush()
	if !req.isKeepAlive() {
		panic(io.EOF)
	}
}

//NotFoundInterceptor 404页面拦截器
type NotFoundInterceptor struct {
}

func (i *NotFoundInterceptor) preHandle(req *Request) {
}

func (i *NotFoundInterceptor) handle(req *Request, resp *Response) bool {
	resp.code = StatusNotFound
	resp.codeMsg = "Not Found"
	resp.headers["Content-Type"] = "text/html;charset=utf-8"
	bodyContent := "404 您访问的页面不存在\r\n"
	resp.bodySize = int64(len([]byte(bodyContent)))
	resp.body = bufio.NewReader(strings.NewReader("404 您访问的页面不存在\r\n"))
	return false
}

func (i *NotFoundInterceptor) postHandle(req *Request, resp *Response) {
}

//StaticFileInterceptor 静态文件处理拦截器
type StaticFileInterceptor struct {
}

func (f *StaticFileInterceptor) preHandle(req *Request) {
}

func (f *StaticFileInterceptor) handle(req *Request, resp *Response) bool {
	//判断uri是否是静态文件
	result, suffix := req.isStaticFile()

	if !result {
		return true
	}

	//html文件夹是否有符合uri的文件路径
	config := GetConfig()
	file, err := os.OpenFile(config.HTMLPath+req.uri, os.O_RDONLY, 0666)
	if err != nil {
		log.Println("打开静态文件出错:", err)
		return true
	}

	resp.code = StatusOK
	resp.codeMsg = "OK"
	resp.headers = make(map[string]string, 0)
	//处理静态html文件
	resp.headers["Content-Type"] = config.contentTypeMap[suffix]
	// resp.headers["Date"] = time.Now().Format("Mon, 02 Jan 2006 15:04:05 GMT")
	fileInfo, err1 := os.Lstat(config.HTMLPath + req.uri)
	resp.bodySize = fileInfo.Size()
	if err1 != nil {
		log.Println(err)
		return true
	}
	if suffix == ".woff2" {
		resp.headers["cache-control"] = "max-age=2592000"
	}
	resp.body = file
	return false
}

func (f *StaticFileInterceptor) postHandle(req *Request, resp *Response) {
}

//RouteIntercetor 自定义路由拦截器
type RouteIntercetor struct {
}

func (r *RouteIntercetor) preHandle(req *Request) {
}

func (r *RouteIntercetor) handler(req *Request, resp *Response) bool {
	//处理自定义router
	handle1 := GetConfig().routers[req.uri]
	if handle1 != nil {
		handle1(req, resp)
		return false
	}
	return true
}

func (r *RouteIntercetor) postHandle(req *Request, resp *Response) {
}

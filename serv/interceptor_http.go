package serv

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

//PreHTTPInterceptor ...
type PreHTTPInterceptor struct {
	Type  int8
	Order int
}

//Handle 返回值用于判断是否继续执行链路 true:继续执行
func (h *PreHTTPInterceptor) Handle(req *Request, resp *Response) bool {
	log.Println("进入PreHTTPInterceptor...")
	conn := req.Conn
	reader := bufio.NewReader(conn)
	//处理http request 第一行
	line, err1 := reader.ReadString('\n')
	if err1 != nil {
		panic(err1)
	}
	line = strings.TrimRight(line, "\r")
	firstLine := strings.Split(line, " ")
	req.Method = firstLine[0]
	req.URI = firstLine[1]
	req.Protocal = firstLine[2]

	//处理http request headers
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
		req.Headers[header[0]] = header[1]
	}

	//处理http request body
	req.Body = conn
	return true
}

//PostHTTPInterceptor ...
type PostHTTPInterceptor struct {
	Type  int8
	Order int
}

//Handle 返回值用于判断是否继续执行链路 true:继续执行
func (h *PostHTTPInterceptor) Handle(req *Request, res *Response) bool {
	log.Println("进入PostHTTPInterceptor...")
	file, ok := res.Body.(*os.File)
	if ok {
		//如果res.Body是os.File类型
		defer file.Close()
		fileInfo, err1 := os.Lstat(GetConfig().StaticFilePath + req.URI)
		res.BodySize = fileInfo.Size()
		if err1 != nil {
			panic(err1)
		}
	}

	reader, ok := res.Body.(*strings.Reader)
	if ok {
		//如果res.Body是strings.Reader类型
		res.BodySize = reader.Size()
	}

	if res.Code == 0 {
		res.Code = StatusOK
		res.CodeMsg = "OK"
	}
	if req.IsKeepAlive() {
		res.Headers["Connection"] = "keep-alive"
	} else {
		res.Headers["Connection"] = "close"
	}
	res.Headers["Server"] = "simple-server"
	res.Headers["Accept-Ranges"] = "bytes"
	res.Headers["Content-Length"] = fmt.Sprintf("%d", res.BodySize)

	writer := bufio.NewWriter(req.Conn)
	writer.Write([]byte(fmt.Sprintf("HTTP/1.1 %d %s\r\n", res.Code, res.CodeMsg)))
	for k, v := range res.Headers {
		writer.Write([]byte(fmt.Sprintf("%s: %s\r\n", k, v)))
	}
	writer.Write([]byte("\r\n"))
	buf := make([]byte, 4096)
	for {
		len, err := res.Body.Read(buf)
		if err != nil {
			if err != io.EOF {
				log.Println("response流，写出出错:", err)
			}
			break
		}
		writer.Write(buf[0:len])
	}
	writer.Flush()
	if !req.IsKeepAlive() {
		panic(io.EOF)
	}
	return true
}

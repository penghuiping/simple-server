package serv

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

//HTTPInterceptor ...
type HTTPInterceptor struct {
}

func (h *HTTPInterceptor) preHandle(req *Request) {
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

func (h *HTTPInterceptor) handle(req *Request, resp *Response) bool {
	return true
}

func (h *HTTPInterceptor) postHandle(req *Request, res *Response) {
	file, ok := res.body.(*os.File)
	if ok {
		defer file.Close()
		fileInfo, err1 := os.Lstat(config.StaticFilePath + req.uri)
		res.bodySize = fileInfo.Size()
		if err1 != nil {
			panic(err1)
		}
	}

	reader, ok := res.body.(*strings.Reader)
	if ok {
		res.bodySize = reader.Size()
	}

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
	res.headers["Content-Length"] = fmt.Sprintf("%d", res.bodySize)

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

//HTTPValidationInterceptor ...
type HTTPValidationInterceptor struct {
}

func (h *HTTPValidationInterceptor) preHandle(req *Request) {

}

func (h *HTTPValidationInterceptor) handle(req *Request, resp *Response) bool {
	return false
}

func (h *HTTPValidationInterceptor) postHandle(req *Request, resp *Response) {

}

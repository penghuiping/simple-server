package serv

import (
	"bufio"
	"strings"
)

//Interceptor 拦截器接口
type Interceptor interface {
	preHandle(req *Request)

	//返回值用于判断是否继续执行链路 true:继续执行
	handle(req *Request, resp *Response) bool

	postHandle(req *Request, resp *Response)
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
	resp.body = bufio.NewReader(strings.NewReader(bodyContent))
	return false
}

func (i *NotFoundInterceptor) postHandle(req *Request, resp *Response) {
}

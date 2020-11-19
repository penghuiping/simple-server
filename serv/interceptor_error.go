package serv

import (
	"log"
	"strings"
)

//NotFoundInterceptor 404页面拦截器
type NotFoundInterceptor struct {
}

//Handle 返回值用于判断是否继续执行链路 true:继续执行
func (i *NotFoundInterceptor) Handle(req *Request, resp *Response) bool {
	resp.Code = StatusNotFound
	resp.CodeMsg = "Not Found"
	resp.Headers["Content-Type"] = "text/html;charset=utf-8"
	resp.Body = strings.NewReader("404 您访问的页面不存在\r\n")
	resp.BodySize = int64(len([]byte("404 您访问的页面不存在\r\n")))
	return false
}

//InterceptorOrder 用于判断此拦截器的先后顺序，数字越小优先级越高,此拦截器就会被优先执行
func (i *NotFoundInterceptor) InterceptorOrder() int {
	return 100
}

//ServerErrorInterceptor 500页面拦截器
type ServerErrorInterceptor struct {
}

//Handle 返回值用于判断是否继续执行链路 true:继续执行
func (i *ServerErrorInterceptor) Handle(req *Request, resp *Response) bool {
	log.Println("进入ServerErrorInterceptor...")
	resp.Code = StatusInternalServerError
	resp.CodeMsg = "server error"
	resp.Headers["Content-Type"] = "text/html;charset=utf-8"
	bodyContent := "500 服务器错误，请联系相关技术人员解决\r\n"
	resp.BodySize = int64(len([]byte(bodyContent)))
	resp.Body = strings.NewReader(bodyContent)
	return false
}

//InterceptorOrder 用于判断此拦截器的先后顺序，数字越小优先级越高,此拦截器就会被优先执行
func (i *ServerErrorInterceptor) InterceptorOrder() int {
	return 100
}

package serv

import (
	"strings"
)

//NotFoundInterceptor 404页面拦截器
type NotFoundInterceptor struct {
	Type  int8
	Order int32
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

//InterceptorType 用于判断此拦截器的类型 0:前置拦截器 1:后置拦截器
func (i *NotFoundInterceptor) InterceptorType() uint8 {
	return PostIntercpetor
}

//InterceptorOrder 用于判断此拦截器的先后顺序，数字越小优先级越高,此拦截器就会被优先执行
func (i *NotFoundInterceptor) InterceptorOrder() uint32 {
	return 1
}

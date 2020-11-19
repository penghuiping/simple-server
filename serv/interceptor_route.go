package serv

import (
	"log"
)

//RouteInterceptor 自定义路由拦截器
type RouteInterceptor struct {
}

//Handle 返回值用于判断是否继续执行链路 true:继续执行
func (r *RouteInterceptor) Handle(req *Request, resp *Response) bool {
	log.Println("进入RouteInterceptor...")
	//处理自定义router
	handle1 := req.serv.Routers[req.URI]
	if handle1 != nil {
		handle1(req, resp)
		//本拦截器就可以处理，不需要继续执行
		return false
	}
	return true
}

//InterceptorOrder 用于判断此拦截器的先后顺序，数字越小优先级越高,此拦截器就会被优先执行
func (r *RouteInterceptor) InterceptorOrder() int {
	return 1
}

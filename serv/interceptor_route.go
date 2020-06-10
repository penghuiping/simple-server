package serv

//RouteIntercetor 自定义路由拦截器
type RouteIntercetor struct {
}

func (r *RouteIntercetor) preHandle(req *Request) {
}

func (r *RouteIntercetor) handle(req *Request, resp *Response) bool {
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

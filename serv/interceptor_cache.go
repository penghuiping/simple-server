package serv

//StaticsCacheInterceptor 用于实现http 缓存机制 Last-Modified/If-Modified-Since
type StaticsCacheInterceptor struct {
}

func (t *StaticsCacheInterceptor) Handle(req *Request, resp *Response) bool {
	return false
}

func (t *StaticsCacheInterceptor) InterceptorOrder() int {
	return 2
}

package serv

import (
	"fmt"
	"io"
	"log"
	"runtime/debug"
)

//Interceptor 拦截器接口
type Interceptor interface {

	//Handle 返回值用于判断是否继续执行链路 true:继续执行
	Handle(req *Request, resp *Response) bool

	//InterceptorOrder 用于判断此拦截器的先后顺序，数字越小优先级越高,此拦截器就会被优先执行
	InterceptorOrder() int
}

//InterceptorManager 拦截器管理器
type InterceptorManager struct {
	interceptors    []Interceptor
	serverErrorHTTP *ServerErrorInterceptor
	processor       *HTTPProcessor
}

//Init 初始化拦截器管理器
func (s *InterceptorManager) Init() {
	s.interceptors = make([]Interceptor, 0)
	s.Add(&RouteInterceptor{})
	s.Add(&StaticFileInterceptor{})
	s.Add(&NotFoundInterceptor{})
	s.processor = &HTTPProcessor{}
	s.serverErrorHTTP = &ServerErrorInterceptor{}
}

//Add 添加拦截器
func (s *InterceptorManager) Add(interceptor Interceptor) {
	s.interceptors = append(s.interceptors, interceptor)
}

//Run 运行拦截器
func (s *InterceptorManager) Run(req *Request, res *Response) {
	s.processor.Decode(req, res)
	for _, inter := range s.interceptors {
		result := inter.Handle(req, res)
		if !result {
			break
		}
	}
	s.processor.Encode(req, res)
}

//ServerErrorHandle 服务器错误处理
func (s *InterceptorManager) ServerErrorHandle(req *Request, res *Response) {
	if err := recover(); err != nil {
		if err != io.EOF {
			log.Println("异常是:" + fmt.Sprint(err) + "\n" + string(debug.Stack()))
			s.serverErrorHTTP.Handle(req, res)
			s.processor.Encode(req, res)
		}
		if req.Conn != nil {
			req.Conn.Close()
		}
	}
}

//InterceptorSlice ...
type InterceptorSlice []Interceptor

func (p InterceptorSlice) Len() int { return len(p) }
func (p InterceptorSlice) Less(i, j int) bool {
	return p[i].InterceptorOrder() < p[j].InterceptorOrder()
}
func (p InterceptorSlice) Swap(i, j int) { p[i], p[j] = p[j], p[i] }

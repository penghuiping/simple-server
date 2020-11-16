package serv

//Interceptor 拦截器接口
type Interceptor interface {

	//Handle 返回值用于判断是否继续执行链路 true:继续执行
	Handle(req *Request, resp *Response) bool

	//InterceptorType 用于判断此拦截器的类型 0:前置拦截器 1:后置拦截器
	InterceptorType() uint8

	//InterceptorOrder 用于判断此拦截器的先后顺序，数字越小优先级越高,此拦截器就会被优先执行
	InterceptorOrder() uint32
}

const (
	//PreIntercpetor 前置拦截器
	PreIntercpetor = 0
	//PostIntercpetor 后置拦截器
	PostIntercpetor = 1
)

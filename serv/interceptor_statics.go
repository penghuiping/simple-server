package serv

import (
	"log"
	"os"
)

//StaticFileInterceptor 静态文件处理拦截器
type StaticFileInterceptor struct {
}

//Handle 返回值用于判断是否继续执行链路 true:继续执行
func (f *StaticFileInterceptor) Handle(req *Request, resp *Response) bool {
	log.Println("进入StaticFileInterceptor...")
	//判断uri是否是静态文件
	result, suffix := req.IsStaticFile()

	if !result {
		//不是静态资源 本拦截器无法处理 直接跳过
		return true
	}

	//html文件夹是否有符合uri的文件路径
	file, err := os.OpenFile(req.serv.StaticFilePath+req.URI, os.O_RDONLY, 0666)
	if err != nil {
		panic(err)
	}

	resp.Code = StatusOK
	resp.CodeMsg = "OK"
	resp.Headers["Content-Type"] = req.serv.ContentTypeMap[suffix]
	//处理静态html文件
	if suffix != ".html" {
		resp.Headers["cache-control"] = "max-age=2592000"
	}
	resp.Body = file
	resp.BodySize = 0
	//本拦截器就可以处理，不需要继续执行
	return false
}

//InterceptorOrder 用于判断此拦截器的先后顺序，数字越小优先级越高,此拦截器就会被优先执行
func (f *StaticFileInterceptor) InterceptorOrder() int {
	return 2
}

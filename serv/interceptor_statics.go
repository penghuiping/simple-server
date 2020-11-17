package serv

import (
	"log"
	"os"
)

//StaticFileInterceptor 静态文件处理拦截器
type StaticFileInterceptor struct {
	Type  int8
	Order int32
}

//Handle 返回值用于判断是否继续执行链路 true:继续执行
func (f *StaticFileInterceptor) Handle(req *Request, resp *Response) bool {
	log.Println("进入StaticFileInterceptor...")
	//判断uri是否是静态文件
	result, suffix := req.IsStaticFile()

	if !result {
		return true
	}

	//html文件夹是否有符合uri的文件路径
	config := GetConfig()

	file, err := os.OpenFile(config.StaticFilePath+req.URI, os.O_RDONLY, 0666)
	if err != nil {
		if os.IsExist(err) {
			log.Println("打开静态文件出错:", err)
		}
		return true
	}
	resp.Code = StatusOK
	resp.CodeMsg = "OK"

	resp.Headers["Content-Type"] = config.ContentTypeMap[suffix]

	//处理静态html文件
	if suffix == ".woff2" {
		resp.Headers["cache-control"] = "max-age=2592000"
	}
	resp.Body = file
	resp.BodySize = 0
	return false
}

//InterceptorOrder 用于判断此拦截器的先后顺序，数字越小优先级越高,此拦截器就会被优先执行
func (f *StaticFileInterceptor) InterceptorOrder() int {
	return 2
}

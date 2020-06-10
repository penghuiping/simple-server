package serv

import (
	"log"
	"os"
)

//StaticFileInterceptor 静态文件处理拦截器
type StaticFileInterceptor struct {
}

func (f *StaticFileInterceptor) preHandle(req *Request) {
}

func (f *StaticFileInterceptor) handle(req *Request, resp *Response) bool {
	//判断uri是否是静态文件
	result, suffix := req.isStaticFile()

	if !result {
		return true
	}

	//html文件夹是否有符合uri的文件路径
	config := GetConfig()

	file, err := os.OpenFile(config.StaticFilePath+req.uri, os.O_RDONLY, 0666)
	if err != nil {
		if os.IsExist(err) {
			log.Println("打开静态文件出错:", err)
		}
		return true
	}

	resp.code = StatusOK
	resp.codeMsg = "OK"
	//处理静态html文件
	resp.headers["Content-Type"] = config.contentTypeMap[suffix]
	if suffix == ".woff2" {
		resp.headers["cache-control"] = "max-age=2592000"
	}
	resp.body = file
	return false
}

func (f *StaticFileInterceptor) postHandle(req *Request, resp *Response) {
}

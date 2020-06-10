package serv

import (
	"fmt"
	"io"
	"log"
	"net"
)

//HTTPServer http服务器
type HTTPServer struct {
}

//Init 初始化
func (serv *HTTPServer) Init(htmlPath string, goroutinNum int, port int) {
	conf := GetConfig()
	conf.StaticFilePath = htmlPath
	conf.Port = port
	conf.GoroutineNum = goroutinNum
	serv.AddInterceptor(&HTTPInterceptor{})
	serv.AddInterceptor(&RouteIntercetor{})
	serv.AddInterceptor(&StaticFileInterceptor{})
	serv.AddInterceptor(&NotFoundInterceptor{})
}

//Start 启动服务器
func (serv *HTTPServer) Start() {
	conf := GetConfig()
	addr := &net.TCPAddr{}
	addr.Port = conf.Port

	listener, err := net.ListenTCP("tcp", addr)
	if err != nil {
		log.Println("tcp listen出错:", err)
		return
	}
	defer listener.Close()

	boss := &Boss{}
	boss.Start(config.GoroutineNum)

	boss.AddJobHandler("net.conn", func(job *Job) {
		conn := *(job.Content.(*net.Conn))
		conn1 := conn.(*net.TCPConn)
		conn1.SetLinger(-1)
		conn1.SetKeepAlive(true)
		conn1.SetNoDelay(true)

		for {
			req := &Request{}
			req.conn = conn
			req.headers = make(map[string]string, 0)
			resp := &Response{}
			resp.headers = make(map[string]string, 0)

			defer func() {
				if err := recover(); err != nil {
					if err != io.EOF {
						log.Println("job异常是:" + fmt.Sprint(err))
						handleError(err, req, resp)
						httpInterceptor := &HTTPInterceptor{}
						httpInterceptor.postHandle(req, resp)
					}
					conn.Close()
				}
			}()

			interceptors := GetConfig().interceptors

			//preHandle
			for _, interceptor := range interceptors {
				interceptor.preHandle(req)
			}

			//handle
			for _, interceptor := range interceptors {
				result := interceptor.handle(req, resp)
				if !result {
					break
				}
			}

			//postHandle
			for _, interceptor := range interceptors {
				interceptor.postHandle(req, resp)
			}
		}
	})

	for {
		conn, err1 := listener.Accept()
		if err1 != nil {
			log.Println("server listen accept:", err1)
			continue
		}

		job := &Job{}
		job.JobType = "net.conn"
		job.Content = &conn
		boss.AddJob(job)
	}
}

//AddRoute 加入路径路由
func (serv *HTTPServer) AddRoute(path string, handler func(*Request, *Response)) {
	GetConfig().routers[path] = handler
}

//AddInterceptor 添加拦截器
func (serv *HTTPServer) AddInterceptor(interceptor Interceptor) {
	GetConfig().interceptors = append(GetConfig().interceptors, interceptor)
}

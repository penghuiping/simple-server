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
func (s *HTTPServer) Init(htmlPath string, goroutinNum int, port int) {
	conf := GetConfig()
	conf.StaticFilePath = htmlPath
	conf.Port = port
	conf.GoroutineNum = goroutinNum
	s.AddInterceptor(&PreHTTPInterceptor{})
	s.AddInterceptor(&PostHTTPInterceptor{})
	s.AddInterceptor(&RouteIntercetor{})
	s.AddInterceptor(&StaticFileInterceptor{})
	s.AddInterceptor(&NotFoundInterceptor{})
}

//Start 启动服务器
func (s *HTTPServer) Start() {
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
	boss.Start(conf.GoroutineNum)

	boss.AddJobHandler("net.conn", func(job *Job) {
		conn := *(job.Content.(*net.Conn))
		conn1 := conn.(*net.TCPConn)
		conn1.SetLinger(-1)
		conn1.SetKeepAlive(true)
		conn1.SetNoDelay(true)

		for {
			req := &Request{}
			req.Conn = conn
			req.Headers = make(map[string]string, 0)
			resp := &Response{}
			resp.Headers = make(map[string]string, 0)

			defer func() {
				if err := recover(); err != nil {
					if err != io.EOF {
						log.Println("job异常是:" + fmt.Sprint(err))
						handleError(err, req, resp)
						httpInterceptor := &PostHTTPInterceptor{}
						httpInterceptor.Handle(req, resp)
					}
					conn.Close()
				}
			}()

			interceptors := GetConfig().Interceptors
			//handle
			for _, interceptor := range interceptors {
				result := interceptor.Handle(req, resp)
				if !result {
					break
				}
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
func (s *HTTPServer) AddRoute(path string, handler func(*Request, *Response)) {
	GetConfig().Routers[path] = handler
}

//AddInterceptor 添加拦截器
func (s *HTTPServer) AddInterceptor(interceptor Interceptor) {
	GetConfig().Interceptors = append(GetConfig().Interceptors, interceptor)
}

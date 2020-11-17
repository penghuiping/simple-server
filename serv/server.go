package serv

import (
	"fmt"
	"io"
	"log"
	"net"
)

//HTTPServer http服务器
type HTTPServer struct {
	im   *InterceptorManager
	conf *Config
}

//Init 初始化
func (s *HTTPServer) Init(htmlPath string, goroutinNum int, port int) {
	conf := GetConfig()
	conf.StaticFilePath = htmlPath
	conf.Port = port
	conf.GoroutineNum = goroutinNum
	s.conf = conf
	s.im = &InterceptorManager{}
	s.im.Init()
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

	boss := s.initBossWorkers(conf.GoroutineNum)
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
	s.im.Add(interceptor)
}

//initBossWorkers 初始化工作线程
func (s *HTTPServer) initBossWorkers(workersNumber int) *Boss {
	boss := &Boss{}
	boss.Start(workersNumber)

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
						serverErrorInterceptor := &ServerErrorInterceptor{}
						serverErrorInterceptor.Handle(req, resp)
						httpInterceptor := &PostHTTPInterceptor{}
						httpInterceptor.Handle(req, resp)
					}
					conn.Close()
				}
			}()
			//运行拦截器
			s.im.Run(req, resp)
		}
	})
	return boss
}

package serv

import (
	"log"
	"net"
)

//HTTPServer http服务器
type HTTPServer struct {
	im             *InterceptorManager
	StaticFilePath string
	Port           int
	GoroutineNum   int
	ContentTypeMap map[string]string
	Routers        map[string]func(*Request, *Response)
}

//Init 初始化
func (s *HTTPServer) Init(htmlPath string, workers int, port int) {
	s.StaticFilePath = htmlPath
	s.Port = port
	s.GoroutineNum = workers
	s.im = &InterceptorManager{}
	s.im.Init()
	s.ContentTypeMap = InitContentType()
	s.Routers = make(map[string]func(*Request, *Response), 0)
}

//Start 启动服务器
func (s *HTTPServer) Start() {
	addr := &net.TCPAddr{}
	addr.Port = s.Port

	listener, err := net.ListenTCP("tcp", addr)
	if err != nil {
		log.Println("tcp listen出错:", err)
		return
	}
	defer listener.Close()

	boss := s.initBossWorkers(s.GoroutineNum)
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
	s.Routers[path] = handler
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
			req.serv = s
			resp := &Response{}
			resp.Headers = make(map[string]string, 0)
			defer s.im.ServerErrorHandle(req, resp)
			s.im.processor.Decode(req, resp)
			//运行拦截器
			s.im.Run(req, resp)
			s.im.processor.Encode(req, resp)
		}
	})
	return boss
}

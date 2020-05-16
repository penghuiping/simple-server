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

//Start 启动服务器
func (serv *HTTPServer) Start() {
	conf := GetConfig()
	conf.routers = make(map[string]func(*Request, *Response))
	conf.interceptors = make([]Interceptor, 0)
	serv.AddInterceptor(&FirstInterceptor{})
	serv.AddInterceptor(&StaticFileInterceptor{})

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

		defer func() {
			if err := recover(); err != nil {
				if err != io.EOF {

				}
				log.Println("job异常是:" + fmt.Sprint(err))
				conn.Close()
			}
		}()

		for {
			req := &Request{}
			req.conn = conn
			resp := &Response{}
			serv.handle(req, resp)
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

func (serv *HTTPServer) handle(req *Request, resp *Response) {
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

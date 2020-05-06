package serv

import (
	"bytes"
	"fmt"
	"log"
	"net"
)

//Dispatcher ...
type Dispatcher struct {
	jobs      chan *net.Conn
	workerNum int
	workers   chan *Worker
}

//NewDispatcher ...
func newDispatcher(workerNum int) *Dispatcher {
	dispatcher := &Dispatcher{}
	dispatcher.workerNum = workerNum
	dispatcher.jobs = make(chan *net.Conn, 2048)
	dispatcher.workers = make(chan *Worker, workerNum)
	for i := 0; i < dispatcher.workerNum; i++ {
		worker := &Worker{}
		worker.jobQueue = make(chan *net.Conn, 10)
		worker.dispatcher = dispatcher
		worker.run()
	}
	return dispatcher
}

//Run ...
func (d *Dispatcher) run() {
	go func() {
		for {
			//从全部job中选一个分发给一个空闲可用的worker
			job := <-d.jobs
			worker := <-d.workers
			worker.jobQueue <- job
		}
	}()
}

//AddJob ...
func (d *Dispatcher) addJob(conn *net.Conn) {
	d.jobs <- conn
}

//Worker ...
type Worker struct {
	dispatcher *Dispatcher
	jobQueue   chan *net.Conn
}

//Run ...
func (w *Worker) run() {
	go func() {
		for {
			//把worker自己注册到dispatcher
			w.dispatcher.workers <- w

			//等待dispatcher给自己分配任务
			conn := *(<-w.jobQueue)

			defer func() {
				if err := recover(); err != nil {
					defer conn.Close()
					log.Println("panic异常是:" + fmt.Sprint(err))
				}
			}()

			total := make([]byte, 0)
			for {
				buf := make([]byte, 512)
				len, err3 := conn.Read(buf)
				if err3 != nil {
					log.Println(err3)
					conn.Close()
					break
				}
				if len > 0 {
					total = bytes.Join([][]byte{total, buf}, []byte{})
					if len < 512 {
						//一个请求完结
						req, _ := parseRequest(total)
						req.remoteAddr = conn.RemoteAddr().String()
						total = make([]byte, 0)

						resp := &Response{}
						resp.init(conn)
						handle(req, resp)
					}
				}
			}

		}
	}()
}

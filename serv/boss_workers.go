package serv

import (
	"fmt"
	"log"
)

//Job ...
type Job struct {
	JobType string
	Content interface{}
}

//Boss ...
type Boss struct {
	jobs      chan *Job
	workerNum int
	workers   chan *worker
	handlers  map[string]func(job *Job)
}

//Start 开启boss线程
func (boss *Boss) Start(workerNum int) {
	boss.workerNum = workerNum
	boss.jobs = make(chan *Job, 2048)
	boss.workers = make(chan *worker, workerNum)
	boss.handlers = make(map[string]func(job *Job), 0)

	//运行worker线程
	for i := 0; i < boss.workerNum; i++ {
		worker := &worker{}
		worker.jobQueue = make(chan *Job, 10)
		worker.boss = boss
		go worker.start()
	}

	go func() {
		for {
			//从全部job中选一个分发给一个空闲可用的worker
			job := <-boss.jobs
			worker := <-boss.workers
			worker.jobQueue <- job
		}
	}()
}

//AddJob ...
func (boss *Boss) AddJob(job *Job) {
	boss.jobs <- job
}

//AddJobHandler ...
func (boss *Boss) AddJobHandler(jobType string, handler func(job *Job)) {
	boss.handlers[jobType] = handler
}

//Worker ...
type worker struct {
	boss     *Boss
	jobQueue chan *Job
}

//开启worker线程
func (w *worker) start() {
	for {
		//把worker自己注册到boss
		w.boss.workers <- w

		//等待boss给自己分配任务
		job := <-w.jobQueue

		defer func() {
			if err := recover(); err != nil {
				log.Println("work工作线程出错，异常是:" + fmt.Sprint(err))
			}
		}()

		handler := w.boss.handlers[job.JobType]
		if handler != nil {
			handler(job)
		}

	}
}

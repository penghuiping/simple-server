package serv

import (
	"fmt"
	"strconv"
	"sync"
	"testing"
)

func Test1(in *testing.T) {
	countDown := &sync.WaitGroup{}
	countDown.Add(10)
	boss := &Boss{}
	boss.Start(10)
	for i := 0; i < 10; i++ {
		job := &Job{JobType: "msg", Content: "this is a job " + strconv.Itoa(i)}
		boss.AddJob(job)
	}
	boss.AddJobHandler("msg", func(job *Job) {
		defer countDown.Done()
		fmt.Println(job.Content)
	})
	countDown.Wait()
}

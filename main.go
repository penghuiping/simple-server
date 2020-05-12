package main

import (
	"log"
	_ "net/http/pprof"
	"os"

	"github.com/penghuiping/simple-server/serv"
)

func main() {
	// file, _ := os.OpenFile("./server.log", os.O_CREATE|os.O_RDWR|os.O_APPEND, 0666)
	// defer file.Close()

	log.SetOutput(os.Stdout)

	//配置pprof
	// go func() {
	// 	http.ListenAndServe("0.0.0.0:15672", nil)
	// }()

	config := &serv.Config{}
	config.HTMLPath = "./html"
	config.Port = 8080
	config.GoroutineNum = 1000
	serv.SetConfig(config)

	// serv.AddRoute("/hello", func(a *serv.Request, b *serv.Response) {
	// 	b.Body("你好")
	// 	b.Header("Content-Type", "text/html;charset=utf-8")
	// })

	// serv.AddRoute("/world", func(a *serv.Request, b *serv.Response) {
	// 	b.Body("世界")
	// 	b.Header("Content-Type", "text/html;charset=utf-8")
	// })

	httpServ := serv.HTTPServer{}
	httpServ.Start()
}

package main

import (
	"log"
	_ "net/http/pprof"
	"os"
	"strings"

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

	httpServ := serv.HTTPServer{}

	httpServ.Init("./html", 1000, 8080)

	httpServ.AddRoute("/hello", func(req *serv.Request, resp *serv.Response) {
		resp.Header("Content-Type", "text/html;charset=utf-8")
		resp.Body(strings.NewReader("你好"))
	})

	httpServ.AddRoute("/world", func(req *serv.Request, resp *serv.Response) {
		resp.Header("Content-Type", "text/html;charset=utf-8")
		resp.Body(strings.NewReader("世界"))
	})

	httpServ.Start()

}

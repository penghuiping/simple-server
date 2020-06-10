package serv

import (
	"bufio"
	"strings"
)

func handleError(err interface{}, req *Request, resp *Response) {
	switch err {
	default:
		handle500Error(req, resp)
		break
	}
}

func handle500Error(req *Request, resp *Response) {
	resp.code = StatusInternalServerError
	resp.codeMsg = "server error"
	resp.headers["Content-Type"] = "text/html;charset=utf-8"
	bodyContent := "500 服务器错误，请联系相关技术人员解决\r\n"
	resp.bodySize = int64(len([]byte(bodyContent)))
	resp.body = bufio.NewReader(strings.NewReader(bodyContent))
}

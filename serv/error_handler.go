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
	resp.Code = StatusInternalServerError
	resp.CodeMsg = "server error"
	resp.Headers["Content-Type"] = "text/html;charset=utf-8"
	bodyContent := "500 服务器错误，请联系相关技术人员解决\r\n"
	resp.BodySize = int64(len([]byte(bodyContent)))
	resp.Body = bufio.NewReader(strings.NewReader(bodyContent))
}

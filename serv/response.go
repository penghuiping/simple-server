package serv

import (
	"io"
)

//Response ...
type Response struct {
	headers  map[string]string
	code     int
	codeMsg  string
	bodySize int64
	body     io.Reader
}

//Body ...
func (res *Response) Body(body io.Reader) {
	res.body = body
}

//Header ...
func (res *Response) Header(name string, value string) {
	res.headers[name] = value
}

//Code ...
func (res *Response) Code(code int) {
	res.code = code
}

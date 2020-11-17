package serv

//ProtocolProcessor 协议处理器
type ProtocolProcessor interface {

	//Encode  编码
	Encode(request *Request, response *Response)

	//Decode 解码
	Decode(request *Request, response *Response)
}

package serv

import (
	"bytes"
	"strings"
)

//Request ...
type Request struct {
	protocal   string
	uri        string
	method     string
	headers    map[string]string
	body       []byte
	remoteAddr string
}

func parseRequest(content []byte) (*Request, error) {
	req := &Request{}
	req.headers = make(map[string]string, 0)

	contentStr := string(content)
	lines := strings.Split(contentStr, "\r\n")

	isHeaderPart := false

	for i, line := range lines {
		if i == 0 {
			firstLine := strings.Split(line, " ")
			req.method = firstLine[0]
			req.uri = firstLine[1]
			req.protocal = firstLine[2]
			isHeaderPart = true
			continue
		}

		if i > 0 {
			if !IsBlankStr(line) && isHeaderPart {
				header := strings.Split(line, ":")
				req.headers[header[0]] = header[1]
				continue
			} else if IsBlankStr(line) && isHeaderPart && IsBlankStr(lines[i+1]) {
				isHeaderPart = false
			} else {
				req.body = bytes.Join([][]byte{req.body, []byte(line)}, []byte{})
			}
		}

	}
	return req, nil
}

package serv

import "sync"

//Config 全局配置类
type Config struct {
	HTMLPath       string
	Port           int
	GoroutineNum   int
	contentTypeMap map[string]string
	routers        map[string]func(*Request, *Response)
	interceptors   []Interceptor
}

var config *Config
var mutex sync.RWMutex = sync.RWMutex{}

//GetConfig 获取全局配置
func GetConfig() *Config {
	mutex.RLock()
	defer mutex.RUnlock()
	return config
}

//SetConfig 设置全局配置
func SetConfig(newConfig *Config) {
	mutex.Lock()
	defer mutex.Unlock()
	newConfig.contentTypeMap = make(map[string]string, 0)
	newConfig.contentTypeMap[".html"] = "text/html;charset=utf-8"
	newConfig.contentTypeMap[".css"] = "text/css;charset=utf-8"
	newConfig.contentTypeMap[".js"] = "application/x-javascript"
	newConfig.contentTypeMap[".gif"] = "image/gif"
	newConfig.contentTypeMap[".png"] = "image/png"
	newConfig.contentTypeMap[".woff"] = "application/x-font-woff"
	newConfig.contentTypeMap[".woff2"] = "application/x-font-woff"
	newConfig.contentTypeMap[".gz"] = "application/x-gzip"
	newConfig.contentTypeMap[".zip"] = "application/x-zip-compressed"
	newConfig.contentTypeMap[".rar"] = "application/octet-stream"
	newConfig.contentTypeMap[".mp4"] = "video/mp4"
	newConfig.contentTypeMap[".mp3"] = "audio/mpeg"
	newConfig.contentTypeMap[".7z"] = "application/x-7z-compressed"
	newConfig.contentTypeMap[".pdf"] = "application/pdf"
	config = newConfig
}

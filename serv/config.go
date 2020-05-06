package serv

import "sync"

//Config ...
type Config struct {
	HTMLPath     string
	Port         int
	GoroutineNum int
}

var config Config
var mutex sync.RWMutex = sync.RWMutex{}

//GetConfig ...
func GetConfig() Config {
	mutex.RLock()
	defer mutex.RUnlock()
	return config
}

//SetConfig ...
func SetConfig(newConfig Config) {
	mutex.Lock()
	defer mutex.Unlock()
	config = newConfig
}

package conf

import (
	"github.com/BurntSushi/toml"
	"log"
	"sync"
)

const (
	Ratio = 100000000
	//MaxLockBlockCount = 16
	//configFile = "/Users/logan/Workspace/github/forks/kahf/config.toml"
	configFile = "config.toml"
)

var Setting *Config
var once sync.Once

func init() {

	once.Do(func() {
		//cfg = &Config{}
		_, err := toml.DecodeFile(configFile, &Setting)
		if err != nil {
			log.Println(err)
		}
		DecodeSetting("")
	})
}

func DecodeSetting(key string) {
	Setting.Miner.decode(key)
	Setting.Rpc.decode(key)
}

package conf

import "strings"

type Config struct {
	Rpc     Rpc    `toml:"rpc"`
	Log     Log    `toml:"log"`
	Version string `toml:"version"`
}

type Rpc struct {
	Host string `toml:"host"`
	Auth string `toml:"auth"`
	User string
	Pwd  string
}

type Log struct {
	Mode       string `toml:"mode"`
	Level      string `toml:"level"`
	BufferSize int    `toml:"buffersize"`
}

func (r *Rpc) decode(key string) {
	auth := r.Auth
	if len(key) > 0 {
		//TODO: auth=
	}
	acct := strings.Split(auth, ":")
	r.User = acct[0]
	r.Pwd = acct[1]
}

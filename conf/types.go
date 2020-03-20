package conf

import "strings"

type Config struct {
	Miner   Miner  `toml:"miner"`
	Rpc     Rpc    `toml:"rpc"`
	Api     Api    `toml:"api"`
	Auth    Auth   `toml:"auth"`
	Mysql   Mysql  `toml:"mysql"`
	Log     Log    `toml:"log"`
	Version string `toml:"version"`
}

type Miner struct {
	Address    string `toml:"address"`
	Key        string `toml:"key"`
	PrivateKey string
}

type Rpc struct {
	Host                      string `toml:"host"`
	Auth                      string `toml:"auth"`
	SyncBlockInterval         int64  `toml:"syncBlockInterval"`
	SyncMemoryInterval        int64  `toml:"syncMemoryInterval"`
	SyncPeerInfoInterval      int64  `toml:"syncPeerInfoInterval"`
	SyncNodeInfoInterval      int64  `toml:"syncNodeInfoInterval"`
	SyncLockedBlockInterval   int64  `toml:"syncLockedBlockInterval"`
	SyncInvalidBlockInterval  int64  `toml:"syncInvalidBlockInterval"`
	SyncMemoryTxStateInterval int64  `toml:"syncMemoryTxStateInterval"`
	User                      string
	Pwd                       string
}

type Api struct {
	ListenPort string `toml:"listenPort"`
}

type Auth struct {
	Jwt            string `toml:"jwt"`
	ExpirationTime int64  `toml:"expirationTime"`
	Issuer         string `toml:"issuer"`
	SecretKey      string `toml:"secretKey"`
}

type Mysql struct {
	User         string `toml:"user"`
	Password     string `toml:"password"`
	Address      string `toml:"address"`
	DBName       string `toml:"dbname"`
	Prefix       string `toml:"prefix"`
	MaxIdleConns int    `toml:"maxidleconns"`
	MaxOpenConns int    `toml:"maxopenconns"`
}

type Log struct {
	Mode       string `toml:"mode"`
	Level      string `toml:"level"`
	BufferSize int    `toml:"buffersize"`
}

func (m *Miner) decode(key string) {
	m.PrivateKey = m.Key
	if len(key) > 0 {
		//TODO: m.PrivateKey =
	}
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

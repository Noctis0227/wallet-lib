package core

import (
	"bytes"
	"encoding/binary"
	"encoding/gob"
	"git.diabin.com/BlockChain/wallet-lib/log"
)

var empty []byte

type UintLong uint64
type Bytes []byte
type Output struct {
	Id         int64  `xorm:"bigint autoincr pk"`
	TxId       string `xorm:"varchar(64) index"`
	OutIndex   int    `xorm:"int index"`
	BlockHash  string `xorm:"varchar(64)"`
	BlockOrder uint64 `xorm:"bigInt"`
	BlockId    uint64 `json:"-" xorm:"BigInt"`
	Address    string `xorm:"varchar(35) index"`
	Amount     uint64 `xorm:"bigInt"`
	Timestamp  string `xorm:"varchar(19)"`
	SpentTxId  string `xorm:"varchar(64)"`
	Stat       int8   `xorm:"int"`
}

func (o *Output) Error() string {
	panic("implement me")
}

type HistoryId struct {
	Id                  uint64 `xorm:"bigint pk"`
	LastTxBlockId       uint64 `xorm:"bigint"`
	LastCoinBaseBlockId uint64 `xorm:"bigint"`
}

func init() {
	empty = []byte("")
}

func (o *Output) ToBytes() []byte {
	return toBytes(o)
}

func (i UintLong) ToBytes() []byte {
	var buf = make([]byte, 8)
	binary.BigEndian.PutUint64(buf, uint64(i))
	return buf
}

func (b Bytes) ToOutput() *Output {
	var rs = &Output{}
	fromBytes(b, rs)
	return rs
}

func (b Bytes) ToUint64() uint64 {
	if bytes.Equal(b, empty) {
		return 0
	}
	return binary.BigEndian.Uint64(b)
}

func Uint64ToBytes(val uint64) []byte {
	var buf = make([]byte, 8)
	binary.BigEndian.PutUint64(buf, val)
	return buf
}

func BytesToUint64(val []byte) uint64 {
	return binary.BigEndian.Uint64(val)
}

func toBytes(s interface{}) []byte {
	var result bytes.Buffer
	encoder := gob.NewEncoder(&result)
	err := encoder.Encode(s)
	if err != nil {
		log.Errorf("An error occurred when to bytes,%s", err)
	}
	return result.Bytes()
}

func fromBytes(val []byte, obj interface{}) {
	decoder := gob.NewDecoder(bytes.NewReader(val))
	err := decoder.Decode(obj)
	if err != nil {
		log.Errorf("An error occurred when from bytes,%s", err)
	}
}

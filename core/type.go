package core

import (
	"bytes"
	"encoding/binary"
	"encoding/gob"
	"kahf/log"
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

type TxInputs struct {
	TxId     string
	OutIndex int
	InTime   string
}

type FaucetRecord struct {
	Address string
	Amount  uint64
	Status  int
	Memo    string
	InTime  string
}

type AddressMap struct {
	UsableAmount   uint64
	LockedAmount   uint64
	MemoryAmount   uint64
	SpentAmount    uint64
	Unspents       UniqueList
	Spents         UniqueList
	MemoryUnspents UniqueList
	MemorySpents   UniqueList
	Outs           UniqueList
	TxRecords      UniqueList
}

type AddressInfo struct {
	UseableAmount float64     `json:"useableAmount"`
	LockedAmount  float64     `json:"lockedAmount"`
	MemoryAmount  float64     `json:"memoryAmount"`
	SpentAmount   float64     `json:"spentAmount"`
	Balance       float64     `json:"balance"`
	TxRecords     interface{} `json:"txrecords"`
}

//transaction record
type TxRecord struct {
	Id         int64  `xorm:"bigint autoincr pk"`
	TxId       string `json:"txid" xorm:"Varchar(64) unique index"`
	From       string `json:"from" xorm:"Varchar(35) 'addr_from'"`
	To         string `json:"to" xorm:"Varchar(35) 'addr_to'"`
	Amount     uint64 `json:"amount" xorm:"BigInt"`
	ToCount    int    `json:"tocount" xorm:"Int"`
	Spent      uint64 `json:"spent" xorm:"BigInt"`
	Change     uint64 `json:"change" xorm:"BigInt 'change_amount'"`
	Fee        uint64 `json:"fee" xorm:"BigInt"`
	BlockOrder uint64 `json:"blockorder" xorm:"BigInt"`
	BlockId    uint64 `json:"-" xorm:"BigInt"`
	InTime     string `json:"inTime" xorm:"Varchar(19)"`
	Stat       int8   `json:"stat" xorm:"int"`
}

type JsonTxRecord struct {
	Id         int64            `xorm:"bigint autoincr pk"`
	TxId       string           `json:"txid" xorm:"Varchar(64) unique index"`
	From       string           `json:"from" xorm:"Varchar(35) 'addr_from'"`
	Amount     uint64           `json:"amount" xorm:"BigInt"`
	ToCount    int              `json:"tocount" xorm:"Int"`
	Spent      uint64           `json:"spent" xorm:"BigInt"`
	Change     uint64           `json:"change" xorm:"BigInt 'change_amount'"`
	Fee        uint64           `json:"fee" xorm:"BigInt"`
	BlockOrder uint64           `json:"blockorder" xorm:"BigInt"`
	InTime     string           `json:"inTime" xorm:"Varchar(19)"`
	Stat       int8             `json:"stat" xorm:"int"`
	Outs       []*JsonOutRecord `json:"outs" xorm:"-"`
}

type JsonOutRecord struct {
	OutIndex  int    `json:"index"`
	Address   string `json:"address"`
	Amount    uint64 `json:"amount"`
	SpentTxId string `json:"spenttx"`
	IsChange  bool   `json:"ischange"`
}

type BlockRecord struct {
	Id            uint64 `xorm:"bigint pk"`
	Hash          string `json:"hash" xorm:"Varchar(64) unique index"`
	BlockOrder    uint64 `json:"blockorder" xorm:"BigInt"`
	Transactions  int    `json:"transactions" xorm:"Int"`
	CreateTime    string `json:"createTime" xorm:"Varchar(19)"`
	Size          int    `json:"size" xorm:"Int"`
	Miner         string `json:"miner" xorm:"varchar(35)"`
	Confirmations uint32 `json:"confirmations" xorm:"int"`
	Stat          int8   `json:"stat" xorm:"int"` // 0:normal,1:unconfirmed,2:inmemory,3:invalid
}

type PeerRecord struct {
	Id         uint64 `xorm:"bigint autoincr pk"`
	Addr       string `json:"addr" xorm:"Varchar(35) index"`
	ConnTime   string `json:"conntime" xorm:"Varchar(30)"`
	SubVer     string `json:"subver" xorm:"Varchar(64)"`
	MainOrder  uint64 `json:"mainorder" xorm:"BigInt"`
	Layer      uint64 `json:"layer" xorm:"BigInt"`
	MainHeight uint64 `json:"mainheight" xorm:"BigInt"`
}

type MinerRecord struct {
	Id         int64      `xorm:"bigint autoincr pk"`
	Addr       string     `json:"address" xorm:"varchar(35)"`
	CoinBase   uint64     `json:"balance" xorm:"-"`
	Orders     UniqueList `json:"-" xorm:"-"`
	BlockCount int        `json:"count" xorm:"-"`
}

type AccountRecord struct {
	Addr    string  `json:"addr"`
	Balance float64 `json:"balance"`
}

//type Stats struct {
//	LatestHeight  UintLong  `json:"latestHeight"`
//	LatestBlocks  StringStack `json:"latestBlocks"`
//	LatestTxs     StringStack `json:"latestTxs"`
//	ConfirmedTx   UintLong  `json:"confirmedTx"`
//	UnconfirmedTx UintLong  `json:"unconfirmedTx"`
//}

type StatsKeys struct {
	LatestId         []byte
	LatestHeight     []byte
	TxNum            []byte
	UnconfirmedTxNum []byte
	BlockNum         []byte
	LatestTxRecords  []byte
	PeerInfoCount    []byte
	AccountNumber    []byte
}

func init() {
	empty = []byte("")
}

func (a *AddressMap) ToBytes() []byte {
	return toBytes(a)
}

func (o *Output) ToBytes() []byte {
	return toBytes(o)
}

func (t *TxRecord) ToBytes() []byte {
	return toBytes(t)
}

func (t *BlockRecord) ToBytes() []byte {
	return toBytes(t)
}

func (i UintLong) ToBytes() []byte {
	var buf = make([]byte, 8)
	binary.BigEndian.PutUint64(buf, uint64(i))
	return buf
}

func (p *PeerRecord) ToBytes() []byte {
	return toBytes(p)
}

//func (s *Stats) toBytes() []byte {
//	return toBytes(s)
//}

func (b Bytes) ToAddressMap() *AddressMap {
	var rs = &AddressMap{}
	fromBytes(b, rs)
	return rs
}

func (b Bytes) ToOutput() *Output {
	var rs = &Output{}
	fromBytes(b, rs)
	return rs
}

func (b Bytes) ToTxRecord() *TxRecord {
	var rs = &TxRecord{}
	fromBytes(b, rs)
	return rs
}

func (b Bytes) ToBlockRecord() *BlockRecord {
	var rs = &BlockRecord{}
	fromBytes(b, rs)
	return rs
}

func (b Bytes) ToUint64() uint64 {
	if bytes.Equal(b, empty) {
		return 0
	}
	return binary.BigEndian.Uint64(b)
}

func (b Bytes) ToPeerRecord() *PeerRecord {
	var rs = &PeerRecord{}
	fromBytes(b, rs)
	return rs
}

//func (s *Stats) fromBytes(val []byte) {
//	fromBytes(val, s)
//}

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

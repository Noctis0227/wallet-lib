package core

import (
	"github.com/Noctis0227/wallet-lib/rpc"
)

var Storage IStorage
var List IList

type IStorage interface {
	Init()
	Start()
	End()
	UpdateOutputs(outs []Output) error
	UpdateOutputsState(outs []Output) error
	UpdateTxOutputSpent(txId string, outIndex int, spentTxId string) error
	UpdateMainHeight(mainHeight uint64) error
	ClearTxOutputSpent(spentId string) error
	GetOutput(key string) *Output
	GetTxOutputs(txId string) []Output
	GetLastId() *HistoryId
	UpdateLastId(lastId *HistoryId)
}

type IList interface {
	QueryOuts(condition interface{}, value []interface{}) []Output
}

type context struct {
	Hash         string
	Order        uint64
	Id           uint64
	Timestamp    string
	IsSyncMemory bool
	Invalid      bool
	Outs         []Output
	SpentOuts    []Output
	TxRecords    UniqueList
	Miner        string
	BlockColor   int
	BlockStat    int8
}

func isCoinBase(tx *rpc.Transaction) bool {
	if tx != nil && len(tx.Vin) > 0 && tx.Vin[0].Coinbase != "" {
		return true
	}
	return false
}

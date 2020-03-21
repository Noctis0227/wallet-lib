package core

import (
	"fmt"
	"git.diabin.com/BlockChain/wallet-lib/log"
	"git.diabin.com/BlockChain/wallet-lib/rpc"
	"strconv"
	"strings"
	"time"
)

var Storage IStorage
var Stats IStats
var List IList
var StatsKey StatsKeys
var QueryOutput IQueryOutput

type IStorage interface {
	Init()
	Start()
	End()
	UpdateOutputs(outs []Output) error
	UpdateOutputsState(outs []Output) error
	UpdateTxOutputSpent(txId string, outIndex int, spentTxId string) error
	UpdateAddressMap(addrs map[string]*AddressMap) error
	UpdateLockedMap(id uint64, stat int8) error
	UpdateInvalidMap(id uint64, stat int8) error
	UpdateTxRecord(txs UniqueList) error
	UpdateBlockRecord(b BlockRecord) error
	UpdatePeerRecord(p *PeerRecord) error
	UpdateMainHeight(mainHeight uint64) error
	ClearTxOutputSpent(spentId string) error
	GetOutput(key string) *Output
	GetTxOutputs(txId string) []Output
	GetAddressMap(key string) *AddressMap
	GetLockedMap() []uint64
	GetInvalidMap() []uint64
	GetTxRecord(key string) *TxRecord
	GetBlockRecordByHeight(height uint64) *BlockRecord
	GetBlockRecordByHash(key string) *BlockRecord
	GetPeerRecord(key string) *PeerRecord
	GetMinerRecordByHeight(height uint64) *MinerRecord
	GetMinerRecordByAddr(addr string) *MinerRecord
	GetMinerRecords(start, end uint64) []map[string]string
	GeMinersBlocks(addr string, start, end uint64) []map[string]string
	GetMemoryTxs() []TxRecord
	GetFailedTxs() []TxRecord
	ClearPeerRecord()
	//GetBlockRecords() []BlockRecord
}

type IList interface {
	QueryLatestTx(idx int, size int) []TxRecord
	QueryLatestRecordOuts(txid string, idx int, size int) []Output
	QueryAddrLastestTx(addr string, idx int, size int) []TxRecord
	QueryAddrRecordCount(addr string) int64
	QueryLatestBlocks(idx int, size int) []BlockRecord
	QueryLatestPeers(idx int, size int) []PeerRecord
	QueryLatestAccounts(idx int, size int) []AccountRecord
	QueryLatestMiners(idx int, size int) []MinerRecord
	QueryOuts(condition interface{}, value []interface{}) []Output
	QueryTx(condition interface{}, value []interface{}) []TxRecord
	QueryBlocks(condition interface{}, value []interface{}) []BlockRecord
}

type IStats interface {
	AddLatestTx(txId string)
	GetStats() map[string]string
	GetStatsItem(key []byte) Bytes
	UpdateStats(key []byte, val string)
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
	AddrMaps     map[string]*AddressMap
	TxRecords    UniqueList
	Miner        string
	BlockColor   int
	BlockStat    int8
}

func init() {
	StatsKey = StatsKeys{
		LatestId:         []byte("latestId"),
		LatestHeight:     []byte("latestHeight"),
		TxNum:            []byte("txNum"),
		UnconfirmedTxNum: []byte("unconfirmedTxNum"),
		BlockNum:         []byte("blockNum"),
		LatestTxRecords:  []byte("latestTxRecords"),
		PeerInfoCount:    []byte("peerInfoCount"),
		AccountNumber:    []byte("accountNumber"),
	}
}

func NewContext(hash string, order uint64, id uint64, time string, isMemory bool, invalid bool, blockColor int, blockStat int8) *context {
	return &context{Hash: hash,
		Order:        order,
		Id:           id,
		Timestamp:    time,
		IsSyncMemory: isMemory,
		Invalid:      invalid,
		Miner:        "",
		Outs:         []Output{},
		SpentOuts:    []Output{},
		AddrMaps:     map[string]*AddressMap{},
		TxRecords:    UniqueList{},
		BlockColor:   blockColor,
		BlockStat:    blockStat,
	}
}

func (cx *context) saveTransactions(txs []rpc.Transaction) error {
	for _, tx := range txs {
		cx.updateTempOuts(tx)
	}

	if err := Storage.UpdateOutputs(cx.Outs); err != nil {
		return err
	}

	for _, tx := range txs {
		cx.updateTmpSpentOuts(tx)
	}
	for _, out := range cx.SpentOuts {
		if err := Storage.UpdateTxOutputSpent(out.TxId, out.OutIndex, out.SpentTxId); err != nil {
			return err
		}
	}

	cx.updateAddrMap()
	for _, tx := range txs {
		cx.updateTxRecord(tx)
	}
	if err := Storage.UpdateTxRecord(cx.TxRecords); err != nil {
		return err
	}
	if err := Storage.UpdateAddressMap(cx.AddrMaps); err != nil {
		return err
	}
	return nil
}

func (cx *context) updateTempOuts(tx rpc.Transaction) {
	var stat = getTxStat(tx.Confirmations, isCoinBase(&tx), cx.Invalid, cx.BlockColor, cx.BlockStat)
	for idx, vo := range tx.Vout {
		if vo.ScriptPubKey.Addresses == nil || len(vo.ScriptPubKey.Addresses) == 0 {
			continue
		}
		out := Output{
			TxId:       tx.Txid,
			BlockOrder: cx.Order,
			BlockId:    cx.Id,
			BlockHash:  cx.Hash,
			Address:    vo.ScriptPubKey.Addresses[0],
			Amount:     vo.Amount,
			OutIndex:   idx,
			SpentTxId:  "",
			Timestamp:  cx.Timestamp,
			Stat:       stat,
		}
		cx.Outs = append(cx.Outs, out)
	}
}

func (cx *context) updateTmpSpentOuts(tx rpc.Transaction) {
	for _, in := range tx.Vin {
		key := GetOutKey(in.Txid, in.Vout)
		if it := Storage.GetOutput(key); it != nil {
			it.SpentTxId = tx.Txid
			cx.SpentOuts = append(cx.SpentOuts, *it)
		}
	}
}

func (cx *context) updateCoinBaseTxRecord(tx rpc.Transaction) {
	state := getTxStat(tx.Confirmations, true, cx.Invalid, cx.BlockColor, cx.BlockStat)
	txrd := TxRecord{
		TxId:       tx.Txid,
		BlockOrder: cx.Order,
		From:       "Miner",
		To:         "",
		ToCount:    1,
		Amount:     uint64(0),
		Spent:      uint64(0),
		Change:     uint64(0),
		Stat:       state,
		BlockId:    cx.Id,
		InTime:     cx.Timestamp}
	for _, vout := range tx.Vout {
		if vout.ScriptPubKey.Addresses == nil || len(vout.ScriptPubKey.Addresses) == 0 {
			continue
		}
		if vout.Amount > txrd.Amount {
			txrd.Amount = vout.Amount
			cx.Miner = vout.ScriptPubKey.Addresses[0]
			txrd.To = cx.Miner
		}
	}
	cx.updateAddrTx(txrd.To, txrd.TxId)
	cx.TxRecords.Add(txrd)
	Stats.AddLatestTx(tx.Txid)
}

func (cx *context) updateTxRecord(tx rpc.Transaction) {
	if isCoinBase(&tx) {
		cx.updateCoinBaseTxRecord(tx)
		return
	}
	s := getTxStat(tx.Confirmations, false, cx.Invalid, cx.BlockColor, cx.BlockStat)
	toCount := 0
	txrd := TxRecord{
		TxId:       tx.Txid,
		BlockOrder: cx.Order,
		BlockId:    cx.Id,
		From:       "",
		To:         "",
		Amount:     uint64(0),
		Spent:      uint64(0),
		Change:     uint64(0),
		Stat:       s,
		InTime:     cx.Timestamp}

	for _, in := range tx.Vin {
		key := GetOutKey(in.Txid, in.Vout)
		if it := Storage.GetOutput(key); it != nil {
			txrd.From = it.Address
		}

		if in.Amountin == 0 {
			if out := Storage.GetOutput(key); out != nil {
				in.Amountin = out.Amount
			}
		}
		txrd.Spent += in.Amountin
	}
	for _, out := range tx.Vout {
		addrs := out.ScriptPubKey.Addresses
		if addrs == nil || len(addrs) == 0 {
			continue
		}
		if addrs[0] == txrd.From {
			txrd.Change += out.Amount
		} else {
			txrd.To = addrs[0]
			txrd.Amount += out.Amount
			toCount++
		}
	}
	txrd.ToCount = toCount
	txrd.Fee = txrd.Spent - txrd.Amount - txrd.Change
	if int64(txrd.Fee) < 0 {
		txrd.Fee = 0
	}
	cx.updateAddrTx(txrd.From, txrd.TxId)
	cx.updateAddrTx(txrd.To, txrd.TxId)
	cx.TxRecords.Add(txrd)
	Stats.AddLatestTx(tx.Txid)
}

func (cx *context) updateAddrMap() {
	if IsServer {
		return
	}
	outs := append(cx.Outs, cx.SpentOuts...)
	for _, out := range outs {
		key := GetOutKey(out.TxId, out.OutIndex)
		it, has := cx.AddrMaps[out.Address]
		if !has {
			it = Storage.GetAddressMap(out.Address)
		}

		if nil != it {
			it.Outs.Add(key)
		} else {
			it = &AddressMap{Outs: UniqueList{}}
			it.Outs.Add(key)
		}
		cx.AddrMaps[out.Address] = it
	}

	for _, it := range cx.AddrMaps {
		addToAddrMap(it)
	}
}

func (cx *context) updateAddrTx(addr string, txid string) {
	if IsServer {
		return
	}
	if addr == "" || addr == "Miner" {
		return
	}
	_, ok := cx.AddrMaps[addr]
	if !ok {
		log.Debugf("update addr %s tx addrMap is nil, txid = %s", addr, txid)
		item := Storage.GetAddressMap(addr)
		if item == nil {
			item = &AddressMap{Outs: UniqueList{}}
			cx.AddrMaps[addr] = item
		}
	}

	it := cx.AddrMaps[addr]
	if it.TxRecords == nil {
		it.TxRecords = UniqueList{}
	}
	it.TxRecords.Add(txid)
}

func GetOutKey(txid string, idx interface{}) string {
	return fmt.Sprintf("%s-%d", txid, idx)
}

func SaveFromBlock(b *rpc.Block) error {
	blockColor, err := rpc.IsBlue(b.Hash)
	if err != nil {
		return err
	}
	blockStat := getBlockStat(b.Confirmations, b.Txsvalid, blockColor)
	ds := NewContext(b.Hash, b.Order, b.Id, strings.Split(b.Timestamp.UTC().String(), "+")[0], false, b.Txsvalid, blockColor, blockStat)
	if err := ds.saveTransactions(b.Transactions); err != nil {
		return err
	}
	brd := BlockRecord{
		Id:            b.Id,
		Hash:          b.Hash,
		BlockOrder:    b.Order,
		Transactions:  len(b.Transactions),
		Size:          b.Size(),
		CreateTime:    ds.Timestamp,
		Miner:         ds.Miner,
		Confirmations: b.Confirmations,
		Stat:          blockStat,
	}
	if err := Storage.UpdateLockedMap(brd.Id, brd.Stat); err != nil {
		return err
	}
	if err := Storage.UpdateInvalidMap(brd.Id, brd.Stat); err != nil {
		return err
	}
	if err := Storage.UpdateBlockRecord(brd); err != nil {
		return err
	}
	if err := Storage.UpdateMainHeight(NodeInfo.MainHeight); err != nil {
		return err
	}
	return nil
}

func SaveFromMemory(txs []rpc.Transaction) error {
	ds := NewContext("", 0, 0, "", true, true, 1, State_InMemory)
	if err := ds.saveTransactions(txs); err != nil {
		return err
	}
	Stats.UpdateStats(StatsKey.UnconfirmedTxNum, strconv.Itoa(len(txs)))
	return nil
}

func SaveFromPeers(peerInfo []rpc.PeerInfo) {
	Storage.ClearPeerRecord()
	peerRecord := &PeerRecord{}
	addrMap := make(map[string]bool)
	for _, peer := range peerInfo {
		peerRecord.Id = peer.Id
		peerRecord.Addr = peer.Addr
		peerRecord.ConnTime = time.Unix(int64(peer.Conntime), 0).String()
		peerRecord.SubVer = peer.Subver
		peerRecord.MainHeight = peer.GraphState.MainHeight
		peerRecord.Layer = peer.GraphState.Layer
		peerRecord.MainOrder = peer.GraphState.Mainorder
		Storage.UpdatePeerRecord(peerRecord)
		addrMap[peerRecord.Addr] = true
	}
	Stats.UpdateStats(StatsKey.PeerInfoCount, strconv.Itoa(len(addrMap)))
}

func SaveLastBlockId(id uint64) {
	Stats.UpdateStats(StatsKey.LatestId, strconv.FormatUint(id, 10))
}

func addToAddrMap(item *AddressMap) {
	for i, key := range item.Outs {
		if i == 0 {
			item.Unspents = UniqueList{}
			item.Spents = UniqueList{}
			item.MemoryUnspents = UniqueList{}
			item.MemorySpents = UniqueList{}

			item.UsableAmount = 0
			item.MemoryAmount = 0
			item.SpentAmount = 0
		}
		if out := Storage.GetOutput(key.(string)); out != nil {
			if out.BlockOrder != 0 {
				if out.SpentTxId != "" {
					item.Spents.Add(key)
				} else {
					item.Unspents.Add(key)
				}
			} else {
				if out.SpentTxId != "" {
					item.MemorySpents.Add(key)
				} else {
					item.MemoryUnspents.Add(key)
				}
			}
		}
	}
}

func getBlockStat(confirmations uint32, invalid bool, blockColor int) int8 {
	if !invalid {
		return State_Invaild
	}
	if confirmations <= uint32(NodeInfo.Confirmations) || confirmations <= uint32(NodeInfo.Coinbasematurity) {
		return State_Unconfirmed
	}
	switch blockColor {
	case 0:
		return State_Red
	case 1:
		return State_Confirmed
	case 2:
		return State_Unconfirmed
	}
	return State_Confirmed
}

func isNeedChangeStat(confirmations uint32) bool {
	if confirmations > uint32(NodeInfo.Confirmations) || confirmations > uint32(NodeInfo.Coinbasematurity) {
		return true
	}
	return false
}

func getTxStat(confirmations uint32, isCoinBase bool, invalid bool, blockColor int, blockStat int8) int8 {
	if isCoinBase {
		if confirmations <= uint32(NodeInfo.Coinbasematurity) {
			return State_Unconfirmed
		}
		if !invalid {
			return State_Failed
		}
		switch blockColor {
		case 0:
			return State_Failed
		case 1:
			return State_Confirmed
		case 2:
			return State_Unconfirmed
		}
	} else {
		if blockStat == State_InMemory {
			return State_InMemory
		}
		if confirmations <= uint32(NodeInfo.Confirmations) {
			return State_Unconfirmed
		} else {
			if !invalid {
				return State_Failed
			}
			switch blockStat {
			case State_Confirmed:
				return State_Confirmed
			case State_Unconfirmed:
				return State_Confirmed
			case State_Invaild:
				return State_Failed
			}
		}
	}
	return State_Confirmed
}

func isCoinBase(tx *rpc.Transaction) bool {
	if tx != nil && len(tx.Vin) > 0 && tx.Vin[0].Coinbase != "" {
		return true
	}
	return false
}

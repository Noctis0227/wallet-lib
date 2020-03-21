package core

import (
	"fmt"
	"git.diabin.com/BlockChain/wallet-lib/conf"
	"git.diabin.com/BlockChain/wallet-lib/log"
	"git.diabin.com/BlockChain/wallet-lib/rpc"
	"strconv"
	"strings"
	"sync"
	"time"
)

//type Height in
var NodeInfo *rpc.NodeInfo
var currId = uint64(1)
var LatestHeight = uint64(0)
var FirstShow bool
var showFlag = 0
var IsServer = false

var notSyncedTxs map[string]rpc.Transaction
var lock = &sync.Mutex{}

func rmSyncedTx(txs []rpc.Transaction) {
	lock.Lock()
	defer lock.Unlock()
	for _, it := range txs {
		delete(notSyncedTxs, it.Txid)
	}
}

func rmSyncedTxById(txIds []string) {
	lock.Lock()
	defer lock.Unlock()
	for _, txid := range txIds {
		delete(notSyncedTxs, txid)
	}
}

func getNotSyncedTxs() []rpc.Transaction {
	rs := make([]rpc.Transaction, 0)
	for _, it := range notSyncedTxs {
		rs = append(rs, it)
	}
	return rs
}

func testIsBlue() {
	outs, _ := GetUsableOuts("TmjTooPeHr27TLkJzvXM9NabbyTEqXY2Bay")
	for _, o := range outs {
		if o.Stat == 5 {
			rs, err := rpc.IsBlue(o.BlockHash)
			if err != nil {
				fmt.Println("rpc error ", err)
				continue
			}
			fmt.Println(o.BlockHash, ":", rs)
		}
	}
}

func StartSync() {
	notSyncedTxs = make(map[string]rpc.Transaction)
	SyncNodeInfo()
	SyncPeerInfo()
	SyncLockedBlock()
	//SyncInvalidBlock()
	SyncBlocks()
	SyncMemoryPool()
	SyncMemoryTxState()
}

func RequestBlock(blockID uint64, blockChan chan *rpc.Block, stop chan bool) {

	for {
		rs, ok := rpc.GetBlockById(blockID)
		if !ok {
			t := time.NewTicker(time.Second)
			for {
				select {
				case <-t.C:
					if len(blockChan) == 0 {
						stop <- true
						return
					}
				}
			}
		}
		//fmt.Printf("%s : request id %d\n", time.Now().String(), blockID)
		rs.Id = blockID
		blockChan <- rs
		blockID++
	}
}

func RequestBlockByIds(ids []uint64, blockChan chan *rpc.Block, stop chan bool) {
	for _, id := range ids {
		rs, ok := rpc.GetBlockById(id)
		if ok {
			rs.Id = id
			blockChan <- rs
		}
	}
	t := time.NewTicker(time.Second)
	for {
		select {
		case <-t.C:
			if len(blockChan) == 0 {
				stop <- true
				return
			}
		}
	}
}

func SyncBlocks() {
	iv := time.Duration(conf.Setting.Rpc.SyncBlockInterval)
	defer time.AfterFunc(iv*time.Second, SyncBlocks)
	if s := Stats.GetStatsItem(StatsKey.LatestId); s != nil {
		currId, _ = strconv.ParseUint(string(s), 10, 64)
	}

	LatestHeight, _ = strconv.ParseUint(rpc.GetBlockCount(), 10, 64)
	LatestHeight -= 1
	fmt.Printf("start sync block,last height:%d;current height:%d\n", LatestHeight, currId)

	s := (LatestHeight - currId) / 99
	i := 0
	Storage.Start()
	if currId > 5 {
		currId -= 5
	}
	blockChan := make(chan *rpc.Block, 100)
	stopChan := make(chan bool, 1)
	isEnd := false
	go RequestBlock(currId, blockChan, stopChan)

	for {
		select {
		case <-stopChan:
			for !isEnd {
				time.Sleep(time.Second * 1)
			}
			return
		case block := <-blockChan:
			isEnd = false
			if err := SaveFromBlock(block); err != nil {
				fmt.Printf("save block %d failed! %v, update at next interval", currId, err)
				isEnd = true
				break
			}
			SaveLastBlockId(currId)
			rmSyncedTx(block.Transactions)
			printProgress("Synchronizing blocks", s, &i, currId, LatestHeight)
			//fmt.Printf("%s : save id %d\n", time.Now().String(), currId)
			currId++
			isEnd = true
		}
	}
	fmt.Println()
}

func SyncMemoryPool() {
	fmt.Printf("start sync pool transactions， the time interval is %d\n", conf.Setting.Rpc.SyncMemoryInterval)
	iv := time.Duration(conf.Setting.Rpc.SyncMemoryInterval)
	defer time.AfterFunc(iv*time.Second, SyncMemoryPool)
	trans, err := rpc.GetMemoryPool()
	if err != nil {
		log.Error(err.Error())
	}
	h := len(trans)
	s := h / 99
	i := 0

	for idx, txid := range trans {
		//If it is already in the cache, it is not repeated
		if _, has := notSyncedTxs[txid]; has {
			continue
		}
		tran, err := rpc.GetTransaction(txid)
		if err != nil {
			log.Error(err.Error())
			break
		}
		lock.Lock()
		notSyncedTxs[txid] = *tran
		lock.Unlock()
		printProgress("Synchronizing transactions of pool", uint64(s), &i, uint64(idx+1), uint64(h))
	}
	fmt.Println()
	Storage.Start()
	//defer Storage.End()
	if err := SaveFromMemory(getNotSyncedTxs()); err != nil {
		log.Errorf("save from memory failed! %v, update at next interval", err)
		rmSyncedTxById(trans)
	}
}

func SyncPeerInfo() {
	iv := time.Duration(conf.Setting.Rpc.SyncPeerInfoInterval)
	defer time.AfterFunc(iv*time.Second, SyncPeerInfo)
	fmt.Printf("start sync peer info， the time interval is %d\n", conf.Setting.Rpc.SyncPeerInfoInterval)
	peerInfo, err := rpc.GetPeerInfo()
	if err != nil {
		log.Error(err.Error())
		return
	}
	Storage.Start()
	//defer Storage.End()
	SaveFromPeers(peerInfo)
}

func SyncNodeInfo() {
	var err error
	iv := time.Duration(conf.Setting.Rpc.SyncNodeInfoInterval)
	defer time.AfterFunc(iv*time.Second, SyncNodeInfo)
	fmt.Printf("start sync node info， the time interval is %d\n", conf.Setting.Rpc.SyncPeerInfoInterval)
	NodeInfo, err = rpc.GetNodeInfo()
	if err != nil {
		log.Error(err.Error())
		NodeInfo = new(rpc.NodeInfo)
		NodeInfo.Confirmations = 10
		NodeInfo.Coinbasematurity = 720
		return
	}
}

func SyncLockedBlock() {
	iv := time.Duration(conf.Setting.Rpc.SyncLockedBlockInterval)
	defer time.AfterFunc(iv*time.Second, SyncLockedBlock)
	fmt.Printf("start sync locked block， the time interval is %d\n", conf.Setting.Rpc.SyncLockedBlockInterval)
	Storage.Start()
	ids := Storage.GetLockedMap()
	if len(ids) == 0 {
		return
	}
	s := len(ids) / 99
	i := 0
	var index uint64 = 0
	blockChan := make(chan *rpc.Block, 100)
	stopChan := make(chan bool, 1)
	isEnd := false
	go RequestBlockByIds(ids, blockChan, stopChan)

	for {
		select {
		case <-stopChan:
			for !isEnd {
				time.Sleep(time.Second * 1)
			}
			return
		case block := <-blockChan:
			isEnd = false
			if err := SaveFromBlock(block); err != nil {
				fmt.Printf("save locked block %d failed! %v, update at next interval", block.Id, err)
				isEnd = true
				break
			}
			index++
			printProgress("Synchronizing locked blocks", uint64(s), &i, index, uint64(len(ids)))
			isEnd = true
		}
	}
	fmt.Println()
}

func SyncMemoryTxState() {
	iv := time.Duration(conf.Setting.Rpc.SyncMemoryTxStateInterval)
	defer time.AfterFunc(iv*time.Second, SyncMemoryTxState)
	fmt.Printf("start sync memory transaction state， the time interval is %d\n", conf.Setting.Rpc.SyncMemoryTxStateInterval)
	txs := Storage.GetMemoryTxs()
	var rsTxs UniqueList
	for _, tx := range txs {
		_, err := rpc.GetTransaction(tx.TxId)
		if err != nil {
			if isNoTx(err) {
				tx.Stat = State_Failed
				rsTxs.Add(tx)
			}
		}
	}
	var rsOuts []Output
	for _, e := range rsTxs {
		tx := e.(TxRecord)
		outs := Storage.GetTxOutputs(tx.TxId)
		for _, out := range outs {
			out.Stat = tx.Stat
			rsOuts = append(rsOuts, out)
		}
	}
	if len(rsTxs) > 0 {
		for _, e := range rsTxs {
			tx := e.(TxRecord)
			Storage.ClearTxOutputSpent(tx.TxId)
		}
		Storage.UpdateTxRecord(rsTxs)
		Storage.UpdateOutputsState(rsOuts)
	}
}

func SyncInvalidBlock() {
	iv := time.Duration(conf.Setting.Rpc.SyncInvalidBlockInterval)
	defer time.AfterFunc(iv*time.Second, SyncInvalidBlock)
	fmt.Printf("start sync invalid block， the time interval is %d\n", conf.Setting.Rpc.SyncInvalidBlockInterval)
	Storage.Start()
	ids := Storage.GetInvalidMap()
	if len(ids) == 0 {
		return
	}
	s := len(ids) / 99
	i := 0
	var index uint64 = 0
	blockChan := make(chan *rpc.Block, 100)
	stopChan := make(chan bool, 1)
	isEnd := false
	go RequestBlockByIds(ids, blockChan, stopChan)

	for {
		select {
		case <-stopChan:
			for !isEnd {
				time.Sleep(time.Second * 1)
			}
			return
		case block := <-blockChan:
			isEnd = false
			if err := SaveFromBlock(block); err != nil {
				fmt.Printf("save invalid block %d failed! %v, update at next interval", block.Id, err)
				isEnd = true
				break
			}
			index++
			isEnd = true
			printProgress("Synchronizing invalid blocks", uint64(s), &i, index, uint64(len(ids)))
		}
	}
	fmt.Println()
}

func isNoTx(err error) bool {
	if strings.Contains(err.Error(), "No information available about transaction") {
		return true
	}
	return false
}

func isInMemory(txId string, memTxIds []string) bool {
	for _, memTxid := range memTxIds {
		if memTxid == txId {
			return true
		}
	}
	return false
}

func printProgress(text string, interval uint64, progress *int, last uint64, total uint64) {
	if FirstShow == true && showFlag > 0 {
		return
	}

	if interval == 0 {
		*progress = 100
	} else if last%interval == 0 {
		*progress++
		if *progress > 100 {
			*progress = 100
		}
	}
	fmt.Printf("%s: %d%% [%5d/%d]\r", text, *progress, last, total)
	showFlag++
}

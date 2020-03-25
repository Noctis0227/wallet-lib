package core

import (
	"fmt"
	"github.com/Noctis0227/wallet-lib/rpc"
	sync2 "github.com/Noctis0227/wallet-lib/sync"
	"sync"
	"time"
)

//type Height in

var synchronizer *sync2.Synchronizer
var wg sync.WaitGroup

func StartSync() {
	opt := &sync2.Options{
		RpcUser: "admin",
		RpcPwd:  "123",
		RpcAddr: "127.0.0.1:1234",
		TxChLen: 100,
	}
	historyId := Storage.GetLastId()
	synchronizer = sync2.NewSynchronizer(opt)
	txChan, err := synchronizer.Start(&sync2.HistoryId{historyId.LastTxBlockId, historyId.LastCoinBaseBlockId})
	if err != nil {
		fmt.Printf(err.Error())
		return
	}

	wg.Add(2)

	go saveTxs(txChan)
	go saveHistoryId()

	wg.Wait()
}

func saveTxs(txChan <-chan []rpc.Transaction) {
	for {
		txs := <-txChan
		if txs != nil {
			for range txs {
				// save tx
			}
		}
	}
	wg.Done()
}

func saveHistoryId() {
	for {
		historyId := synchronizer.GetHistoryId()
		if historyId.LastTxBlockId != 0 {
			Storage.UpdateLastId(&HistoryId{
				LastCoinBaseBlockId: historyId.LastCoinBaseBlockId,
				LastTxBlockId:       historyId.LastTxBlockId,
			})
		}

		time.Sleep(time.Second * 10)
	}
	wg.Done()
}

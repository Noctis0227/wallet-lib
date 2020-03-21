package sync

import (
	"fmt"
	"git.diabin.com/BlockChain/wallet-lib/rpc"
	"testing"
)

func TestSynchronizer_Start(t *testing.T) {
	opt := &Options{
		Auth: &rpc.Auth{
			User: "admin",
			Pwd:  "123",
		},
		Host:    "http://127.0.0.1:1234",
		TxChLen: 100,
	}
	sync := NewSynchronizer(opt)
	txChan, err := sync.Start(&HistoryId{200, 101})
	if err != nil {
		t.Error(err.Error())
	}
	for i := 0; i < 100; i++ {
		select {
		case txs := <-txChan:
			for _, tx := range txs {
				fmt.Printf("%s %d\n", tx.Txhash, tx.Confirmations)
			}
		}
	}
	historyId := sync.GetHistoryId()
	sync.Stop()
	fmt.Println(historyId.LastTxBlockId, historyId.LastCoinBaseBlockId)
}

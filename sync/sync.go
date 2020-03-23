package sync

import (
	"fmt"
	"git.diabin.com/BlockChain/wallet-lib/rpc"
	"time"
)

const (
	defaultHost                 = "127.0.0.1:1234"
	defaultTxChLen              = 100
	defaultRepeatCount          = 5
	defaultCoinBaseThreshold    = 720
	defaultTransactionThreshold = 10
)

type Synchronizer struct {
	rpcClient          *rpc.Client
	opt                *Options
	threshold          *threshold
	TxChannel          chan []rpc.Transaction
	stopSyncTxCh       chan bool
	stopSyncCoinBaseCh chan bool
	curTxBlockId       uint64
	curCoinBaseBlockId uint64
}

type Options struct {
	// Rpc option
	RpcAddr string
	RpcUser string
	RpcPwd  string

	// tx channel length
	TxChLen uint
}

type HistoryId struct {
	LastTxBlockId       uint64
	LastCoinBaseBlockId uint64
}

func NewSynchronizer(opt *Options) *Synchronizer {
	if opt.RpcAddr == "" {
		opt.RpcAddr = defaultHost
	}
	if opt.TxChLen == 0 {
		opt.TxChLen = defaultTxChLen
	}

	client := rpc.NewClient(&rpc.RpcConfig{
		Address: opt.RpcAddr,
		User:    opt.RpcUser,
		Pwd:     opt.RpcPwd,
	})
	return &Synchronizer{
		rpcClient:          client,
		opt:                opt,
		TxChannel:          make(chan []rpc.Transaction, opt.TxChLen),
		stopSyncTxCh:       make(chan bool),
		stopSyncCoinBaseCh: make(chan bool),
		threshold: &threshold{
			coinBaseThreshold:    defaultCoinBaseThreshold,
			transactionThreshold: defaultTransactionThreshold,
		},
	}
}

// start syncing at 0
// or start syncing at last stop return id
func (s *Synchronizer) Start(id *HistoryId) (<-chan []rpc.Transaction, error) {
	if err := s.setThreshold(); err != nil {
		return nil, fmt.Errorf("failed to set threshold %s", err.Error())
	}

	go s.startSync(id)

	return s.TxChannel, nil
}

// use the return value as the parameter for the next startup
func (s *Synchronizer) Stop() {
	s.stopSyncTxCh <- true
	s.stopSyncCoinBaseCh <- true
}

func (s *Synchronizer) GetHistoryId() *HistoryId {
	return &HistoryId{
		LastTxBlockId:       s.curTxBlockId,
		LastCoinBaseBlockId: s.curCoinBaseBlockId,
	}
}

func (s *Synchronizer) startSync(hisId *HistoryId) {
	s.curTxBlockId = hisId.LastTxBlockId
	if s.curTxBlockId >= defaultRepeatCount {
		s.curTxBlockId -= defaultRepeatCount
	} else {
		s.curTxBlockId = 0
	}

	s.curCoinBaseBlockId = hisId.LastCoinBaseBlockId
	if s.curCoinBaseBlockId >= defaultRepeatCount {
		s.curCoinBaseBlockId -= defaultRepeatCount
	} else {
		s.curCoinBaseBlockId = 0
	}

	go s.SyncTxs()
	go s.SyncCoinBaseTx()
}

func (s *Synchronizer) SyncTxs() {
	s.requestTxs()
}

func (s *Synchronizer) SyncCoinBaseTx() {
	for {
		select {
		case _ = <-s.stopSyncCoinBaseCh:
			return
		default:
			block, err := s.rpcClient.GetBlockById(s.curCoinBaseBlockId)
			if err != nil {
				time.Sleep(time.Second * 5)
				break
			}
			if !s.isBlockConfirmed(block) {
				time.Sleep(time.Second * 1)
				break
			}
			if usable, err := s.IsCoinBaseUsable(block); err != nil {
				time.Sleep(time.Second * 5)
				break
			} else {
				if usable {
					txs := getConfirmedCoinBase(block)
					if len(txs) != 0 {
						s.TxChannel <- txs
					}
				}
				s.curCoinBaseBlockId++
			}
		}
	}
}

func (s *Synchronizer) requestTxs() {
	for {
		select {
		case _ = <-s.stopSyncTxCh:
			return
		default:
			block, err := s.rpcClient.GetBlockById(s.curTxBlockId)
			if err != nil {
				time.Sleep(time.Second * 5)
				break
			}
			if block.Txsvalid {
				if s.isTxConfirmed(block) {
					txs := getConfirmedTx(block)
					if len(txs) != 0 {
						s.TxChannel <- txs
					}
					s.curTxBlockId++
				} else {
					time.Sleep(time.Second * 1)
				}
			} else {
				if s.isTxConfirmed(block) {
					s.curTxBlockId++
				} else {
					time.Sleep(time.Second * 1)
				}
			}
		}
	}
}

func (s *Synchronizer) isBlockConfirmed(block *rpc.Block) bool {
	return block.Confirmations >= s.threshold.coinBaseThreshold
}

func (s *Synchronizer) isTxConfirmed(block *rpc.Block) bool {
	return block.Confirmations >= s.threshold.transactionThreshold
}

func (s *Synchronizer) IsCoinBaseUsable(block *rpc.Block) (bool, error) {
	color, err := s.rpcClient.IsBlue(block.Hash)
	if err != nil {
		return false, err
	}
	switch color {
	case 0:
		return false, nil
	case 1:
		return true, nil
	}
	return false, nil
}

type threshold struct {
	coinBaseThreshold    uint32
	transactionThreshold uint32
}

func (s *Synchronizer) setThreshold() error {
	nodeInfo, err := s.rpcClient.GetNodeInfo()
	if err != nil {
		return err
	}
	s.threshold.coinBaseThreshold = nodeInfo.Coinbasematurity
	s.threshold.transactionThreshold = nodeInfo.Confirmations
	return nil
}

func getConfirmedTx(block *rpc.Block) []rpc.Transaction {
	for i, tx := range block.Transactions {
		if isCoinBase(&tx) {
			return append(block.Transactions[0:i], block.Transactions[i+1:]...)
		}
	}
	return []rpc.Transaction{}
}

func getConfirmedCoinBase(block *rpc.Block) []rpc.Transaction {
	for _, tx := range block.Transactions {
		if isCoinBase(&tx) {
			return []rpc.Transaction{tx}
		}
	}
	return []rpc.Transaction{}
}

func isCoinBase(tx *rpc.Transaction) bool {
	if tx != nil && len(tx.Vin) > 0 && tx.Vin[0].Coinbase != "" {
		return true
	}
	return false
}

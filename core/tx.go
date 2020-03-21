package core

import (
	"errors"
	"fmt"
	"git.diabin.com/BlockChain/wallet-lib/conf"
	"git.diabin.com/BlockChain/wallet-lib/log"
	"git.diabin.com/BlockChain/wallet-lib/rpc"
	"github.com/Qitmeer/qitmeer/core/types"
	"github.com/Qitmeer/qitmeer/qx"
	"strings"
)

const (
	Trading_Speed_Slow    = 10
	Trading_Speed_Faster  = 100
	Trading_Speed_Fastest = 1000
)

const (
	State_Confirmed   = 0
	State_Unconfirmed = 1
	State_InMemory    = 2
	State_Invaild     = 3
	State_Failed      = 4
	State_Red         = 5
)

const MinAmount = 5970

func SendTransaction(fromAddr string, key string, toAddr string, amount float64, speed float64) (string, bool) {
	if len(key) < 1 {
		log.Debug("key is required")
		return "send transaction failed! key is required", false
	}
	if len(fromAddr) < 1 || len(toAddr) < 1 || amount <= 0 {
		log.Debugf("from:%s,to:%s,amount:%f", fromAddr, toAddr, amount)
		return "send transaction failed! fromAddr,key,toAddr,amount is required", false
	}

	if speed > 10000 || speed < 1 {
		return "the coefficient of trading speed is not in the range(1, 10000)", false
	}

	strSign, stat := SignedTx(fromAddr, key, toAddr, amount, 0.00, speed)
	if !stat {
		return strSign, stat
	}

	log.Info("signed string:", strSign)
	txId, err := RpcClient.SendTransaction(strSign)
	if err != nil {
		stat = false
	}
	return strings.Trim(txId, "\""), stat
}

func SignedTx(fromAddr string, key string, toAddr string, amount float64, fees uint64, speed float64) (string, bool) {
	utxo, usable := GetUsableOuts(fromAddr)
	if utxo == nil || 0 == len(utxo) {
		return "There is not enough balance", false
	}

	sum := uint64(0)
	iUsable, err := NewAmount(usable)
	if err != nil {
		return "Fail to create the currency amount", false
	}
	iAmount, err := NewAmount(amount)
	if err != nil {
		return "Fail to create the currency amount", false
	}

	strSign, _, stat := SignEncode(utxo, fromAddr, key, toAddr, iAmount, fees)
	if !stat {
		return strSign, stat
	}
	fees = GetFees(strSign, speed)
	for {
		strSign, sum, stat = SignEncode(utxo, fromAddr, key, toAddr, iAmount, fees)
		if !stat {
			return strSign, stat
		}

		fees = GetFees(strSign, speed)
		if sum >= iAmount+fees || (sum < iAmount+fees && iUsable >= iAmount) {
			break
		}
	}
	return strSign, stat
}

func SignEncode(utxo []Output, fromAddr string, key string, toAddr string, iAmount uint64, fees uint64) (string, uint64, bool) {
	strEncode, sum, stat := EncodeTransaction(utxo, fromAddr, toAddr, iAmount, fees)
	if !stat {
		return strEncode, sum, stat
	}
	strSign, stat := SignTransaction(strEncode, key)
	if !stat {
		return strEncode, sum, stat
	}
	if len(strSign)%2 != 0 {
		strSign = "0" + strSign
	}
	return strSign, sum, true
}

func EncodeTransaction(utxo []Output, fromAddr string, toAddr string, iAmount uint64, fees uint64) (string, uint64, bool) {
	if iAmount < fees {
		return fmt.Sprintf("The transfer amount is less than the service fee %f", float64(fees)/conf.Ratio), 0, false
	}
	if iAmount < MinAmount {
		return fmt.Sprintf("The minimum transfer amount shall not be less than %f after the fee is deducted", float64(MinAmount)/conf.Ratio), 0, false
	}

	iNeedAmount := iAmount + fees
	inputs := make(map[string]uint32)
	sum := uint64(0)
	for _, it := range utxo {
		sum += it.Amount
		inputs[it.TxId] = uint32(it.OutIndex)
		if sum >= iNeedAmount {
			break
		}
	}

	outputs := map[string]uint64{toAddr: iAmount}
	if sum < iNeedAmount && sum > iAmount {
		outputs[toAddr] = sum - fees
	} else if sum > iNeedAmount {
		change := sum - iNeedAmount
		outputs[fromAddr] = change
	} else if sum < iNeedAmount {
		return fmt.Sprintf("There is not enough balance, fees is %d", fees), sum, false
	}

	encode, err := qx.TxEncode(1, 0, inputs, outputs)
	if err != nil {
		return err.Error(), sum, false
	}
	return strings.Trim(encode, "\""), sum, true
}

func SignTransaction(encode string, key string) (string, bool) {

	rs, err := qx.TxSign(key, encode, conf.Setting.Version)
	if err != nil {
		return err.Error(), false
	}
	return rs, true
}

func GetFees(strSign string, speed float64) uint64 {
	return uint64(float64(len(strSign)) / 2 / 1000 * 10000 * speed)
}

func GetLatestTxRecords(idx int, size int) []TxRecord {
	if idx == 0 {
		idx = 1
	}
	if size == 0 {
		size = 50
	}
	return List.QueryLatestTx(idx, size)
}

func GetLatestTxRecordOuts(txid string, idx int, size int) []Output {
	if idx == 0 {
		idx = 1
	}
	if size == 0 {
		size = 50
	}
	return List.QueryLatestRecordOuts(txid, idx, size)
}

func GetTxRecord(txId string) *TxRecord {
	txRecord := Storage.GetTxRecord(txId)
	if txRecord == nil {
		return nil
	}

	return txRecord
}

func GetTransaction(id string) (*rpc.Transaction, error) {
	if len(id) < 1 {
		return nil, errors.New("query transaction failed! txid is required")
	}
	return RpcClient.GetTransaction(id)
}

func GetAddrRecord(addr string, idx, size int) []TxRecord {
	if idx == 0 {
		idx = 1
	}
	if size == 0 {
		size = 50
	}
	return List.QueryAddrLastestTx(addr, idx, size)
}

func NewAmount(amount float64) (uint64, error) {
	iAmount, err := types.NewAmount(amount)
	if err != nil {
		return 0, err
	}
	return uint64(iAmount), nil
}

func TxRecordToJsonTxRecord(txRecord *TxRecord) *JsonTxRecord {
	outs := Storage.GetTxOutputs(txRecord.TxId)
	jsonOuts := make([]*JsonOutRecord, 0)
	index := 0
	for i, out := range outs {
		if out.Address == txRecord.From {
			index = i
		}
		jsonOuts = append(jsonOuts, &JsonOutRecord{
			OutIndex:  out.OutIndex,
			Address:   out.Address,
			Amount:    out.Amount,
			SpentTxId: out.SpentTxId,
			IsChange:  false,
		})
	}
	if len(outs) > 1 {
		jsonOuts[index].IsChange = true
	}

	jsonTx := &JsonTxRecord{
		Id:         txRecord.Id,
		TxId:       txRecord.TxId,
		From:       txRecord.From,
		Amount:     txRecord.Amount,
		ToCount:    txRecord.ToCount,
		Spent:      txRecord.Spent,
		Change:     txRecord.Change,
		Fee:        txRecord.Fee,
		BlockOrder: txRecord.BlockOrder,
		InTime:     txRecord.InTime,
		Stat:       txRecord.Stat,
		Outs:       jsonOuts,
	}
	return jsonTx
}

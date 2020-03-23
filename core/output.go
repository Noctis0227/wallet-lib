package core

import "github.com/Qitmeer/qitmeer/core/types"

type IQueryOutput interface {
	GetUsable(addr string) ([]Output, float64)
	GetUnspent(addr string) ([]Output, float64)
	GetSpent(addr string) ([]Output, float64)
	GetMemoryUnspent(addr string) ([]Output, float64)
	GetMemorySpent(addr string) ([]Output, float64)
	GetLockUnspent(addr string) ([]Output, float64)
	GetBalance(addr string) ([]Output, float64)
	GetAllOuts(addr string) ([]Output, float64)
	GetOutput(key string) *Output
}

type QueryOutput struct {
}

func (m *QueryOutput) GetUsable(addr string) ([]Output, float64) {
	condition := "address = ? and spent_tx_id = ? and stat = ?"
	outs := List.QueryOuts(condition, []interface{}{addr, "", 0})
	return outs, getOutAmount(outs)
}

func (m *QueryOutput) GetUnspent(addr string) ([]Output, float64) {
	condition := "address = ? and spent_tx_id = ? and stat <> ?"
	outs := List.QueryOuts(condition, []interface{}{addr, "", 3})
	return outs, getOutAmount(outs)
}

func (m *QueryOutput) GetSpent(addr string) ([]Output, float64) {
	condition := "address = ? and spent_tx_id <> ?"
	outs := List.QueryOuts(condition, []interface{}{addr, ""})
	return outs, getOutAmount(outs)
}

func (m *QueryOutput) GetMemoryUnspent(addr string) ([]Output, float64) {
	condition := "address = ? and spent_tx_id = ? and stat = ?"
	outs := List.QueryOuts(condition, []interface{}{addr, "", 2})
	return outs, getOutAmount(outs)
}

func (m *QueryOutput) GetMemorySpent(addr string) ([]Output, float64) {
	condition := "address = ? and spent_tx_id <> ? and stat = ?"
	outs := List.QueryOuts(condition, []interface{}{addr, "", 2})
	return outs, getOutAmount(outs)
}

func (m *QueryOutput) GetLockUnspent(addr string) ([]Output, float64) {
	condition := "address = ? and spent_tx_id = ? and stat = ?"
	outs := List.QueryOuts(condition, []interface{}{addr, "", 1})
	return outs, getOutAmount(outs)
}

func (m *QueryOutput) GetBalance(addr string) ([]Output, float64) {
	condition := "address = ? and spent_tx_id = ? and stat < ?"
	outs := List.QueryOuts(condition, []interface{}{addr, "", 3})
	return outs, getOutAmount(outs)
}

func (m *QueryOutput) GetAllOuts(addr string) ([]Output, float64) {
	condition := "address = ?"
	outs := List.QueryOuts(condition, []interface{}{addr})
	return outs, getOutAmount(outs)
}

func (m *QueryOutput) GetOutput(key string) *Output {
	return Storage.GetOutput(key)
}

func getOutAmount(outs []Output) float64 {
	var total uint64
	for _, out := range outs {
		total += out.Amount
	}
	return types.Amount(total).ToCoin()
}

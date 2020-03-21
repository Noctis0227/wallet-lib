package output

import (
	"git.diabin.com/BlockChain/wallet-lib/core"
	"github.com/Qitmeer/qitmeer/core/types"
)

type MysqlOutput struct {
}

func (m *MysqlOutput) GetUsable(addr string) ([]core.Output, float64) {
	condition := "address = ? and spent_tx_id = ? and stat = ?"
	outs := core.List.QueryOuts(condition, []interface{}{addr, "", 0})
	return outs, getOutAmount(outs)
}

func (m *MysqlOutput) GetUnspent(addr string) ([]core.Output, float64) {
	condition := "address = ? and spent_tx_id = ? and stat <> ?"
	outs := core.List.QueryOuts(condition, []interface{}{addr, "", 3})
	return outs, getOutAmount(outs)
}

func (m *MysqlOutput) GetSpent(addr string) ([]core.Output, float64) {
	condition := "address = ? and spent_tx_id <> ?"
	outs := core.List.QueryOuts(condition, []interface{}{addr, ""})
	return outs, getOutAmount(outs)
}

func (m *MysqlOutput) GetMemoryUnspent(addr string) ([]core.Output, float64) {
	condition := "address = ? and spent_tx_id = ? and stat = ?"
	outs := core.List.QueryOuts(condition, []interface{}{addr, "", 2})
	return outs, getOutAmount(outs)
}

func (m *MysqlOutput) GetMemorySpent(addr string) ([]core.Output, float64) {
	condition := "address = ? and spent_tx_id <> ? and stat = ?"
	outs := core.List.QueryOuts(condition, []interface{}{addr, "", 2})
	return outs, getOutAmount(outs)
}

func (m *MysqlOutput) GetLockUnspent(addr string) ([]core.Output, float64) {
	condition := "address = ? and spent_tx_id = ? and stat = ?"
	outs := core.List.QueryOuts(condition, []interface{}{addr, "", 1})
	return outs, getOutAmount(outs)
}

func (m *MysqlOutput) GetBalance(addr string) ([]core.Output, float64) {
	condition := "address = ? and spent_tx_id = ? and stat < ?"
	outs := core.List.QueryOuts(condition, []interface{}{addr, "", 3})
	return outs, getOutAmount(outs)
}

func (m *MysqlOutput) GetAllOuts(addr string) ([]core.Output, float64) {
	condition := "address = ?"
	outs := core.List.QueryOuts(condition, []interface{}{addr})
	return outs, getOutAmount(outs)
}

func (m *MysqlOutput) GetOutput(key string) *core.Output {
	return core.Storage.GetOutput(key)
}

func getOutAmount(outs []core.Output) float64 {
	var total uint64
	for _, out := range outs {
		total += out.Amount
	}
	return types.Amount(total).ToCoin()
}

package sqlite

import (
	"fmt"
	"git.diabin.com/BlockChain/wallet-lib/core"
	"git.diabin.com/BlockChain/wallet-lib/log"
	"strings"
)

var mainHeight uint64 = 0

// IStorage
func (m *Sqlite) Start() {
}

func (m *Sqlite) End() {
}

func (m *Sqlite) Init() {
}

func (m *Sqlite) UpdateOutputs(outs []core.Output) error {
	session := m.engine.NewSession()
	defer session.Close()

	sql := "call insertOrUpdateOutput('%s', %d, '%s', %d, %d, '%s', %d, '%s', '%s', %d)"
	for _, out := range outs {
		_, err := session.Exec(fmt.Sprintf(sql, out.TxId, out.OutIndex, out.BlockHash, out.BlockOrder, out.BlockId,
			out.Address, out.Amount, out.Timestamp, out.SpentTxId, out.Stat))
		if err != nil {
			log.Errorf("update output %v failed! %s", out, err)
			return err
		}
	}
	return nil
}

func (m *Sqlite) UpdateOutputsState(outs []core.Output) error {
	session := m.engine.NewSession()
	defer session.Close()

	for _, out := range outs {
		m.engine.Id(out.Id).Update(out)
	}
	return nil
}

func (m *Sqlite) UpdateTxOutputSpent(txId string, outIndex int, spentTxId string) error {
	session := m.engine.NewSession()
	defer session.Close()

	_, err := m.engine.Table(new(core.Output)).Where("tx_id = ? and out_index = ?", txId, outIndex).Update(map[string]string{"spent_tx_id": spentTxId})
	return err
}

func (m *Sqlite) UpdateLockedMap(id uint64, stat int8) error {
	return nil
}

func (m *Sqlite) UpdateInvalidMap(id uint64, stat int8) error {
	return nil
}

func (m *Sqlite) UpdateMainHeight(height uint64) error {
	mainHeight = height
	return nil
}

func (m *Sqlite) ClearTxOutputSpent(spentId string) error {
	session := m.engine.NewSession()
	defer session.Close()

	_, err := m.engine.Table(new(core.Output)).Where("spent_tx_id = ?", spentId).Update(map[string]string{"spent_tx_id": ""})
	return err
}

func (m *Sqlite) GetOutput(key string) *core.Output {
	out := new(core.Output)
	keyList := strings.Split(key, "-")
	if len(keyList) != 2 {
		return nil
	}
	session := m.engine.Where("tx_id=? and out_index=? and stat < ?", keyList[0], keyList[1], core.State_Invaild)
	defer session.Close()

	exist, err := session.Get(out)
	if err != nil || !exist {
		return nil
	}
	return out
}

func (m *Sqlite) GetTxOutputs(txId string) []core.Output {
	return m.QueryOuts("tx_id = ? and stat < ?", []interface{}{txId, core.State_Invaild})
}

func (m *Sqlite) QueryLatestRecordOuts(txid string, idx int, size int) []core.Output {
	idx -= 1
	outs := make([]core.Output, 0)
	start := idx * size
	session := m.engine.Where("tx_id = ? and stat < ?", txid, core.State_Invaild).Desc("id").Limit(size, start)
	defer session.Close()

	session.Find(&outs)
	return outs
}

// IStat
func (m *Sqlite) AddLatestTx(txId string) {
	// TODO
}

func (m *Sqlite) UpdateStats(key []byte, val string) {
}

func (m *Sqlite) ClearDatabase() {
	for _, v := range m.engine.Tables {
		err := m.engine.DropTables(v.Name)
		if err != nil {
			log.Errorf("clear database error! %s", err)
		}
	}
}

func (m *Sqlite) GetLastId() *core.HistoryId {
	lastId := &core.HistoryId{}
	m.engine.Table(new(core.HistoryId)).Where("id = ?", 1).Get(lastId)
	return lastId
}

func (m *Sqlite) UpdateLastId(lastId *core.HistoryId) {
	lastId.Id = 1
	_, err := m.engine.Insert(lastId)
	if err != nil {
		m.engine.Table(new(core.HistoryId)).Where("id = ? ", 1).
			Update(map[string]uint64{"last_tx_block_id": lastId.LastTxBlockId, "last_coin_base_block_id": lastId.LastCoinBaseBlockId})
	}

}

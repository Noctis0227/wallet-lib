package sqlite

import (
	"git.diabin.com/BlockChain/wallet-lib/core"
	"github.com/go-xorm/xorm"
	_ "github.com/mattn/go-sqlite3"
	xormcore "xorm.io/core"
)

type Sqlite struct {
	engine *xorm.Engine
}

func (m *Sqlite) Connect(connect string) error {
	var err error
	m.engine, err = xorm.NewEngine("sqlite", connect)
	if err != nil {
		return err
	}
	return nil
}

func (m *Sqlite) Close() error {
	return m.engine.Close()
}

func (m *Sqlite) Open(connect string) error {
	err := m.Connect(connect)
	if err != nil {
		return err
	}
	m.engine.ShowSQL(false)
	tbMapper := xormcore.NewPrefixMapper(xormcore.SameMapper{}, "")
	m.engine.SetTableMapper(tbMapper)

	err = m.engine.Sync(new(core.Output), new(core.HistoryId))
	if err != nil {
		return err
	}

	return nil
}

func (m *Sqlite) QueryOuts(condition interface{}, value []interface{}) []core.Output {
	outs := make([]core.Output, 0)
	session := m.engine.Where(condition, value...)
	defer session.Close()

	session.Find(&outs)
	return outs
}

package output

import (
	"git.diabin.com/BlockChain/wallet-lib/core"
	"github.com/Qitmeer/qitmeer/core/types"
)

type BoltOutput struct {
}

func (b *BoltOutput) GetUsable(addr string) ([]core.Output, float64) {
	return b.getOuts(addr, 1)
}

func (b *BoltOutput) GetUnspent(addr string) ([]core.Output, float64) {
	return b.getOuts(addr, 2)
}

func (b *BoltOutput) GetSpent(addr string) ([]core.Output, float64) {
	return b.getOuts(addr, 3)
}

func (b *BoltOutput) GetMemoryUnspent(addr string) ([]core.Output, float64) {
	return b.getOuts(addr, 4)
}

func (b *BoltOutput) GetMemorySpent(addr string) ([]core.Output, float64) {
	return b.getOuts(addr, 5)
}

func (b *BoltOutput) GetLockUnspent(addr string) ([]core.Output, float64) {
	return b.getOuts(addr, 6)
}

func (b *BoltOutput) GetBalance(addr string) ([]core.Output, float64) {
	return b.getOuts(addr, 7)
}

func (b *BoltOutput) GetAllOuts(addr string) ([]core.Output, float64) {
	return b.getOuts(addr, 0)
}

func (b *BoltOutput) GetOutput(key string) *core.Output {
	return core.Storage.GetOutput(key)
}

func (b *BoltOutput) getOuts(addr string, typ uint8) ([]core.Output, float64) {
	var outs core.UniqueList
	total := float64(0)
	if item := core.Storage.GetAddressMap(addr); item != nil {
		var rs []core.Output
		switch typ {
		case 0:
			outs = item.Outs
			break
		case 1:
			outs = splitUnspentOuts(item.Unspents, false)
			break
		case 2:
			outs = item.Unspents
			outs = append(outs, item.MemoryUnspents...)
			break
		case 3:
			outs = item.Spents
			outs = append(outs, item.MemorySpents...)
			break
		case 4:
			outs = item.MemoryUnspents
			break
		case 5:
			outs = item.MemorySpents
			break
		case 6:
			outs = splitUnspentOuts(item.Unspents, true)
			break
		case 7:
			outs = item.Unspents
			break
		}

		for _, key := range outs {
			if it := core.Storage.GetOutput(key.(string)); it != nil {
				rs = append(rs, *it)
				total += types.Amount(it.Amount).ToCoin()
			}
		}
		return rs, total
	}
	return nil, total
}

func splitUnspentOuts(unSpentOuts core.UniqueList, locked bool) core.UniqueList {
	var outs core.UniqueList
	for _, key := range unSpentOuts {
		if out := core.Storage.GetOutput(key.(string)); out != nil {
			if locked {
				if out.Stat == 1 {
					outs = append(outs, key)
				}
			} else {
				if out.Stat == 0 {
					outs = append(outs, key)
				}
			}
		}
	}
	return outs
}

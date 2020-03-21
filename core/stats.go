package core

import "github.com/Qitmeer/qitmeer/core/types"

func GetAddressAmounts(addr string) *AddressInfo {
	_, usableAmount := GetUsableOuts(addr)
	_, lockedAmount := GetLockedUnspent(addr)
	_, memoryUnspentAmount := GetMemoryUnspent(addr)
	_, spentAmount := GetSpentOuts(addr)
	var rs = AddressInfo{
		UseableAmount: usableAmount,
		LockedAmount:  lockedAmount,
		MemoryAmount:  memoryUnspentAmount,
		SpentAmount:   spentAmount,
	}
	//_, rs.Balance = GetBalance(addr)
	u, _ := types.NewAmount(usableAmount)
	l, _ := types.NewAmount(lockedAmount)
	m, _ := types.NewAmount(memoryUnspentAmount)
	balance := u + l + m
	rs.Balance = balance.ToCoin()
	return &rs
}

func GetAddressInfo(addr string) *AddressInfo {
	_, usableAmount := GetUsableOuts(addr)
	_, lockedAmount := GetLockedUnspent(addr)
	_, memoryUnspentAmount := GetMemoryUnspent(addr)
	_, spentAmount := GetSpentOuts(addr)
	var rs = AddressInfo{
		UseableAmount: usableAmount,
		LockedAmount:  lockedAmount,
		MemoryAmount:  memoryUnspentAmount,
		SpentAmount:   spentAmount,
	}
	//_, rs.Balance = GetBalance(addr)
	u, _ := types.NewAmount(usableAmount)
	l, _ := types.NewAmount(lockedAmount)
	m, _ := types.NewAmount(memoryUnspentAmount)
	balance := u + l + m
	rs.Balance = balance.ToCoin()
	txs := GetAddrRecord(addr, 1, 20)
	rs.TxRecords = txs
	return &rs
}

func GetLatestBlocks(idx int, size int) []BlockRecord {
	if idx == 0 {
		idx = 1
	}
	if size == 0 {
		size = 50
	}

	return List.QueryLatestBlocks(idx, size)
}

func GetLatestPeers(idx int, size int) []PeerRecord {
	if idx == 0 {
		idx = 1
	}
	if size == 0 {
		size = 50
	}

	return List.QueryLatestPeers(idx, size)
}

func GetChainStats() map[string]string {
	return Stats.GetStats()
}

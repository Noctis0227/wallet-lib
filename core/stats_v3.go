package core

func GetAddressInfoV3(addr string) *AddressInfo {
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
	rs.Balance = usableAmount + lockedAmount + memoryUnspentAmount
	txs := GetAddrRecordV3(addr, 1, 20)
	rs.TxRecords = txs
	return &rs
}

package core

func GetLatestMiners(idx int, size int) []MinerRecord {
	if idx == 0 {
		idx = 1
	}
	if size == 0 {
		size = 50
	}

	return List.QueryLatestMiners(idx, size)
}

func GetMinerRecordByHeight(height uint64) *MinerRecord {
	return Storage.GetMinerRecordByHeight(height)
}

func GetMinerRecordByAddr(addr string) *MinerRecord {
	return Storage.GetMinerRecordByAddr(addr)
}

func GetMinerRecords(iStart, iEnd uint64) []map[string]string {
	return Storage.GetMinerRecords(iStart, iEnd)
}

func GetMinerBlocksV3(address string, start, end uint64) []map[string]string {
	return Storage.GeMinersBlocks(address, start, end)
}

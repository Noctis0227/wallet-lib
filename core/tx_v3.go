package core

func GetLatestTxRecordsV3(idx int, size int) []*JsonTxRecord {
	if idx == 0 {
		idx = 1
	}
	if size == 0 {
		size = 50
	}
	txRecords := List.QueryLatestTx(idx, size)
	var jsonTxRecords = make([]*JsonTxRecord, 0)
	for _, tx := range txRecords {
		jsonTxRecords = append(jsonTxRecords, TxRecordToJsonTxRecord(&tx))
	}
	return jsonTxRecords
}

func GetAddrRecordV3(addr string, idx, size int) []*JsonTxRecord {
	if idx == 0 {
		idx = 1
	}
	if size == 0 {
		size = 50
	}
	txRecords := List.QueryAddrLastestTx(addr, idx, size)
	var jsonTxRecords = make([]*JsonTxRecord, 0)
	for _, tx := range txRecords {
		jsonTxRecords = append(jsonTxRecords, TxRecordToJsonTxRecord(&tx))
	}
	return jsonTxRecords
}

func GetAddrRecordCountV3(addr string) int64 {
	return List.QueryAddrRecordCount(addr)
}

func GetTxRecordV3(txId string) *JsonTxRecord {
	txRecord := Storage.GetTxRecord(txId)
	if txRecord == nil {
		return nil
	}

	return TxRecordToJsonTxRecord(txRecord)
}

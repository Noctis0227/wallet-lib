package core

import (
	"strconv"
)

func GetBlockRecordByHash(hash string) *BlockRecord {
	return Storage.GetBlockRecordByHash(hash)
}

func GetBlockRecordByHeight(height string) *BlockRecord {
	uHeight, err := strconv.ParseUint(height, 10, 64)
	if err != nil {
		return nil
	}
	block, err := RpcClient.GetBlockByOrder(uHeight)
	if err != nil {
		return nil
	}
	return GetBlockRecordByHash(block.Hash)
}

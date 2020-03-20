package core

import (
	"kahf/rpc"
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
	block, ok := rpc.GetBlock(uHeight)
	if !ok {
		return nil
	}
	return GetBlockRecordByHash(block.Hash)
}

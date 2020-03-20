package rpc

import (
	"encoding/json"
	"kahf/log"
	"testing"
)

func TestGetBlockCount(t *testing.T) {
	rs := GetBlockCount()
	log.Debug(rs)
}

func TestGetBlock(t *testing.T) {
	b, ok := GetBlock(12)
	if ok {
		rs, _ := json.Marshal(b)
		log.Debugf("%s", rs)
	} else {
		log.Debugf("failed")
	}
}

func TestGetMemoryPool(t *testing.T) {
	ls, err := GetMemoryPool()
	if err == nil {
		rs, _ := json.Marshal(ls)
		log.Debugf("%s", rs)
	} else {
		log.Debugf("failed")
	}
}

func TestSendTransaction(t *testing.T) {
	ls, err := GetMemoryPool()
	if err == nil {
		rs, _ := json.Marshal(ls)
		log.Debugf("%s", rs)
	} else {
		log.Debugf("failed")
	}
}

func TestGetTransaction(t *testing.T) {
	ls, err := GetTransaction("")
	if err == nil {
		rs, _ := json.Marshal(ls)
		log.Debugf("%s", rs)
	} else {
		log.Debugf("failed")
	}
}

package rpc

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"errors"
	"git.diabin.com/BlockChain/wallet-lib/conf"
	"git.diabin.com/BlockChain/wallet-lib/log"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
)

func GetBlock(h uint64) (*Block, bool) {
	params := []interface{}{h, true}
	resp := NewReqeust(params).SetMethod("getBlockByOrder").call()
	blk := new(Block)
	if resp.Error != nil {
		//log.Error(resp.Error.Message)
		return blk, false
	}
	if err := json.Unmarshal(resp.Result, blk); err != nil {
		log.Error(err.Error())
		return blk, false
	}
	return blk, true
}

func GetBlockCount() string {
	var params []interface{}
	resp := NewReqeust(params).SetMethod("getBlockCount").call()
	if resp.Error != nil {
		return "-1"
	}
	return string(resp.Result)
}

func SendTransaction(tx string) (string, bool) {
	params := []interface{}{strings.Trim(tx, "\n"), false}
	resp := NewReqeust(params).SetMethod("sendRawTransaction").call()
	if resp.Error != nil {
		log.Errorf("send raw transaction failed! %s", resp.Error.Message)
		return resp.Error.Message, false
	}
	txId := string(resp.Result)
	return txId, true
}

func GetTransaction(txid string) (*Transaction, error) {
	params := []interface{}{txid, true}
	resp := NewReqeust(params).SetMethod("getRawTransaction").call()
	if resp.Error != nil {
		return nil, errors.New(resp.Error.Message)
	}
	var rs *Transaction
	if err := json.Unmarshal(resp.Result, &rs); err != nil {
		return nil, err
	}
	return rs, nil
}

func CreateTransaction(inputs []TransactionInput, amounts Amounts) (string, error) {
	jsonInput, err := json.Marshal(inputs)
	if err != nil {
		return "", err
	}
	jsonAmount, err := json.Marshal(amounts)
	if err != nil {
		return "", err
	}
	params := []interface{}{json.RawMessage(jsonInput), json.RawMessage(jsonAmount)}
	resp := NewReqeust(params).SetMethod("createRawTransaction").call()
	if resp.Error != nil {
		return "", errors.New(resp.Error.Message)
	}
	encode := string(resp.Result)
	return encode, nil
}

func GetMemoryPool() ([]string, error) {
	params := []interface{}{"", false}
	resp := NewReqeust(params).SetMethod("getMempool").call()
	if resp.Error != nil {
		return nil, errors.New(resp.Error.Message)
	}
	var rs []string
	if err := json.Unmarshal(resp.Result, &rs); err != nil {
		return nil, err
	}
	return rs, nil
}

func GetPeerInfo() ([]PeerInfo, error) {
	var params []interface{}
	resp := NewReqeust(params).SetMethod("getPeerInfo").call()
	if resp.Error != nil {
		return nil, errors.New(resp.Error.Message)
	}
	var rs []PeerInfo
	if err := json.Unmarshal(resp.Result, &rs); err != nil {
		return nil, err
	}
	return rs, nil
}

func GetBlockById(id uint64) (*Block, bool) {
	params := []interface{}{id, true}
	resp := NewReqeust(params).SetMethod("getBlockByID").call()
	blk := new(Block)
	if resp.Error != nil {
		return blk, false
	}
	if err := json.Unmarshal(resp.Result, blk); err != nil {
		log.Error(err.Error())
		return blk, false
	}
	return blk, true
}

func GetNodeInfo() (*NodeInfo, error) {
	params := []interface{}{}
	resp := NewReqeust(params).SetMethod("getNodeInfo").call()
	nodeInfo := new(NodeInfo)
	if resp.Error != nil {
		return nodeInfo, errors.New(resp.Error.Message)
	}
	if err := json.Unmarshal(resp.Result, nodeInfo); err != nil {
		log.Error(err.Error())
		return nodeInfo, err
	}
	return nodeInfo, nil
}

func IsBlue(hash string) (int, error) {
	params := []interface{}{hash}
	resp := NewReqeust(params).SetMethod("isBlue").call()
	if resp.Error != nil {
		return 0, errors.New(resp.Error.Message)
	}
	state, err := strconv.Atoi(string(resp.Result))
	if err != nil {
		return 0, err
	}
	return state, nil
}

func (req *ClientRequest) call() *ClientResponse {
	cfg := conf.Setting.Rpc
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	client := &http.Client{Transport: tr}

	//convert struct to []byte
	marshaledData, err := json.Marshal(req)
	if err != nil {
		log.Error(err.Error())
	}
	//log.Debugf("rpc call starting, Host:%s, params:%s", cfg.Host, marshaledData)

	httpRequest, err :=
		http.NewRequest(http.MethodPost, cfg.Host, bytes.NewReader(marshaledData))
	if err != nil {
		log.Error(err.Error())
	}

	if httpRequest == nil {
		log.Error("the httpRequest is nil")
		return &ClientResponse{Error: &Error{Message: "the httpRequest is nil"}}
	}
	httpRequest.Close = true
	httpRequest.Header.Set("Content-Type", "application/json")
	httpRequest.SetBasicAuth(conf.Setting.Rpc.User, conf.Setting.Rpc.Pwd)
	//log.Debugf("u:%s;p:%s", cfg.User, cfg.Pwd)

	response, err := client.Do(httpRequest)
	if err != nil {
		log.Error(err.Error())
		return &ClientResponse{Error: &Error{Message: err.Error()}}
	}

	body := response.Body

	bodyBytes, err := ioutil.ReadAll(body)
	if err != nil {
		log.Error("io read error", err.Error())
	}

	//log.Info("rpc call successful! ", string(bodyBytes))

	resp := &ClientResponse{}
	//convert []byte to struct
	if err := json.Unmarshal(bodyBytes, resp); err != nil {
		log.Errorf("json unmarshal failed; value:%s; error:%s", string(bodyBytes), err.Error())
	}

	err = response.Body.Close()
	if err != nil {
		log.Error(err.Error())
	}

	if resp.Error != nil {
		//log.Fail(resp.Error.Message)
	}
	return resp
}

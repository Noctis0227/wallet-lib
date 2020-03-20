package core

import (
	"bytes"
	"encoding/gob"
	"encoding/hex"
	"fmt"
	"kahf/log"
	"testing"
	"unsafe"
)

func TestStackList_Push(t *testing.T) {
	ls := StringStack{IgnoreRepeat: true, MaxSize: 3}
	ls.Push("1")
	ls.Push("2")
	ls.Push("3")
	ls.Push("4")

	ls.Each(func(i int, it string) {
		fmt.Printf("index:%d\titem:%s\n", i, it)
	})
}

func TestUniqueList_Add(t *testing.T) {
	ls := UniqueList{}
	ls.Add("test1")
	ls.Add("test2")
	ls.Add("test2")

	for _, it := range ls {
		println(it.(string))
	}

}
func TestBytes_ToOutput(t *testing.T) {
	out := Output{
		BlockOrder: 10,
		BlockHash:  "12333333",
	}

	fmt.Println(unsafe.Sizeof(out) / 8 / 1024)

	var bs Bytes
	bs = out.ToBytes()
	fmt.Printf("to bytes success,%s\n", bs)
	//rs:=&Output{}
	//fromBytes(bs,rs)
	rs := bs.ToOutput()
	fmt.Printf("%s,%d", rs.BlockHash, rs.BlockOrder)
}

func TestBytes_Convert(t *testing.T) {

	s := StringStack{MaxSize: 10}
	s.Push("test 1")
	s.Push("test 2")
	s.Push("test 3")
	var result bytes.Buffer
	encoder := gob.NewEncoder(&result)
	err := encoder.Encode(s)
	if err != nil {
		log.Errorf("An error occurred when to bytes,%s", err)
	}
	bs := result.Bytes()

	r := StringStack{}
	decoder := gob.NewDecoder(bytes.NewReader(bs))
	e := decoder.Decode(&r)
	if e != nil {
		log.Errorf("An error occurred when from bytes,%s", err)
	}
}

func TestGenerateAddr(t *testing.T) {
	//LoadKeystore()
	//GenerateAddr()
	//LoadKeystore()
	ks := GetAddrs()
	if ks != nil {
		//fmt.Println("def addr:", ks.DefAddr)
		for k, v := range ks {
			fmt.Printf("pk: %s\t addr:%s\n", k, v)
		}
	}
	//LoadKeystore()
	//if ks != nil {
	//	for k, v := range ks {
	//		fmt.Printf("pk: %s\t addr:%s\n", k, v)
	//	}
	//}
}

func TestHexString(t *testing.T) {
	ks := readData()

	s := hex.EncodeToString(toBytes(ks))
	fmt.Println(s)

	bs, err := hex.DecodeString(s)
	if err != nil {
		fmt.Printf("err,%s", err)
		return
	}

	rs := map[string]string{}
	fromBytes(bs, &rs)
	for k, v := range rs {
		fmt.Printf("%s,%s\n", k, v)
	}
}

func TestGetBlockStats(t *testing.T) {
	tf := "%-20s%-20s%-20s%-20s\n"
	fmt.Printf(tf, "BlockNum", "TxNum", "UnconfirmedTxNum", "LatestHeight")
	bf := "%-20d%-20d%-20d%-20d\n"
	rs := GetChainStats()
	if rs != nil {
		fmt.Printf(bf, rs[string(StatsKey.BlockNum)], rs[string(StatsKey.TxNum)], rs[string(StatsKey.UnconfirmedTxNum)], rs[string(StatsKey.LatestHeight)])
	}
}

func TestArrayLength(t *testing.T) {
	//a:=[]int{10,20,30,40,50}
	//for i,it:=range a[3:6]{
	//	println(i,it)
	//}
	println(len("47ceb894a86a2fa765029c5ede037152c8b1e7d14207be8c538322955cf243eb"))
	println(len("0b5de3e2d81f52d05b3c9956ce704b05e61f7752e943b9323dcfbff0a6cb5739"))

	tmp := StringStack{}
	tmp.Push("10")
	tmp.Push("20")
	tmp.Push("30")
	tmp.Push("40")
	tmp.Push("50")
	tmp.EachByPage(2, 3, func(i int, it string) {
		println(i, it)
	})
}

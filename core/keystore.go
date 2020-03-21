package core

import (
	"bufio"
	"encoding/hex"
	"fmt"
	"git.diabin.com/BlockChain/wallet-lib/conf"
	"git.diabin.com/BlockChain/wallet-lib/log"
	"github.com/Qitmeer/qitmeer/qx"
	"io"
	"os"
)

const ksFile = ".keystore"

//const ksFile = "/Users/logan/Workspace/github/forks/kahf/.Keystore"

type Keystore struct {
	AddrKeyMap map[string]string
	DefAddr    string
}

var kstore *Keystore

func LoadKeystore() {
	kstore = readData()
}

func GenerateAddr() string {
	ver := conf.Setting.Version

	s, err := qx.NewEntropy(32)
	if err != nil {
		log.Errorf("An error occurred generating s，%s", err)
	}
	prv, err := qx.EcNew("secp256k1", s)
	if err != nil {
		log.Errorf("An error occurred generating private key，%s", err)
	}
	pub, err := qx.EcPrivateKeyToEcPublicKey(false, prv)
	if err != nil {
		log.Errorf("An error occurred generating public key，%s", err)
	}
	addr, err := qx.EcPubKeyToAddress(ver, pub)
	if err != nil {
		log.Errorf("An error occurred generating address，%s", err)
	}

	if kstore == nil {
		kstore = &Keystore{AddrKeyMap: map[string]string{}}
	}
	kstore.AddrKeyMap[addr] = prv
	kstore.DefAddr = addr
	saveKeystore()
	return addr
}

func GetPK(addr string) string {
	if kstore == nil {
		return ""
	}
	if it, ok := kstore.AddrKeyMap[addr]; ok {
		return it
	}
	return ""
}

func GetDefAddr() string {
	if kstore == nil {
		return ""
	}
	return kstore.DefAddr
}

func GetAddrs() map[string]string {
	if kstore == nil {
		return nil
	}
	return kstore.AddrKeyMap
}

func SetDefAddr(addr string) error {
	if kstore == nil {
		return fmt.Errorf("addr '%s' has not be generated", addr)
	}

	if _, ok := kstore.AddrKeyMap[addr]; ok {
		kstore.DefAddr = addr
	}
	return fmt.Errorf("addr '%s' has not be generated", addr)
}

func readData() *Keystore {
	f, err := os.Open(ksFile)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		log.Error("read data file failed", err.Error())
		return nil
	}
	if f != nil {
		defer f.Close()
	}

	reader := bufio.NewReader(f)
	var str string
	for {
		//_, err := reader.Read(buf[:])
		str, err = reader.ReadString('\n')
		if err == io.EOF {
			//fmt.Println("read the file finished")
			break
		}
		if err != nil {
			log.Errorf("%s", err)
			os.Exit(2)
		}
		//fmt.Println(string(buf[:n]))
	}
	rs, _ := hex.DecodeString(str)
	ks := Keystore{}
	fromBytes(rs, &ks)
	return &ks
}

func saveKeystore() {
	f, err := os.Create(ksFile)
	if err != nil {
		log.Error("", err.Error())
		return
	}
	if f != nil {
		defer f.Close()
	}

	s := hex.EncodeToString(toBytes(kstore))
	if _, err := f.WriteString(s); err != nil {
		log.Errorf("write Keystore failed!")
	}
}

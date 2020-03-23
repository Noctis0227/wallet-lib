package sign

import (
	"github.com/Qitmeer/qitmeer/qx"
)

func TxEncode(version uint32, lockTime uint32, inputs map[string]uint32, outputs map[string]uint64) (string, error) {
	return qx.TxEncode(version, lockTime, inputs, outputs)
}

func TxSign(encode string, key string, network string) (string, bool) {
	rs, err := qx.TxSign(key, encode, network)
	if err != nil {
		return err.Error(), false
	}
	return rs, true
}

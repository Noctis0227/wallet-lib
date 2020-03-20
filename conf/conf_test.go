package conf

import (
	"encoding/base32"
	"fmt"
	"testing"
)

func TestSettings(t *testing.T) {
	//dist:=[]byte("hlc")
	src := []byte("this is test fadsfas asdfasd fasdfas faasdfdas ")
	encoding := base32.NewEncoding("qwertyuiopasdfghjklzxcvbnm123456")
	rs := encoding.EncodeToString(src)
	fmt.Println(rs)
	s, _ := encoding.DecodeString(rs)
	fmt.Println(string(s))
}

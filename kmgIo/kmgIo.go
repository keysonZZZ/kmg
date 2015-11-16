package kmgIo

import (
	"io"
	"io/ioutil"
)

func MustReadAll(r io.Reader) (b []byte) {
	b, err := ioutil.ReadAll(r)
	if err != nil {
		panic(err)
	}
	return b
}

// 将所有内容都读出来,但是将结果全部扔掉,只返回读取的字节数量
var DiscardReadFrom = ioutil.Discard.(io.ReaderFrom).ReadFrom
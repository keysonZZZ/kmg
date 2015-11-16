package kmgRand

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"io"
	"sync"
)

//这个东西的首要目标就是快,内容看上去很随机而已 有 200M/s 左右
var FastRandReader io.Reader = &fastRandReader{}

const randBlockSize = 256

type fastRandReader struct {
	stream  cipher.Stream
	buf     [randBlockSize]byte
	lock    sync.Mutex
	hasInit bool
}

func (r *fastRandReader) Read(dst []byte) (n int, err error) {
	r.lock.Lock()
	defer r.lock.Unlock()
	r.init()
	remainSize := len(dst)
	for {
		if remainSize >= randBlockSize {
			r.stream.XORKeyStream(dst[remainSize-randBlockSize:remainSize], r.buf[:])
			remainSize -= randBlockSize
			continue
		}
		r.stream.XORKeyStream(dst[0:remainSize], r.buf[:remainSize])
		break
	}
	return len(dst), nil
}

func (r *fastRandReader) init() {
	if r.stream != nil {
		return
	}
	_, err := io.ReadFull(rand.Reader, r.buf[:])
	if err != nil {
		panic(err)
	}
	block, err := aes.NewCipher(r.buf[:32])
	if err != nil {
		panic(err)
	}
	r.stream = cipher.NewCTR(block, r.buf[32:48])
}

//这个东西的首要目标就是快,内容看上去很随机而已 现在速度有 200M/s 左右
func NewLimitedFastRandReader(size int) io.Reader{
	return io.LimitReader(FastRandReader,int64(size))
	//return newRandSizeReader(size)
}

/*
const oneBlockSize = 1024 * 1024

func randSizeWrite(w io.Writer, size int) (err error) {
	remainByteNum := size
	n := 0
	for {
		if remainByteNum <= 0 {
			return
		}
		if remainByteNum < oneBlockSize {
			n, err = w.Write(rand1M[:remainByteNum])
			remainByteNum -= n
			if err != nil {
				return
			}
		} else {
			n, err = w.Write(rand1M)
			remainByteNum -= n
			if err != nil {
				return
			}
		}
	}
}

func newRandSizeReader(size int) (reader io.Reader) {
	randInit()
	r, w := io.Pipe()
	go func() {
		defer w.Close()
		err := randSizeWrite(w, size)
		w.CloseWithError(err)
	}()
	return r
	//return &randReader{remainByteNum: size}
}

type randReader struct{
	remainByteNum int
	pos int
}

func (r *randReader) Read(dst []byte) (n int, err error) {
	if r.remainByteNum<=0{
		return n,io.EOF
	}
	writeByteNum:=len(dst)
	if writeByteNum>r.remainByteNum{
		writeByteNum = r.remainByteNum
	}
	if writeByteNum>len(rand1M){
		writeByteNum = len(rand1M)
	}
	copy(dst[:writeByteNum],rand1M[:writeByteNum])
	r.remainByteNum -= writeByteNum
	if r.remainByteNum<=0{
		return writeByteNum,io.EOF
	}
	return writeByteNum,nil
}

var rand1M []byte
var once sync.Once

func randInit() {
	once.Do(func(){
		rand1M = make([]byte, oneBlockSize)
		//明文风险
		_, err := rand.Read(rand1M[:])
		if err != nil {
			panic(err)
		}
	})
}
*/
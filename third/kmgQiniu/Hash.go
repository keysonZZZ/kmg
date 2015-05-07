package kmgQiniu

import (
	"bytes"
	"crypto/sha1"
	"encoding/base64"
	"io"
	"os"
)

const (
	hash_BLOCK_BITS = 22 // Indicate that the blocksize is 4M
	hash_BLOCK_SIZE = 1 << hash_BLOCK_BITS
)

func hashBlockCount(fsize int64) int {

	return int((fsize + (hash_BLOCK_SIZE - 1)) >> hash_BLOCK_BITS)
}

func calSha1(b []byte, r io.Reader) ([]byte, error) {

	h := sha1.New()
	_, err := io.Copy(h, r)
	if err != nil {
		return nil, err
	}
	return h.Sum(b), nil
}

//计算从文件计算七牛hash值
func ComputeHashFromFile(filename string) (etag string, err error) {
	f, err := os.Open(filename)
	if err != nil {
		return
	}
	defer f.Close()

	fi, err := f.Stat()
	if err != nil {
		return
	}

	fsize := fi.Size()
	blockCnt := hashBlockCount(fsize)
	sha1Buf := make([]byte, 0, 21)

	if blockCnt <= 1 { // file size <= 4M
		sha1Buf = append(sha1Buf, 0x16)
		sha1Buf, err = calSha1(sha1Buf, f)
		if err != nil {
			return
		}
	} else { // file size > 4M
		sha1Buf = append(sha1Buf, 0x96)
		sha1BlockBuf := make([]byte, 0, blockCnt*20)
		for i := 0; i < blockCnt; i++ {
			body := io.LimitReader(f, hash_BLOCK_SIZE)
			sha1BlockBuf, err = calSha1(sha1BlockBuf, body)
			if err != nil {
				return
			}
		}
		sha1Buf, _ = calSha1(sha1Buf, bytes.NewReader(sha1BlockBuf))
	}
	etag = base64.URLEncoding.EncodeToString(sha1Buf)
	return
}

func ComputeHashFromBytes(b []byte) (etag string) {
	f := bytes.NewReader(b)
	blockCnt := hashBlockCount(int64(len(b)))
	sha1Buf := make([]byte, 0, 21)

	var err error
	if blockCnt <= 1 { // file size <= 4M
		sha1Buf = append(sha1Buf, 0x16)
		sha1Buf, err = calSha1(sha1Buf, f)
		if err != nil {
			panic(err)
		}
	} else { // file size > 4M
		sha1Buf = append(sha1Buf, 0x96)
		sha1BlockBuf := make([]byte, 0, blockCnt*20)
		for i := 0; i < blockCnt; i++ {
			body := io.LimitReader(f, hash_BLOCK_SIZE)
			sha1BlockBuf, err = calSha1(sha1BlockBuf, body)
			if err != nil {
				panic(err)
			}
		}
		sha1Buf, _ = calSha1(sha1Buf, bytes.NewReader(sha1BlockBuf))
	}
	etag = base64.URLEncoding.EncodeToString(sha1Buf)
	return
}

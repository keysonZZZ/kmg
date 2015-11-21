package kmgHttp

import (
	"errors"
	"flag"
	"fmt"
	"github.com/bronze1man/kmg/kmgCache"
	"github.com/bronze1man/kmg/kmgConsole"
	"github.com/bronze1man/kmg/kmgCrypto"
	"github.com/bronze1man/kmg/kmgErr"
	"github.com/bronze1man/kmg/kmgFile"
	"github.com/bronze1man/kmg/kmgStrings"
	"net/http"
	"os"
	"path/filepath"
	"sync"
)

func AddCommandList() {
	kmgConsole.AddAction(kmgConsole.Command{
		Name:   "FileHttpServer",
		Runner: runFileHttpServer,
	})
	kmgConsole.AddCommandWithName("HttpGet", func() {
		requestUrl := ""
		key := ""
		flag.StringVar(&requestUrl, "url", "", "")
		flag.StringVar(&key, "key", "", "crypto key use to decrypt respond")
		flag.Parse()
		if requestUrl == "" {
			kmgConsole.ExitOnErr(errors.New("Usage: kmg HttpGet -url http://xxx"))
		}
		b := MustUrlGetContent(requestUrl)
		var err error
		if key != "" {
			b, err = kmgCrypto.CompressAndEncryptBytesDecodeV2(kmgCrypto.Get32PskFromString(key), b)
			if err != nil {
				panic(err)
			}
		}
		fmt.Print(string(b))
	})
}

var lock *sync.Mutex = &sync.Mutex{}
var cacheFilePathSlice []string = []string{}
var cacheFilePathEncryptMap map[string][]byte = map[string][]byte{}

func runFileHttpServer() {
	listenAddr := ""
	path := ""
	key := ""
	flag.StringVar(&listenAddr, "l", ":80", "listen address")
	flag.StringVar(&path, "path", "", "root path of the file server")
	flag.StringVar(&key, "key", "", "crypto key use to encrypt all request of this server")
	flag.Parse()
	var err error
	if path == "" {
		path, err = os.Getwd()
		if err != nil {
			fmt.Printf("os.Getwd() fail %s", err)
			return
		}
	} else {
		kmgErr.PanicIfError(os.Chdir(path))
	}
	if key == "" {
		http.Handle("/", http.FileServer(http.Dir(path)))
	}
	if key != "" {
		http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			realPath := filepath.Join(path, r.URL.Path)
			if !kmgFile.MustFileExist(realPath) {
				w.Write([]byte("File Not Exist"))
				return
			}
			if !kmgStrings.IsInSlice(cacheFilePathSlice, realPath) {
				cacheFilePathSlice = append(cacheFilePathSlice, realPath)
			}
			updateCache := func() {
				cacheFilePathEncryptMap[realPath] = kmgCrypto.CompressAndEncryptBytesEncodeV2(
					kmgCrypto.Get32PskFromString(key),
					kmgFile.MustReadFile(realPath),
				)
			}
			checkCache := func() {
				lock.Lock()
				defer lock.Unlock()
				kmgCache.MustMd5FileChangeCache(realPath, []string{realPath}, func() {
					updateCache()
				})
			}
			checkCache()
			//进程重启后，内存中的缓存掉了，但是文件系统的缓存还在
			_, exist := cacheFilePathEncryptMap[realPath]
			if !exist {
				updateCache()
			}
			w.Write(cacheFilePathEncryptMap[realPath])
		})
	}
	fmt.Println("start server at", listenAddr)
	err = http.ListenAndServe(listenAddr, nil)
	if err != nil {
		fmt.Printf("http.ListenAndServe() fail %s", err)
		return
	}
	return
}

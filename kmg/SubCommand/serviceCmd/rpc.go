package serviceCmd

import (
	"fmt"
	"github.com/bronze1man/kmg/errors"
	"github.com/bronze1man/kmg/kmgCrypto"
	"github.com/bronze1man/kmg/kmgErr"
	"github.com/bronze1man/kmg/kmgProcessMutex"
	"time"
)

type StartStatus string

const (
	StartStatusSuccess StartStatus = "StartStatusSuccess"
	StartStatusFail                = "StartStatusFail"
)

var statusChannel = make(chan StartStatus)

var rpcPsk = kmgCrypto.Get32PskFromString("w4n4ts28cq")

var rpcAddress = "127.0.0.1:2777"

type ServiceRpc struct{}

func (sr *ServiceRpc) Send(status StartStatus) {
	statusChannel <- status
}

//TODO 对于不会发 RPC 的进程，应该可以用一个 flag 来表明不需要等待 RPC 响应
func waitRpcRespond() chan error {
	returnChan := make(chan error)
	lock := &kmgProcessMutex.FileMutex{Name: "kmg_service_lock"}
	lock.Lock()
	rpcCloser := ListenAndServe_ServiceRpc(rpcAddress, &ServiceRpc{}, rpcPsk)
	__closer := func() {
		err := rpcCloser()
		time.Sleep(time.Second)
		lock.UnLock()
		kmgErr.PanicIfError(err)
	}
	go func() {
		select {
		case startStatus := <-statusChannel:
			__closer()
			if startStatus == StartStatusSuccess {
				returnChan <- nil
			} else {
				returnChan <- errors.New("StartFail")
			}
		case <-time.After(time.Minute * 2):
			__closer()
			returnChan <- errors.New("Rpc timeout")
		}
	}()
	return returnChan
}

func ServiceStartSuccess() {
	time.Sleep(time.Millisecond * 100)
	client := NewClient_ServiceRpc("http://"+rpcAddress, rpcPsk)
	err := client.Send(StartStatusSuccess)
	if err != nil {
		fmt.Println("[Warning kmg service] please use kmg service start/restart", err)
	}
}

//暂时没什么用
func ServiceStartFail() {
	time.Sleep(time.Second)
	client := NewClient_ServiceRpc("http://"+rpcAddress, rpcPsk)
	err := client.Send(StartStatusFail)
	if err != nil {
		fmt.Println("src/github.com/bronze1man/kmg/kmg/SubCommand/serviceCmd/rpc.go ServiceStartFail", err)
	}
}

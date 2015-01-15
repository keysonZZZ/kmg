package netTester

import (
	"github.com/bronze1man/kmg/kmgNet"
	"github.com/bronze1man/kmg/kmgTime"
	"net"
	"time"
)

func writeOnly(listenerNewer ListenNewer, Dialer DirectDialer, debug bool) {
	toWrite := []byte("hello world")
	listener := listenAccept(listenerNewer, func(c net.Conn) {
		defer c.Close()
		for i := 0; i < 2; i++ {
			mustReadSame(c, toWrite)
			time.Sleep(time.Microsecond)
		}
	})
	defer listener.Close()
	kmgTime.MustNotTimeout(func() {
		conn1, err := Dialer()
		mustNotError(err)
		if debug {
			conn1 = kmgNet.NewDebugConn(conn1, "writeOnly")
		}
		defer conn1.Close()
		for i := 0; i < 2; i++ {
			_, err = conn1.Write(toWrite)
			mustNotError(err)
			time.Sleep(time.Microsecond)
		}
		conn1.Close()
	}, time.Second)

	listener.Close()
}

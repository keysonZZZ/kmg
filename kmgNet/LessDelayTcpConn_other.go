// +build !linux

package kmgNet

import "net"

func LessDelayTcpConn(conn *net.TCPConn) (connOut net.Conn, err error) {
	return conn, nil
}

var LessDelayDial = net.Dial
var MustLessDelayListen = MustListen

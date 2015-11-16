package kmgNet

import (
	"fmt"
	"net"
	"strconv"
	"strings"
	"time"
)

//从一个net.Listener里面读取需要Dial的地址(测试用的比较多)
func MustGetLocalAddrFromListener(listener net.Listener) string {
	return MustGetLocalAddrFromAddr(listener.Addr())
}

//从一个net.Listener里面读取需要Dial的地址(测试用的比较多)
func MustGetLocalAddrFromAddr(addr net.Addr) string {
	tcpAddr, err := net.ResolveTCPAddr(addr.Network(), addr.String())
	if err != nil {
		panic(err)
	}
	return "127.0.0.1:" + strconv.Itoa(tcpAddr.Port)
}

func MustTcpRandomListen() net.Listener {
	l, err := net.Listen("tcp", ":0")
	if err != nil {
		panic(err)
	}
	return l
}

func SpeedString(byteNum int, dur time.Duration) string {
	bytePerSecond := float64(byteNum) / (float64(dur) / float64(time.Second))
	return SpeedStringWithFloat(bytePerSecond)
}

func SpeedStringWithFloat(bytePerSecond float64) string {
	if bytePerSecond > 1e9 {
		return fmt.Sprintf("%.2fGB/s", bytePerSecond/(1024*1024*1024))
	}
	if bytePerSecond > 1e6 {
		return fmt.Sprintf("%.2fMB/s", bytePerSecond/(1024*1024))
	}
	if bytePerSecond > 1e3 {
		return fmt.Sprintf("%.2fKB/s", bytePerSecond/1024)
	}
	return fmt.Sprintf("%.2fB/s", bytePerSecond)
}

func SizeString(byteNum int64) string {
	if byteNum > 1e15 || byteNum < -1e15 {
		return fmt.Sprintf("%.2fPB", float64(byteNum)/(1024*1024*1024*1024*1024))
	}
	if byteNum > 1e12 || byteNum < -1e12 {
		return fmt.Sprintf("%.2fTB", float64(byteNum)/(1024*1024*1024*1024))
	}
	if byteNum > 1e9 || byteNum < -1e9 {
		return fmt.Sprintf("%.2fGB", float64(byteNum)/(1024*1024*1024))
	}
	if byteNum > 1e6 || byteNum < -1e6 {
		return fmt.Sprintf("%.2fMB", float64(byteNum)/(1024*1024))
	}
	if byteNum > 1e3 || byteNum < -1e3 {
		return fmt.Sprintf("%.2fKB", float64(byteNum)/(1024))
	}
	return fmt.Sprintf("%dB", byteNum)
}

// 在开头加padding,尝试使长度一致,如果数据超级大有可能会坏掉
func SizeStringWithPadding(byteNum int64) string {
	s := SizeString(byteNum)
	if len(s) < 10 {
		return strings.Repeat(" ", 10-len(s)) + s
	}
	return s
}

func CloseRead(conn net.Conn) error {
	tcpC := mustGetTcpConnFromConn(conn)
	return tcpC.CloseRead()
}

func CloseWrite(conn net.Conn) error {
	tcpC := mustGetTcpConnFromConn(conn)
	return tcpC.CloseWrite()
}

func mustGetTcpConnFromConn(conn net.Conn) *net.TCPConn {
	tcpC, ok := conn.(*net.TCPConn)
	if ok {
		return tcpC
	}
	conner, ok := conn.(GetUnderlyingConner)
	if ok {
		return mustGetTcpConnFromConn(conner.GetUnderlyingConn())
	}
	panic(fmt.Errorf("not support conn type %T", conn))
}

type GetUnderlyingConner interface {
	GetUnderlyingConn() net.Conn
}

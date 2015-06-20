package kmgRadius

import (
	"github.com/bronze1man/kmg/kmgLog"
	"io"
	"net"
)

type AcctRequest struct {
	SessionId   string //连接id
	Username    string
	SessionTime uint32 //连接时间
	InputBytes  uint64 //流入字节
	OutputBytes uint64 //流出字节
	NasPort     uint32
}

type Handler struct {
	// 有的协议需要明文密码来做各种hash的事情
	// exist返回false可以踢掉客户端
	Auth func(username string) (password string, exist bool)
	// 计费开始,来了一条新连接
	// 根据协议规定 此处返回给客户端的包,不能发送任何有效信息(比如踢掉客户端,请采用其他办法踢掉客户端)
	AcctStart func(acctReq AcctRequest)
	// 计费数据更新
	// 根据协议规定 此处返回给客户端的包,不能发送任何有效信息(比如踢掉客户端,请采用其他办法踢掉客户端)
	AcctUpdate func(acctReq AcctRequest)
	// 计费结束
	// 根据协议规定 此处返回给客户端的包,不能发送任何有效信息(比如踢掉客户端,请采用其他办法踢掉客户端)
	AcctStop func(acctReq AcctRequest)
}

//异步运行服务器,
// TODO 返回Closer以便可以关闭服务器,所有无法运行的错误panic出来,其他错误丢到kmgLog error里面.
// 如果不需要Closer可以直接忽略
func RunServer(address string, secret []byte, handler Handler) io.Closer {
	s := server{
		mschapMap: map[string]mschapStatus{},
		handler:   handler,
	}
	return RunServerWithPacketHandler(address, secret, s.PacketHandler)
}

type PacketHandler func(request *Packet) *Packet

//异步运行服务器,返回Closer以便可以关闭服务器,所有无法运行的错误panic出来,其他错误丢到kmgLog error里面.
func RunServerWithPacketHandler(address string, secret []byte, handler PacketHandler) io.Closer {
	var conn *net.UDPConn
	go func() {
		addr, err := net.ResolveUDPAddr("udp", address)
		if err != nil {
			panic(err)
		}
		conn, err = net.ListenUDP("udp", addr)
		if err != nil {
			panic(err)
		}

		for {
			b := make([]byte, 4096)
			n, senderAddress, err := conn.ReadFrom(b)
			if err != nil {
				panic(err)
			}
			go func(p []byte, senderAddress net.Addr) {
				pac, err := DecodeRequestPacket(secret, p)
				if err != nil {
					kmgLog.Log("error", "radius.Decode", err.Error())
					return
				}

				npac := handler(pac)
				if npac == nil {
					// 特殊情况,返回nil,表示抛弃这个包.
					return
				}
				err = npac.Send(conn, senderAddress)
				if err != nil {
					kmgLog.Log("error", "radius.Send", err.Error())
					return
				}
			}(b[:n], senderAddress)
		}
		return
	}()
	return conn
}

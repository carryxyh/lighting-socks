package server

import (
	"lighting-socks/core"
	"log"
	"net"
)

type LsServer struct {
	*core.SecureSocket
}

// 新建一个服务端
func New(password *core.Password, listenAddr *net.TCPAddr) *LsServer {
	return &LsServer{
		SecureSocket: &core.SecureSocket{
			Cipher:     core.NewCipher(password),
			ListenAddr: listenAddr,
		},
	}
}

// 运行服务端并且监听来自本地代理客户端的请求
func (lsServer *LsServer) Listen(didListen func(listenAddr net.Addr)) error {
	listener, err := net.ListenTCP("tcp", lsServer.ListenAddr)
	if err != nil {
		return err
	}
	defer listener.Close()

	if didListen != nil {
		didListen(listener.Addr())
	}

	for {
		localConn, err := listener.AcceptTCP()
		if err != nil {
			log.Println(err)
			continue
		}

		// localConn被关闭时直接清除所有数据 不管没有发送的数据
		localConn.SetLinger(0)
		go lsServer.handleConn(localConn)
	}

	return nil
}

func (lsServer *LsServer) handleConn(localConn *net.TCPConn) {

}

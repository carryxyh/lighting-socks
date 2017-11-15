package local

import (
	"lighting-socks/core"
	"log"
	"net"
)

type LsLocal struct {
	*core.SecureSocket
}

// 新建一个本地端
func New(password *core.Password, listenAddr, remoteAddr *net.TCPAddr) *LsLocal {
	return &LsLocal{
		SecureSocket: &core.SecureSocket{
			Cipher:     core.NewCipher(password),
			ListenAddr: listenAddr,
			RemoteAddr: remoteAddr,
		},
	}
}

// 本地端启动监听，接收来自本机浏览器的连接
func (lsLocal *LsLocal) Listen(didListen func(listenAddr net.Addr)) error {
	listener, err := net.ListenTCP("tcp", lsLocal.ListenAddr)
	if err != nil {
		return err
	}
	defer listener.Close()

	if didListen != nil {
		didListen(listener.Addr())
	}

	for {
		userConn, err := listener.AcceptTCP()
		if err != nil {
			log.Println(err)
			continue
		}
		// userConn被关闭时直接清除所有数据 不管没有发送的数据
		userConn.SetLinger(0)
		go lsLocal.handleConn(userConn)
	}
	return nil
}

func (local *LsLocal) handleConn(userConn *net.TCPConn) {
	defer userConn.Close()

	proxyServer, err := local.ConnRemote()
	if err != nil {
		log.Println(err)
		return
	}
	defer proxyServer.Close()
	// Conn被关闭时直接清除所有数据 不管没有发送的数据
	proxyServer.SetLinger(0)

	// 进行转发 从 proxyServer 读取数据发送到 localUser
	go func() {
		err := local.DecodeCopy(userConn, proxyServer)
		if err != nil {
			userConn.Close()
			proxyServer.Close()
		}
	}()
}

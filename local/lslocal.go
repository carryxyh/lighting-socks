package local

import (
	"lighting-socks/core"
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

	}
	return nil
}

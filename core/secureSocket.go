package core

import (
	"errors"
	"fmt"
	"io"
	"net"
)

const (
	BufSize = 1024
)

// 加密传输的 TCP Socket
type SecureSocket struct {
	Cipher     *Cipher
	ListenAddr *net.TCPAddr
	RemoteAddr *net.TCPAddr
}

// 从输入流里读取加密过的数据，解密后把原数据放到bs里
func (secureSocket *SecureSocket) DecodeRead(conn *net.TCPConn, bs []byte) (n int, err error) {
	n, err = conn.Read(bs)
	if err != nil {
		return
	}
	secureSocket.Cipher.decode(bs[:n])
	return
}

// 把放在bs里的数据加密后立即全部写入输出流
func (secureSocket *SecureSocket) EncodeWrite(conn *net.TCPConn, bs []byte) (n int, err error) {
	secureSocket.Cipher.encode(bs)
	return conn.Write(bs)
}

// 从src中源源不断的读取原数据加密后写入到dst，直到src中没有数据可以再读取
func (secureSocket *SecureSocket) EncodeCopy(dst *net.TCPConn, src *net.TCPConn) error {
	buf := make([]byte, BufSize)
	for {
		readCount, readErr := src.Read(buf)
		if readErr != nil {
			if readErr != io.EOF {
				return readErr
			} else {
				return nil
			}
		}
		if readCount > 0 {
			writeCount, writeErr := dst.Write(buf)
			if writeErr != nil {
				return writeErr
			}
			if readCount != writeCount {
				return io.ErrShortWrite
			}
		}
	}
}

// 从src中源源不断的读取加密后的数据解密后写入到dst，直到src中没有数据可以再读取
func (secureSocket *SecureSocket) DecodeCopy(dst *net.TCPConn, src *net.TCPConn) error {
	buf := make([]byte, BufSize)
	for {
		readCount, readErr := secureSocket.DecodeRead(src, buf)
		if readErr != nil {
			if readErr != nil {
				if readErr != io.EOF {
					return readErr
				} else {
					return nil
				}
			}
		}
		if readCount > 0 {
			writeCount, writeErr := dst.Write(buf[0:readCount])
			if writeErr != nil {
				return writeErr
			}
			if readCount != writeCount {
				return io.ErrShortWrite
			}
		}
	}
}

// 和远程的socket建立连接，他们之间的数据传输会加密
func (secureSocket *SecureSocket) ConnRemote() (conn *net.TCPConn, err error) {
	conn, err = net.DialTCP("tcp", nil, secureSocket.RemoteAddr)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("连接到远程服务器 %s 失败:%s", secureSocket.RemoteAddr, err))
	}
	return
}

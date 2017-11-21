package core

type Cipher struct {

	// 编码用的密码
	encodePwd *Password

	// 解码用的密码
	decodePwd *Password
}

// 加密原数据
func (cipher *Cipher) encode(bs []byte) {
	for i, v := range bs {
		bs[i] = cipher.encodePwd[v]
	}
}

// 解码加密后的数据到原数据
func (cipher *Cipher) decode(bs []byte) {
	for i, v := range bs {
		bs[i] = cipher.decodePwd[v]
	}
}

// 创建一个新的Cipher
func NewCipher(encodePwd *Password) *Cipher {
	decodePwd := &Password{}
	for i, v := range encodePwd {
		encodePwd[i] = v
		decodePwd[v] = byte(i)
	}
	return &Cipher{
		encodePwd: encodePwd,
		decodePwd: decodePwd,
	}
}

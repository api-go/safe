package safe

import (
	crand "crypto/rand"
	"encoding/binary"
	"runtime"
	"sync"
	"time"
	"unsafe"

	"apigo.cc/go/rand"
	"golang.org/x/crypto/chacha20"
)

// MakeSafeToken 生成指定大小的随机字节令牌
func MakeSafeToken(size int) []byte {
	// 使用 rand.Int (apigo.cc/go/rand) 生成偏移量
	fixsize := rand.Int(1, size/10+2)
	key := make([]byte, size+fixsize*2)
	crand.Read(key)
	return key[fixsize : fixsize+size]
}

// EncryptChaCha20 使用 ChaCha20 进行加密
func EncryptChaCha20(raw []byte, key []byte, salt []byte) []byte {
	cipherObj, _ := chacha20.NewUnauthenticatedCipher(key[:chacha20.KeySize], salt[:chacha20.NonceSizeX])
	cipher := make([]byte, len(raw))
	cipherObj.XORKeyStream(cipher, raw)
	return cipher
}

// DecryptChaCha20 使用 ChaCha20 进行解密
func DecryptChaCha20(cipher []byte, key []byte, salt []byte) []byte {
	cipherObj, _ := chacha20.NewUnauthenticatedCipher(key[:chacha20.KeySize], salt[:chacha20.NonceSizeX])
	plaintext := make([]byte, len(cipher))
	cipherObj.XORKeyStream(plaintext, cipher)
	return plaintext
}

// ZeroMemory 使用随机种子覆盖内存，确保敏感数据擦除
func ZeroMemory(buf []byte) {
	seed := uint64(time.Now().UnixNano())
	i := 0
	n := len(buf)
	for i <= n-8 {
		seed ^= seed << 13
		seed ^= seed >> 7
		seed ^= seed << 17
		binary.LittleEndian.PutUint64(buf[i:], seed)
		i += 8
	}
	if i < n {
		seed ^= seed << 13
		seed ^= seed >> 7
		seed ^= seed << 17
		var tmp [8]byte
		binary.LittleEndian.PutUint64(tmp[:], seed)
		copy(buf[i:], tmp[:n-i])
	}
	runtime.KeepAlive(buf)
}

var chachaGlobalKey = make([]byte, chacha20.KeySize)

func init() {
	crand.Read(chachaGlobalKey)
}

var safeBufEncrypt = func(raw []byte) ([]byte, []byte) {
	salt := MakeSafeToken(chacha20.NonceSizeX)
	return EncryptChaCha20(raw, chachaGlobalKey, salt), salt
}

var safeBufDecrypt = func(cipher []byte, salt []byte) []byte {
	return DecryptChaCha20(cipher, chachaGlobalKey, salt)
}

var setObfOnce sync.Once

// SetSafeBufObfuscator 设置自定义的 SafeBuf 加解密混淆器
func SetSafeBufObfuscator(encrypt func([]byte) ([]byte, []byte), decrypt func([]byte, []byte) []byte) {
	setObfOnce.Do(func() {
		safeBufEncrypt = encrypt
		safeBufDecrypt = decrypt
	})
}

// SafeBuf 用于存储受保护的敏感数据
type SafeBuf struct {
	buf  []byte
	salt []byte
}

// SecretPlaintext 表示敏感数据的明文副本
type SecretPlaintext struct {
	Data []byte
}

// String 返回明文的字符串表示
func (sp *SecretPlaintext) String() string {
	if len(sp.Data) == 0 {
		return ""
	}
	return unsafe.String(&sp.Data[0], len(sp.Data))
}

// Close 清除明文数据
func (sp *SecretPlaintext) Close() {
	if sp.Data != nil {
		ZeroMemory(sp.Data)
		sp.Data = nil
	}
}

// NewSafeBuf 创建一个新的 SafeBuf
func NewSafeBuf(raw []byte) *SafeBuf {
	cipher, salt := safeBufEncrypt(raw)
	LockMemory(cipher)
	LockMemory(salt)
	return &SafeBuf{buf: cipher, salt: salt}
}

// NewSafeBufAndErase 创建一个新的 SafeBuf 并自动擦除原始数据
func NewSafeBufAndErase(raw []byte) *SafeBuf {
	defer ZeroMemory(raw)
	return NewSafeBuf(raw)
}

// NewSafeBufFromEncrypted 从已加密数据创建 SafeBuf
func NewSafeBufFromEncrypted(cipher, salt []byte) *SafeBuf {
	LockMemory(cipher)
	LockMemory(salt)
	return &SafeBuf{buf: cipher, salt: salt}
}

// Open 解密 SafeBuf 并返回明文副本
func (sb *SafeBuf) Open() *SecretPlaintext {
	data := safeBufDecrypt(sb.buf, sb.salt)
	LockMemory(data)
	sp := &SecretPlaintext{Data: data}
	runtime.SetFinalizer(sp, func(obj *SecretPlaintext) {
		obj.Close()
	})
	return sp
}

// Close 擦除 SafeBuf 中的加密数据
func (sb *SafeBuf) Close() {
	UnlockMemory(sb.buf)
	ZeroMemory(sb.buf)
	UnlockMemory(sb.salt)
	ZeroMemory(sb.salt)
}

// NewSafeString 创建一个临时的敏感字符串
func NewSafeString(raw []byte) (*SecretPlaintext, string) {
	sp := &SecretPlaintext{Data: raw}
	runtime.SetFinalizer(sp, func(obj *SecretPlaintext) {
		obj.Close()
	})
	return sp, sp.String()
}

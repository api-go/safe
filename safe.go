package safe

import (
	crand "crypto/rand"
	"encoding/base64"
	"encoding/binary"
	"encoding/hex"
	"net/url"
	"runtime"
	"sync"
	"time"
	"unsafe"

	"github.com/api-go/cast"
	"github.com/api-go/rand"
	"golang.org/x/crypto/chacha20"
)

// MakeSafeToken 生成指定大小的随机字节令牌
func MakeSafeToken(size int) []byte {
	// 使用 rand.Int (github.com/api-go/rand) 生成偏移量
	fixsize := rand.Int(1, size/10+2)
	key := make([]byte, size+fixsize*2)
	crand.Read(key)
	return key[fixsize : fixsize+size]
}

// EncryptChaCha20 使用 ChaCha20 进行加密
func EncryptChaCha20(raw []byte, key []byte, salt []byte) []byte {
	if len(key) < chacha20.KeySize || len(salt) < chacha20.NonceSizeX {
		return nil
	}
	cipherObj, err := chacha20.NewUnauthenticatedCipher(key[:chacha20.KeySize], salt[:chacha20.NonceSizeX])
	if err != nil {
		return nil
	}
	cipher := make([]byte, len(raw))
	cipherObj.XORKeyStream(cipher, raw)
	return cipher
}

// DecryptChaCha20 使用 ChaCha20 进行解密
func DecryptChaCha20(cipher []byte, key []byte, salt []byte) []byte {
	if len(key) < chacha20.KeySize || len(salt) < chacha20.NonceSizeX {
		return nil
	}
	cipherObj, err := chacha20.NewUnauthenticatedCipher(key[:chacha20.KeySize], salt[:chacha20.NonceSizeX])
	if err != nil {
		return nil
	}
	plaintext := make([]byte, len(cipher))
	cipherObj.XORKeyStream(plaintext, cipher)
	return plaintext
}

// ZeroMemory 使用随机种子覆盖内存，确保敏感数据擦除
func ZeroMemory(buf []byte) {
	if len(buf) == 0 {
		return
	}
	var seedBytes [8]byte
	_, _ = crand.Read(seedBytes[:])
	seed := binary.LittleEndian.Uint64(seedBytes[:]) ^ uint64(time.Now().UnixNano())

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
	_, _ = crand.Read(chachaGlobalKey)
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
func (secret *SecretPlaintext) String() string {
	if secret == nil || len(secret.Data) == 0 {
		return ""
	}
	return unsafe.String(&secret.Data[0], len(secret.Data))
}

// Close 清除明文数据
func (secret *SecretPlaintext) Close() {
	if secret != nil && secret.Data != nil {
		ZeroMemory(secret.Data)
		secret.Data = nil
	}
}

// newSecretPlaintext 内部辅助函数，统一绑定 Finalizer 作为安全兜底
func newSecretPlaintext(data []byte) *SecretPlaintext {
	secret := &SecretPlaintext{Data: data}
	// GC Finalizer 仅作为兜底策略，主生命周期应由调用者显式管理
	runtime.SetFinalizer(secret, func(obj *SecretPlaintext) {
		obj.Close()
	})
	return secret
}

// NewSafeBuf 创建一个新的 SafeBuf
func NewSafeBuf(raw []byte) *SafeBuf {
	cipher, salt := safeBufEncrypt(raw)
	_ = LockMemory(cipher)
	_ = LockMemory(salt)
	return &SafeBuf{buf: cipher, salt: salt}
}

// NewSafeBufAndErase 创建一个新的 SafeBuf 并自动擦除原始数据
func NewSafeBufAndErase(raw []byte) *SafeBuf {
	defer ZeroMemory(raw)
	return NewSafeBuf(raw)
}

// NewSafeBufFromEncrypted 从已加密数据创建 SafeBuf
func NewSafeBufFromEncrypted(cipher, salt []byte) *SafeBuf {
	_ = LockMemory(cipher)
	_ = LockMemory(salt)
	return &SafeBuf{buf: cipher, salt: salt}
}

// Open 解密 SafeBuf 并返回明文副本
func (safeBuf *SafeBuf) Open() *SecretPlaintext {
	if safeBuf == nil {
		return &SecretPlaintext{}
	}
	data := safeBufDecrypt(safeBuf.buf, safeBuf.salt)
	if data == nil {
		return &SecretPlaintext{}
	}
	_ = LockMemory(data)
	return newSecretPlaintext(data)
}

// Close 擦除 SafeBuf 中的加密数据
func (safeBuf *SafeBuf) Close() {
	if safeBuf == nil {
		return
	}
	_ = UnlockMemory(safeBuf.buf)
	ZeroMemory(safeBuf.buf)
	_ = UnlockMemory(safeBuf.salt)
	ZeroMemory(safeBuf.salt)
}

// NewSafeString 创建一个临时的敏感字符串
func NewSafeString(raw []byte) (*SecretPlaintext, string) {
	secret := newSecretPlaintext(raw)
	return secret, secret.String()
}

// Concat 安全地拼接多个部分并返回一个持久加密的 SafeBuf 对象
func Concat(parts ...any) *SafeBuf {
	return NewSafeBufAndErase(buildBytes(parts...))
}

// Base64 对输入进行拼接与编码，并返回一个持久加密的 SafeBuf 对象
func Base64(parts ...any) *SafeBuf {
	src := buildBytes(parts...)
	defer ZeroMemory(src)

	buf := make([]byte, base64.StdEncoding.EncodedLen(len(src)))
	base64.StdEncoding.Encode(buf, src)

	return NewSafeBufAndErase(buf)
}

// Hex 将数据拼接并转换为 Hex 编码，返回一个持久加密的 SafeBuf 对象
func Hex(parts ...any) *SafeBuf {
	src := buildBytes(parts...)
	defer ZeroMemory(src)

	buf := make([]byte, hex.EncodedLen(len(src)))
	hex.Encode(buf, src)

	return NewSafeBufAndErase(buf)
}

// UrlEncode 对数据拼接并进行 URL 编码，返回一个持久加密的 SafeBuf 对象
func UrlEncode(parts ...any) *SafeBuf {
	src := buildBytes(parts...)
	defer ZeroMemory(src)

	// 使用 unsafe.String 避免对 src 产生额外的 string 分配
	srcStr := unsafe.String(&src[0], len(src))
	encoded := url.QueryEscape(srcStr)
	buf := []byte(encoded)

	return NewSafeBufAndErase(buf)
}

func buildBytes(parts ...any) []byte {
	totalLen := 0
	for _, p := range parts {
		totalLen += partLen(p)
	}

	buf := make([]byte, totalLen)
	pos := 0
	for _, p := range parts {
		pos += copyPart(buf[pos:], p)
	}
	return buf
}

func partLen(v any) int {
	switch t := v.(type) {
	case string:
		return len(t)
	case []byte:
		return len(t)
	case *SecretPlaintext:
		if t == nil {
			return 0
		}
		return len(t.Data)
	case *SafeBuf:
		if t == nil {
			return 0
		}
		p := t.Open()
		defer p.Close()
		return len(p.Data)
	default:
		return len(cast.String(v))
	}
}

func copyPart(dst []byte, v any) int {
	switch t := v.(type) {
	case string:
		return copy(dst, t)
	case []byte:
		return copy(dst, t)
	case *SecretPlaintext:
		if t == nil {
			return 0
		}
		return copy(dst, t.Data)
	case *SafeBuf:
		if t == nil {
			return 0
		}
		p := t.Open()
		defer p.Close()
		return copy(dst, p.Data)
	default:
		return copy(dst, cast.String(v))
	}
}

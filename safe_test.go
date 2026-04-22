package safe

import (
	"bytes"
	"runtime"
	"testing"
	"unsafe"
)

func TestSafeBuf(t *testing.T) {
	data := []byte("secret_password_123")
	sb := NewSafeBuf(data)
	defer sb.Close()

	sp := sb.Open()
	if !bytes.Equal(sp.Data, data) {
		t.Fatal("SafeBuf decryption failed")
	}

	if sp.String() != string(data) {
		t.Fatal("String representation mismatch")
	}

	sp.Close()
	if sp.Data != nil {
		t.Error("Data not cleared after Close")
	}
}

func TestMemoryErasure(t *testing.T) {
	data := []byte("secret_sensitive_content")
	sb := NewSafeBuf(data)
	sp := sb.Open()

	// 捕获内存指针
	ptr := unsafe.Pointer(&sp.Data[0])
	dataLen := len(sp.Data)

	sp.Close()

	// 检查内存是否被覆盖
	// 注意：在释放内存后直接检查虽然是未定义行为，但在这里是为了验证安全基石
	raw := unsafe.Slice((*byte)(ptr), dataLen)
	if bytes.Equal(raw, data) {
		t.Fatal("Security Breach: Memory not erased after Close()")
	}
}

func TestSafeBufCustomObfuscator(t *testing.T) {
	encrypt := func(raw []byte) ([]byte, []byte) {
		return append([]byte(nil), raw...), []byte("salt")
	}
	decrypt := func(cipher []byte, salt []byte) []byte {
		return cipher
	}

	SetSafeBufObfuscator(encrypt, decrypt)

	data := []byte("custom_obfuscator_test")
	sb := NewSafeBuf(data)
	defer sb.Close()

	sp := sb.Open()
	if !bytes.Equal(sp.Data, data) {
		t.Error("Custom obfuscator failed to roundtrip")
	}
	sp.Close()
}

func TestDisableCoreDump(t *testing.T) {
	if err := DisableCoreDump(); err != nil {
		t.Errorf("DisableCoreDump failed: %v", err)
	}
}

func BenchmarkSafeBufOpenClose(b *testing.B) {
	data := []byte("benchmark_data")
	sb := NewSafeBuf(data)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		sp := sb.Open()
		sp.Close()
		runtime.KeepAlive(sp)
	}
	sb.Close()
}

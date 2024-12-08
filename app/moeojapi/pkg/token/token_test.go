package token

import (
	"context"
	"encoding/base64"
	"errors"
	"testing"
	"time"

	"github.com/r27153733/ByteMoeOJ/app/moeojapi/pkg/sign"
	"github.com/r27153733/ByteMoeOJ/lib/uuid"
)

func TestZeroAlloc(t *testing.T) {
	signing := sign.NewSign("test-secret")
	id := uuid.NewUUIDPtr()
	expiration := time.Now().Add(time.Hour)

	//allocProfile := pprof.Lookup("allocs")
	//
	//f3, err := os.Create("alloc3.pprof")
	//if err != nil {
	//	panic(err)
	//}
	//defer f3.Close()

	avg := testing.AllocsPerRun(100, func() {
		dst := GetTokenBuf()
		UUIDWriteToken(signing, id, expiration, dst)
		ReleaseTokenBuf(dst)
	})

	//if err := allocProfile.WriteTo(f3, 0); err != nil {
	//	panic(err)
	//}

	if avg > 0 {
		t.Fatal()
	}
}

func TestUUIDToTokenAndVerifyToken(t *testing.T) {
	// 初始化签名对象
	signing := sign.NewSign("test-secret")

	// 创建一个 UUID
	id := uuid.NewUUIDPtr()
	expiration := time.Now().Add(time.Hour) // 设置过期时间为 1 小时后

	// 测试 UUIDWriteToken 函数
	token := GetTokenBuf()
	UUIDWriteToken(signing, id, expiration, token)

	// 解码 token 查看格式
	decoded, err := base64.RawURLEncoding.DecodeString(token.String())
	if err != nil || len(decoded) < 40 {
		t.Fatalf("Encoded token is invalid: %v", err)
	}

	// 测试 VerifyToken 函数
	parsedUUID, err := VerifyToken(signing, token)
	if err != nil {
		t.Fatalf("VerifyToken failed: %v", err)
	}

	// 验证解析出的 UUID 是否正确
	if !uuid.Equal(parsedUUID, id) {
		t.Errorf("Expected UUID: %s, got: %s", id, parsedUUID)
	}

	// 测试过期 token
	UUIDWriteToken(signing, id, time.Now().Add(-time.Minute), token)

	_, err = VerifyToken(signing, token)
	if err == nil || !errors.Is(err, ErrTokenExpired) {
		t.Errorf("Expected error for expired token, got: %v", err)
	}
}

func TestInvalidToken(t *testing.T) {
	// 初始化签名对象
	signing := sign.NewSign("test-secret")

	// 测试无效 token
	_, err := VerifyTokenBytes(signing, []byte("invalid-token"))
	if err == nil {
		t.Error("Expected error for invalid token, got nil")
	}

	// 测试错误签名的 token
	id := uuid.NewUUIDPtr()
	token := GetTokenBuf()
	UUIDWriteToken(signing, id, time.Now().Add(time.Hour), token)

	// 修改 token 的签名部分
	(*token)[len(*token)-2] = 'A'

	_, err = VerifyToken(signing, token)
	if err == nil || !errors.Is(err, ErrInvalidToken) {
		t.Errorf("Expected error for tampered token, got: %v", err)
	}
}

func TestKey(t *testing.T) {
	ctx := context.Background()
	not := PeekUUIDFromContext(ctx)
	if not != nil {
		t.Error("Expected nil uuid, got ", not)
	}
	id := uuid.NewUUIDPtr()
	ctx = context.WithValue(ctx, uuidKey{}, id)
	value := PeekUUIDFromContext(ctx)
	if *value != *id {
		t.Error("Expected value to be 1, got", value)
	}
}

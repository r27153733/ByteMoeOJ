package password

import (
	"golang.org/x/crypto/bcrypt"
)

// 密码加密的成本，值越高加密强度越高但耗时也越长
const bcryptCost = 10

// Hash 对密码进行哈希加密
func Hash(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcryptCost)
	return string(bytes), err
}

// Compare 比较密码和哈希值是否匹配
func Compare(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

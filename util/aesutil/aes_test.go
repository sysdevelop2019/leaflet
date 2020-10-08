package aesutil

import (
	"fmt"
	"testing"
)

func TestAes(t *testing.T)  {
	orig := "hello world"
	key := "1234567891234567"
	fmt.Println("原文：", orig)
	encryptCode,_ := Encrypt(orig, key)
	fmt.Println("密文：" , encryptCode)
	decryptCode,_ := Decrypt(encryptCode, key)
	fmt.Println("解密结果：", decryptCode)
}

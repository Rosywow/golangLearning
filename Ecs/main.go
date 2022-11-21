package main

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"github.com/obscuren/ecies"
)

func main() {
	var privateKey *ecdsa.PrivateKey

	//标准库生成ecdsa密钥对
	privateKey, _ = ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	//转为ecies密钥对
	var eciesPrivateKey *ecies.PrivateKey
	var eciesPublicKey *ecies.PublicKey
	eciesPrivateKey = ecies.ImportECDSA(privateKey)
	eciesPublicKey = &eciesPrivateKey.PublicKey

	var message string = "this is a message, 这是需要加密的数据"
	fmt.Println("原始数据: \n" + message)

	//加密
	cipherBytes, _ := ecies.Encrypt(rand.Reader, eciesPublicKey, []byte(message), nil, nil)
	//密文编码为16进制字符串输出
	cipherString := hex.EncodeToString(cipherBytes)

	fmt.Println("加密数据: \n" + cipherString)

	//解密
	//decrypeMessageBytes, _ := eciesPrivateKey.Decrypt(rand.Reader, cipherBytes, nil, nil)
	bytes, _ := hex.DecodeString(cipherString)
	decrypeMessageBytes, _ := eciesPrivateKey.Decrypt(rand.Reader, bytes, nil, nil)
	decrypeMessageString := string(decrypeMessageBytes[:])

	fmt.Println("解密数据: \n" + decrypeMessageString)
}

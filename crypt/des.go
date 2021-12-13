package crypt

import (
	"bytes"
	"crypto/cipher"
	"crypto/des"
	"errors"
	"fmt"
)

func DesEncrypt(origData, key []byte) ([]byte, error) {
	block, err := des.NewCipher(key[:8])
	if err != nil {
		return nil, err
	}
	origData = PKCS5Padding(origData, block.BlockSize())
	fmt.Println(block.BlockSize(), len(key[8:]))
	// origData = ZeroPadding(origData, block.BlockSize())
	blockMode := cipher.NewCBCEncrypter(block, key[8:])
	crypted := make([]byte, len(origData))
	// 根据CryptBlocks方法的说明，如下方式初始化crypted也可以
	// crypted := origData
	blockMode.CryptBlocks(crypted, origData)
	return crypted, nil
}

//Des加密
func DesEncrypt1(src, key []byte) ([]byte, error) {
	block, err := des.NewCipher(key)
	if err != nil {
		return nil, err
	}
	bs := block.BlockSize()
	src = ZeroPadding(src, bs)
	// src = PKCS5Padding(src, bs)
	if len(src)%bs != 0 {
		return nil, errors.New("Need a multiple of the blocksize")
	}
	out := make([]byte, len(src))
	dst := out
	for len(src) > 0 {
		block.Encrypt(dst, src[:bs])
		src = src[bs:]
		dst = dst[bs:]
	}
	return out, nil
}
func ZeroPadding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{0}, padding)
	return append(ciphertext, padtext...)
}

// func DesEncrypt1(origData, key []byte) ([]byte, error) {
// 	if len(origData) < 1 || len(key) < 1 {
// 		return nil, errors.New("wrong data or key")
// 	}
// 	block, err := des.NewCipher(key)
// 	if err != nil {
// 		return nil, err
// 	}
// 	bs := block.BlockSize()
// 	// if len(origData)%bs != 0 {
// 	// 	return nil, errors.New("wrong padding")
// 	// }
// 	out := make([]byte, len(origData))
// 	dst := out
// 	for len(origData) > 0 {
// 		block.Encrypt(dst, origData[:bs])
// 		origData = origData[bs:]
// 		dst = dst[bs:]
// 	}
// 	return out, nil
// }

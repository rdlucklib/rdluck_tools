package crypt

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha1"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"errors"
	"fmt"
	"crypto/ecdsa"
)

// rsa 签名
func SignPKCS1v15(origData, privateKey []byte, hash crypto.Hash) (string, error) {
	h := hash.New()
	h.Write(origData)
	digest := h.Sum(nil)

	block, _ := pem.Decode(privateKey)
	if block == nil {
		return "", errors.New("privateKey key error")
	}
	var pri *rsa.PrivateKey
	var err error
	if hash == crypto.MD5 {
		pubInterface, err := x509.ParsePKCS8PrivateKey(block.Bytes)
		if err != nil {
			return "", err
		}
		pri = pubInterface.(*rsa.PrivateKey)
	} else {
		pri, err = x509.ParsePKCS1PrivateKey(block.Bytes)
		if err != nil {
			return "", err
		}
	}
	data, err := rsa.SignPKCS1v15(nil, pri, hash, digest)
	if err != nil {
		fmt.Errorf("rsaSign SignPKCS1v15 error")
		return "", err
	}
	return base64.StdEncoding.EncodeToString(data), nil
}

// rsa 验签
func VerifyPKCS1v15(origData, signedData, publicKey []byte, hash crypto.Hash) error {
	h := hash.New()
	h.Write(origData)
	digest := h.Sum(nil)
	block, _ := pem.Decode(publicKey)
	if block == nil {
		return errors.New("public key error")
	}
	pubInterface, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return err
	}
	pub := pubInterface.(*rsa.PublicKey)

	return rsa.VerifyPKCS1v15(pub, hash, digest, signedData)
}

//rsa公钥加密
func EncryptPKCS1v15(origData, publicKey []byte) (string, error) {
	block, _ := pem.Decode(publicKey)
	if block == nil {
		return "", errors.New("publicKey error")
	}
	//var pub *rsa.PublicKey
	var err error
	pubInterface, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return "", err
	}
	// pub = pubInterface.(*rsa.PublicKey)
	// data, err := rsa.EncryptPKCS1v15(rand.Reader, pub, origData)
	// if err != nil {
	// 	fmt.Errorf("EncryptPKCS1v15 error")
	// 	return "", err
	// }
	pub := pubInterface.(*rsa.PublicKey)
	encrypted := make([]byte, 0, len(origData))
	for i := 0; i < len(origData); i += 117 {
		if i+117 < len(origData) {
			partial, _ := rsa.EncryptPKCS1v15(rand.Reader, pub, origData[i:i+117])

			encrypted = append(encrypted, partial...)
		} else {
			partial, _ := rsa.EncryptPKCS1v15(rand.Reader, pub, origData[i:])

			encrypted = append(encrypted, partial...)
		}
	}
	return base64.StdEncoding.EncodeToString(encrypted), nil
}

func RsaEncrypt(origData, publicKey []byte) (string, error) {
	block, _ := pem.Decode(publicKey)
	if block == nil {
		return "", errors.New("publicKey error")
	}
	var pub *rsa.PublicKey
	var err error
	pubInterface, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return "", err
	}
	pub = pubInterface.(*rsa.PublicKey)
	encrypt, err := rsa.EncryptOAEP(sha1.New(), rand.Reader, pub, origData, []byte(""))
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(encrypt), nil
}

//DecryptRSA decrypt given []byte with RSA algorithm
func DecryptRSA(data, privateKey []byte) (string, error) {
	block, _ := pem.Decode(privateKey)
	if block == nil {
		return "", errors.New("privateKey error")
	}
	privInterface, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return "", err
	}
	decrypted := make([]byte, 0, len(data))
	for i := 0; i < len(data); i += 128 {
		if i+128 < len(data) {
			partial, err1 := rsa.DecryptPKCS1v15(rand.Reader, privInterface, data[i:i+128])
			if err1 != nil {
				return "", err1
			}
			decrypted = append(decrypted, partial...)
		} else {
			partial, err1 := rsa.DecryptPKCS1v15(rand.Reader, privInterface, data[i:])
			if err1 != nil {
				return "", err1
			}
			decrypted = append(decrypted, partial...)
		}
	}
	// partial, err1 := rsa.DecryptPKCS1v15(rand.Reader, privInterface, data)
	// if err1 != nil {
	// 	return "", err1
	// }
	// decrypted = append(decrypted, partial...)
	return string(decrypted), nil
}

func RsaDecryptNew(ciphertext,privateKey []byte) ([]byte, error) {
	block, _ := pem.Decode(privateKey) //将密钥解析成私钥实例
	if block == nil {
		return nil, errors.New("private key error!")
	}
	priv, err := parsePrivateKey(block.Bytes) //解析pem.Decode（）返回的Block指针实例
	fmt.Println(err)
	if err != nil {
		return nil, err
	}
	return rsa.DecryptPKCS1v15(rand.Reader, priv.(*rsa.PrivateKey), ciphertext) //RSA算法解密
}


func parsePrivateKey(der []byte) (crypto.PrivateKey, error) {
	if key, err := x509.ParsePKCS1PrivateKey(der); err == nil {
		return key, nil
	}
	if key, err := x509.ParsePKCS8PrivateKey(der); err == nil {
		switch key := key.(type) {
		case *rsa.PrivateKey, *ecdsa.PrivateKey:
			return key, nil
		default:
			return nil, errors.New("crypto/tls: found unknown private key type in PKCS#8 wrapping")
		}
	}
	if key, err := x509.ParseECPrivateKey(der); err == nil {
		return key, nil
	}


	return nil, errors.New("crypto/tls: failed to parse private key")
}
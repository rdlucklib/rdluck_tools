package crypt

import (
	"crypto"
	"crypto/hmac"
	"encoding/base64"
	"encoding/hex"
)

func Encrypt(origData, key []byte, hash crypto.Hash) string {
	mac := hmac.New(hash.New, key)
	mac.Write(origData)
	return hex.EncodeToString(mac.Sum(nil))
}

func HmacSha256EncryptToBase64(origData, key []byte) string {
	mac := hmac.New(crypto.SHA256.New, key)
	mac.Write(origData)
	return base64.StdEncoding.EncodeToString(mac.Sum(nil))
}
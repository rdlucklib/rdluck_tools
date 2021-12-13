package aes

import (
	"crypto/des"
	"crypto/cipher"
	"encoding/hex"
	"strings"
)

//加密
func encrypt(content string, key string) string {
	contents := []byte(content)
	keys := []byte(key)
	block, err := des.NewCipher(keys)
	if err != nil {
		return "加密失败" + err.Error()
	}
	contents = PKCS5Padding(contents, block.BlockSize())
	blockMode := cipher.NewCBCEncrypter(block, keys)
	crypted := make([]byte, len(contents))
	blockMode.CryptBlocks(crypted, contents)
	return byteToHexString(crypted)
}

func byteToHexString(bytes []byte) string {
	str := ""
	for i := 0; i < len(bytes); i++ {
		sTemp := hex.EncodeToString([]byte{bytes[i]})
		if len(sTemp) < 2 {
			str += string(0)
		}
		str += strings.ToUpper(sTemp)
	}
	return str
}

//解密
func decrypt(content string, key string) string {
	contentBytes, err := hex.DecodeString(content)
	if err != nil {
		return "字符串转换16进制数组失败" + err.Error()
	}
	keys := []byte(key)
	block, err := des.NewCipher(keys)
	if err != nil {
		return "解密失败" + err.Error()
	}
	blockMode := cipher.NewCBCDecrypter(block, keys)
	origData := contentBytes
	blockMode.CryptBlocks(origData, contentBytes)
	origData = ZeroUnPadding(origData)
	return string(origData)
}

func ZeroUnPadding(origData []byte) []byte {
	length := len(origData)
	unpadding := int(origData[length-1])
	return origData[:(length - unpadding)]
}
/*
func main() {
	key := "WQ1s7rzc"
	bbbString := "customerId=HJGS1806081013&customerProdId=PROD1806159021508163743&name=吴仁彪&mobile=1582716758&idCardNo=42092319910204529&timestamp=1528527568482"
	fmt.Println("最终加密后： " + encrypt(bbbString, key))
	cccString := "3B42AC43AA308F7F38359CA4E50734759D1A3C9E6D5F925BBAE80CB390FF6E9AB5C8326A433070ABFC2DCF230A27F5C86CC6F54D672A3247CA58204A58E51D3F2042FEE7B779E3BEBC11382DF6A7660C9CFBA82EF63091B71EFC45248F5171A3BD7D98A69578A18EAA516E53F6167A518C5332A4C9426BBCA105DD0A760579C7A50A2A6C8ECCD869819C8AA378AD51B52DD1496F95697729"
	fmt.Println("解密后： " + decrypt(cccString, key))
}*/

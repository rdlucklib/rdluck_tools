package nethttp

import (
	"bytes"
	"net/http"
	"fmt"
	"io/ioutil"
	"crypto/tls"
	"crypto/x509"
)

const (
	AcceptXml="application/xml"
	Content_Type_Text_Xml="text/xml"
)



func HttpPost(postUrl,postData,AcceptType string)([]byte,error) {
	req, err := http.NewRequest("POST", postUrl, bytes.NewReader([]byte(postData)))
	if err!=nil {
		fmt.Println(err)
		return []byte("创建请求失败!"),err
	}
	req.Header.Set("Accept",AcceptType)
	req.Header.Set("Content-Type",AcceptType+";charset=utf-8")
	client:=http.Client{}
	resp,err:=client.Do(req)
	if err!=nil {
		fmt.Println(err)
		return []byte("请求失败!"),err
	}
	defer resp.Body.Close()
	result,err:=ioutil.ReadAll(resp.Body)
	if err!=nil {
		fmt.Println(err)
		return []byte("读取返回值失败!"),err
	}
	return result,err
}

func HttpPostHasCertificate(url,contentType string,caPath,certPath,keyPath string,xmlContent []byte) ([]byte, error) {
	tlsConfig, err := getTLSConfig(caPath,certPath,keyPath)
	if err != nil {
		fmt.Println(err)
		return []byte(""), err
	}
	tr := &http.Transport{TLSClientConfig: tlsConfig}
	client := &http.Client{Transport: tr}
	resp, err := client.Post(url, contentType, bytes.NewBuffer(xmlContent))
	defer resp.Body.Close()
	b, err := ioutil.ReadAll(resp.Body)
	return b, err
}

var tlsConfig *tls.Config
func getTLSConfig(caPath,certPath,keyPath string) (*tls.Config,error) {
	if tlsConfig!=nil {
		return tlsConfig,nil
	}
	//加载证书文件
	cert,err:=tls.LoadX509KeyPair(certPath,keyPath)
	if err!=nil {
		return nil,err
	}
	//加载根证书文件
	caData,err:=ioutil.ReadFile(caPath)
	if err!=nil {
		return nil,err
	}
	pool:=x509.NewCertPool()
	pool.AppendCertsFromPEM(caData)
	tlsConfig=&tls.Config{
		Certificates:[]tls.Certificate{cert},
		RootCAs:pool,
	}
	return tlsConfig,nil
}
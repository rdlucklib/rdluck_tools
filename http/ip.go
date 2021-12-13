package http

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

//获取公网IP地址
func GetPublickIp() string {
	ipv4:=""
	resp, err := http.Get("http://myexternalip.com/raw")
	if err != nil {
		fmt.Println(err)
		return ipv4
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err!=nil {
		return ipv4
	}
	ipv4=string(body)
	return ipv4
}

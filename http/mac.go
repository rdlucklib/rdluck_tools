package http

import (
	"net"
	"fmt"
	"bytes"
)

func GetMacAddress() (addr string) {
	interfaces,err:=net.Interfaces()
	if err!=nil {
		fmt.Println(err)
		return
	}
	for _,i:=range interfaces{
		//& 与运算
		if i.Flags&net.FlagUp != 0 && bytes.Compare(i.HardwareAddr,nil)!=0 {
			addr=i.HardwareAddr.String()
			break
		}
	}
	return
}
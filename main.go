package main

import (
	"fmt"
	"rdluck_tools/uuid"
)

func main()  {
	fmt.Println("build")
	str:=uuid.NewUUID().Hex32()
	fmt.Println(str)
}

package main

import (
	"fmt"

	"github.com/golang/protobuf/proto"

	"gRPC-go/proto/test"
)

func main() {
	t := &myproto.Test{
		Name:    "panda",
		Tizhong: []int32{120, 125, 110},
		Shengao: 180,
		Motto:   "天行健,君子以自强不息;地势坤,君子以厚德载物",
	}
	fmt.Println(*t)
	// proto编码
	data, err := proto.Marshal(t)
	if err != nil {
		fmt.Println("proto编码失败")
	} else {
		fmt.Println(data)
		newT := &myproto.Test{}
		err = proto.Unmarshal(data, newT)
		if err != nil {
			fmt.Println("proto解码失败")
		} else {
			fmt.Println(*newT)
			fmt.Println(newT.String())
			fmt.Println(newT.GetName())
			fmt.Println(newT.GetTizhong())
			fmt.Println(newT.GetShengao())
			fmt.Println(newT.GetMotto())
		}
	}
}

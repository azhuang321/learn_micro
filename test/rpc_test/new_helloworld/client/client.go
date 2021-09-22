package main

import (
	"fmt"
	"log"
	"myrpc/new_helloworld/client_proxy"
)

func main() {
	//1.建立连接
	client := client_proxy.NewHelloServiceClient("tcp", "localhost:1234")
	var reply string
	err := client.Hello("world", &reply)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(reply)
}

package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
)

func main() {
	client := &http.Client{}
	postData := "{\"id\":0,\"params\":[\"test\"],\"method\":\"HelloService.Hello\"}"
	// 先自定义一个 Request
	req, err := http.NewRequest("POST", "http://localhost:1234/jsonrpc", bytes.NewBuffer([]byte(postData)))
	if err != nil {
		panic("连接失败")
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	fmt.Println("response Body:", string(body))
}

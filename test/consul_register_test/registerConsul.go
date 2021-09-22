package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

func main() {
	register("172.17.0.1", 8000)
}

func register(ip string, port int) {
	url := "http://127.0.0.1:8500/v1/agent/service/register"
	//	payloadStr := `{
	//    "Name": "mshop-web",
	//    "ID": "mshop-web",
	//    "Tags": [
	//        "mxshop",
	//        "bobby",
	//        "imooc",
	//        "web"
	//    ],
	//    "Address": "%s",
	//    "Port": %d,
	//	"Check":{
	//		"HTTP": "http://%s:%d/health",
	//		"Timeout":"5s",
	//		"Interval":"5s",
	//		"DeregisterCriticalServiceAfter":"5s"
	//	}
	//}`

	payloadStr := `{
    "Name": "mshop-web",
    "ID": "mshop-web",
    "Tags": [
        "mxshop",
        "bobby",
        "imooc",
        "web"
    ],
    "Address": "%s",
    "Port": %d,
	"Check":{
		"GRPC":"%s:%d",
		"GRPCUseTLS": false,
		"Timeout":"5s",
		"Interval":"5s",
		"DeregisterCriticalServiceAfter":"5s"
	}
}`

	payloadStr = fmt.Sprintf(payloadStr, ip, port, ip, port)

	payload := strings.NewReader(payloadStr)
	req, _ := http.NewRequest("PUT", url, payload)
	req.Header.Add("Content-Type", "application/json")
	resp, _ := http.DefaultClient.Do(req)
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK {
		fmt.Println("请求错误:", resp.StatusCode)
		fmt.Println(string(body))
		return
	}
	fmt.Println(string(body))
}

func deregister(id string) {
	url := "http://127.0.0.1:8500/v1/agent/service/deregister/" + id
	req, _ := http.NewRequest("PUT", url, strings.NewReader(""))
	resp, _ := http.DefaultClient.Do(req)
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK {
		fmt.Println("请求错误:", resp.StatusCode)
		fmt.Println(string(body))
		return
	}
	fmt.Println(string(body))
}

package utils

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"flag"
	"fmt"
	"net"

	. "user_srv/config"
)

// MD5 md5字符串
func MD5(str string) string {
	salt := Config.MD5Salt
	hash := md5.New()
	hash.Write([]byte(str + salt))
	return hex.EncodeToString(hash.Sum(nil))
}

// GetFreePort 获取当前系统空闲的端口
func GetFreePort() (int, error) {
	addr, err := net.ResolveTCPAddr("tcp", "localhost:0")
	if err != nil {
		return 0, err
	}
	l, err := net.ListenTCP("tcp", addr)
	if err != nil {
		return 0, err
	}
	defer l.Close()
	return l.Addr().(*net.TCPAddr).Port, nil
}

// GetArgs 获取参数
func GetArgs() map[string]interface{} {
	var host string
	var port int
	flag.StringVar(&host, "h", "172.17.0.1", "主机名,默认 172.17.0.1")
	flag.IntVar(&port, "p", 9000, "端口号,默认为 9000")
	flag.Parse()

	if Config.RunMod == "debug" {
		port, _ = GetFreePort()
	}

	return map[string]interface{}{
		"host": host,
		"port": port,
	}
}

func PrettyPrint(v interface{}) {
	b, err := json.Marshal(v)
	if err != nil {
		fmt.Printf("打印错误:%s", err.Error())
		return
	}
	var out bytes.Buffer
	err = json.Indent(&out, b, "", "  ")
	if err != nil {
		fmt.Printf("打印错误:%s", err.Error())
		return
	}
	fmt.Println(out.String())
}

package main

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"os"
	"time"
)

const resolvPath = "/etc/resolv.conf"

func main() {
	// Step 1: 读取原始内容
	original, err := os.ReadFile(resolvPath)
	if err != nil {
		fmt.Println(fmt.Errorf("读取失败:%s", err))
	}

	fmt.Println("===== 原始 /etc/resolv.conf =====")
	fmt.Println(string(original))

	// Step 2: 构造新的内容（添加一个自定义 DNS）
	newContent := string(original) + "\nnameserver 1.1.1.1\n"

	// Step 3: 写入新内容
	err = os.WriteFile(resolvPath, []byte(newContent), 0644)
	if err != nil {
		fmt.Println(fmt.Errorf("写入失败:%s", err))
	}

	// Step 4: 读取并打印修改后的内容
	modified, err := os.ReadFile(resolvPath)
	if err != nil {
		fmt.Println(fmt.Errorf("再次读取失败:%s", err))
	}

	fmt.Println("===== 修改后 /etc/resolv.conf =====")
	fmt.Println(string(modified))
	log.Println("test get k8s info")
	i := 1
	for {
		log.Printf("hello world-->%d", i)
		i++
		time.Sleep(2 * time.Second)
	}
}

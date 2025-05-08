package main

import (
	"fmt"
	"github.com/nats-io/nats.go"
	"log"
	"time"
)

func main() {
	// 连接到 NATS 服务器
	natsURL := "nats://nats.nats.svc.cluster.local:4222"
	nc, err := nats.Connect(natsURL)
	if err != nil {
		log.Fatal(err)
	}
	defer nc.Drain()

	// 设置回调接收结果
	_, err = nc.Subscribe("job.A.result", func(m *nats.Msg) {
		fmt.Printf("Job A received result: %s\n", string(m.Data))
	})
	if err != nil {
		log.Fatal(err)
	}

	i := 0
	for {
		fmt.Printf("Iteration %d\n", i)
		// 发布请求给 Job B，回调 `job.A.result` 用于接收结果
		err = nc.PublishMsg(&nats.Msg{
			Subject: "job.B.run",                                      // 请求发送到 job.B.run 主题
			Reply:   "job.A.result",                                   // 设置回调地址为 job.A.result
			Data:    []byte(fmt.Sprintf("Hello, Job B!, I am %d", i)), // 请求的内容
		})
		if err != nil {
			log.Fatal(err)
		}
		time.Sleep(2 * time.Second)
		i++
	}

	// 保持连接，以便接收回调结果
	select {}
}

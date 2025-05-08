package main

import (
	"fmt"
	"github.com/nats-io/nats.go"
	"log"
)

func main() {
	// 连接到 NATS 服务器
	natsURL := "nats://nats.nats.svc.cluster.local:4222"
	nc, err := nats.Connect(natsURL)
	if err != nil {
		log.Fatal(err)
	}
	defer nc.Drain()

	// 订阅 `job.B.run` 主题
	nc.Subscribe("job.B.run", func(m *nats.Msg) {
		fmt.Printf("Job B received: %s\n", string(m.Data))

		// 模拟处理，并将结果发送回调用者（Job A）
		result := []byte("Job B completed successfully")
		if m.Reply != "" {
			// 回复到 Job A 的回调地址
			nc.Publish(m.Reply, result)
		}
	})

	// 保持运行，等待接收消息
	select {}
}

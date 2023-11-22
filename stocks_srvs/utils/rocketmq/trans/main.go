package main

import (
	"context"
	"fmt"
	"github.com/apache/rocketmq-client-go/v2"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	"github.com/apache/rocketmq-client-go/v2/producer"
	"os"
	"time"
)

type OrderListener struct {
}

func (o *OrderListener) ExecuteLocalTransaction(*primitive.Message) primitive.LocalTransactionState {
	fmt.Printf("1111 execute\n")
	return primitive.CommitMessageState
}

func (o *OrderListener) CheckLocalTransaction(*primitive.MessageExt) primitive.LocalTransactionState {
	return primitive.RollbackMessageState
}

func main() {
	p, _ := rocketmq.NewTransactionProducer(
		&OrderListener{},
		producer.WithNsResolver(primitive.NewPassthroughResolver([]string{"192.168.2.112:9876"})),
		producer.WithRetry(2),
	)
	err := p.Start()
	if err != nil {
		fmt.Printf("start producer error: %s", err.Error())
		os.Exit(1)
	}
	topic := "trans"
	msg := &primitive.Message{
		Topic: topic,
		Body:  []byte("Hello RocketMQ Go Client! "),
	}
	res, err := p.SendMessageInTransaction(context.Background(), msg)
	if err != nil {
		fmt.Printf("send message error: %s\n", err)
	} else {
		fmt.Printf("2222 send message success: result=%s\n", res.String())
	}
	time.Sleep(time.Hour)
	err = p.Shutdown()
	if err != nil {
		fmt.Printf("shutdown producer error: %s", err.Error())
	}
}

package pubsub

import (
	"log"
	"testing"
	"time"
)

func TestBasicFunction(t *testing.T) {
	server := NewPubSub(1)
	chan1 := server.Subscribe("topic1")
	server.Publish("content1", "topic1")
	if _, ok := <-chan1; !ok {
		t.Errorf("Error found on subscribe\n")
	}
}

func TestTwoSubscribor(t *testing.T) {
	// chan的容量是1
	server := NewPubSub(1)
	chan1 := server.Subscribe("topic1")
	chan2 := server.Subscribe("topic2")

	server.Publish("content1", "topic1")
	server.Publish("content2", "topic2")

	val, ok := <-chan1
	if !ok || val != "content1" {
		t.Errorf("Error Found \n")
	}

	val, ok = <-chan2
	if !ok || val != "content2" {
		t.Errorf("Error found \n")
	}
}

func TestAddSub(t *testing.T) {
	server := NewPubSub(1)
	chan1 := server.Subscribe("topic1")
	server.AddSubscription(chan1, "topic2")
	server.Publish("content2", "topic2")

	if val, ok := <-chan1; !ok {
		t.Errorf("Error on chan1: %v", val)
	}
}

func TestRemovweSub(t *testing.T) {
	server := NewPubSub(10)
	chan1 := server.Subscribe("topic1", "topic2")
	// 往topic2里发布一条消息
	server.Publish("content2", "topic2")
	if val, ok := <-chan1; !ok {
		t.Errorf("Error on add subscriber: %v", val)
	}
	server.RemoveSubscriprition(chan1, "topic1")
	server.Publish("content2", "topic1")

	select {
	case val := <-chan1:
		t.Errorf("Should not get %v notify on remove topic\n", val)
		break
	case <-time.After(time.Second):
		break
	}
}

func BenchmarkAddSub(b *testing.B) {
	// 新建一个容量10000的chan
	big := NewPubSub(100000)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		// 订阅一个话题
		big.Subscribe("1234567890")
	}
}

func BenchmarkRemoveSub(b *testing.B) {
	big := NewPubSub(100000)
	var subChans []chan interface{}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		// 订阅一个话题
		chan1 := big.Subscribe("1234567890")
		subChans = append(subChans, chan1)
	}

	b.ResetTimer()
	for _, v := range subChans {
		big.RemoveSubscriprition(v, "1234567890")
	}
}

func BenchmarkBasicFunction(b *testing.B) {
	ser := NewPubSub(1000000)
	chan1 := ser.Subscribe("topic1")

	for i := 0; i < b.N; i++ {
		ser.Publish("content1", "topic1")

		if _, ok := <-chan1; !ok {
			log.Println(" Error found on subscribed.")
		}
	}
}

package pubsub

import (
	"fmt"
	"testing"
	"time"
)

func ExamplePubsub(t *testing.T) {
	ser := NewPubSub(1)
	c1 := ser.Subscribe("topic1")
	c2 := ser.Subscribe("topic2")
	ser.Publish("content1", "topic1")
	ser.Publish("content2", "topic2")
	fmt.Println(<-c1)
	//Got "test1"
	fmt.Println(<-c2)
	//Got "test2"

	// Add subscription "topic2" for c1.
	ser.AddSubscription(c1, "topic2")

	// Publish new content in topic2
	ser.Publish("content3", "topic2")

	fmt.Println(<-c1)
	//Got "test3"

	// Remove subscription "topic2" in c1
	ser.RemoveSubscriprition(c1, "topic2")

	// publish 就是往topic关联的所有channel中 channel <- msg
	// Publish new content in topic2
	ser.Publish("content4", "topic2")

	select {
	case val := <-c1:
		fmt.Printf("Should not get %v notify on remove topic\n", val)
		break
	case <-time.After(time.Second):
		//Will go here, because we remove subscription topic2 in c1.
		break
	}
}

func ExamplePubsub_Subscribe(t *testing.T) {
	ser := NewPubSub(1)
	//Subscribe topic1 in c1
	c1 := ser.Subscribe("topic1")
	t.Logf("c1 is a channel %v\n", c1)
	//Subscribe multiple topic "topic1", "topic2" and "topic3" in c2
	c2 := ser.Subscribe("topic1", "topic2", "topic3")
	t.Logf("c2 is a channel %v\n", c2)
}

func ExamplePubsub_AddSubscription(t *testing.T) {
	ser := NewPubSub(1)
	//Subscribe topic1 in c1
	c1 := ser.Subscribe("topic1")
	//Add more Subscription "topic2" and "topic3" for c1
	ser.AddSubscription(c1, "topic2", "topic3")
}

package pubsub

type chanMapStringList map[chan interface{}][]string
type stringMapChanList map[string][]chan interface{}

type Pubsub struct {
	// capacity for each chan buffer
	capacity int // 管道容量
	// map to store "chan --> topic List" for find subscription
	clientMapTopics chanMapStringList // 管道 --> 话题列表
	// map to store "topic --> chan list" to publish
	topicMapClients stringMapChanList // 话题 --> 管道列表
}

func NewPubSub(initChanCapacity int) *Pubsub {
	initClientMapTopics := make(chanMapStringList)
	initTopicMapClients := make(stringMapChanList)

	server := Pubsub{clientMapTopics: initClientMapTopics, topicMapClients: initTopicMapClients}
	server.capacity = initChanCapacity
	return &server
}

// Sub: Subcribe channels, the channels could be a list of channels name
// The channel name could be any, without define in server
func (p *Pubsub) Subscribe(topics ...string) chan interface{} {
	// init new chan using capacity as channel buffer
	workChan := make(chan interface{}, p.capacity)
	p.updateTopicMapClient(workChan, topics)
	return workChan
}

// 更新topicMapChanList 同时更新 chanMapTopics
func (p *Pubsub) updateTopicMapClient(clientChan chan interface{}, topics []string) {
	// 更新 topic --> chanlist，把每一个topic下的channel list新增一个 channel
	var updateChanList []chan interface{}
	for _, topic := range topics {
		updateChanList, _ = p.topicMapClients[topic]
		updateChanList = append(updateChanList, clientChan)
		p.topicMapClients[topic] = updateChanList
	}
	// 更新 chan ---> topic list
	p.clientMapTopics[clientChan] = topics
}

//AddSubscription: Add a new topic subscribe to specific client channel
func (p *Pubsub) AddSubscription(clientChan chan interface{}, topics ...string) {
	p.updateTopicMapClient(clientChan, topics)
}

// RemoveSubscriprition: Remove sub topic list on specific chan
func (p *Pubsub) RemoveSubscriprition(clientChan chan interface{}, topics ...string) {
	for _, topic := range topics {
		//Remove from topic->chan map
		if chanList, ok := p.topicMapClients[topic]; ok {
			//remove one client chan in chan List
			var updateChanList []chan interface{}
			for _, client := range chanList {
				if client != clientChan {
					updateChanList = append(updateChanList, client)
				}
			}
			p.topicMapClients[topic] = updateChanList
		}

		//Remove from chan->topic map
		if topicList, ok := p.clientMapTopics[clientChan]; ok {
			var updateTopicList []string
			for _, updateTopic := range topicList {
				if updateTopic != topic {
					updateTopicList = append(updateTopicList, topic)
				}
			}
			p.clientMapTopics[clientChan] = updateTopicList
		}
	}
}

// Publish: Publish a content to a list of channels
// The content could be any type
func (p *Pubsub) Publish(content interface{}, topics ...string) {
	for _, topic := range topics {
		// 往topic关联的channel中都塞入一条消息
		if chanList, ok := p.topicMapClients[topic]; ok {
			// Someone has subscribeed this topic
			for _, channel := range chanList {
				channel <- content
			}
		}
	}
}

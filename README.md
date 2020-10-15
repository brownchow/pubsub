# Go 实现发布订阅模型

参考：

https://redis.io/topics/pubsub

https://www.runoob.com/redis/redis-pub-sub.html



思路： map 和 channel

0、订阅（必须先订阅，才能增加订阅，解除订阅）

入参：topic list，出参: chan

新建一个chan，并指定容量，然后更新 chan --> topic List 和 topic --> chan list

1、添加订阅：  

入参: chan, topic list，同时更新 chan --> topic list 和  topic --> chan list

2、解除订阅：

入参: chan, topic list, 同时更新 chan  --> topic list 和 topic --> chan list

3、发布消息：

入参: content, topic list

可以往多个 topic 中加入消息

取出 topic 对应的 chan 列表，往每一个 chan 里塞一条消息

 

Client 等价于 channel 

感觉不需要维护 chan  ---> topicList 的 map，只要维护 topic --> chanLit的列表就可以了




























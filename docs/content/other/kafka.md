# kafka消息队列集群相关

> 本项目的kafka集群通过compose启动，具体可以查看 `docker-compose-with-kafka.yaml`. 
> 抛弃了原有的zookeeper，使用基于raft的kraft。减少第三方的依赖。

## 连接kafka

> `github.com/IBM/sarama`

定义一个全局的kafka client

```go
var GobalKafka sarama.Client
```

初始化kafka client

```go
func InitKafka() {
	con := sarama.NewConfig()
	con.Producer.Return.Successes = true
	kafkaClient, err := sarama.NewClient(config.Conf.Kafka.Address, con)
	if err != nil {
		return
	}
	GobalKafka = kafkaClient
}
```

## 生产消息

> 具体生产消息代码在: `pkg/kfk/produce.go`

### KafkaProducer 生产单条消息

topic:对应需要生成的topic \
msg:对应的生产的消息

```go
func KafkaProducer(topic string, msg []byte) (err error) {
	producer, err := sarama.NewSyncProducerFromClient(GobalKafka)
	if err != nil {
		return
	}
	message := &sarama.ProducerMessage{
		Topic: topic,
		Value: sarama.StringEncoder(msg),
	}
	_, _, err = producer.SendMessage(message)
	if err != nil {
		return
	}
	return
}
```

### KafkaProducers 生产多条

生产发送多条消息，topic在messages中，具体可以查看`ProducerMessage`的字段结构。

```go
// KafkaProducers 发送多条，topic在messages中
func KafkaProducers(messages []*sarama.ProducerMessage) (err error) {
	producer, err := sarama.NewSyncProducerFromClient(GobalKafka)
	if err != nil {
		return
	}
	err = producer.SendMessages(messages)
	if err != nil {
		return
	}
	return
}
```

### 实际例子

`app/index_platform/service/index_platform.go`

在构建索引的时候,将正排索引信息推到kafka中

```go
go func(docStruct *types.Document) {
    err = input_data.DocData2Kfk(docStruct)
    if err != nil {
        logs.LogrusObj.Error(err)
    }
}(docStruct)
```


```go
func DocData2Kfk(doc *types.Document) (err error) {
	doctByte, _ := doc.MarshalJSON()
	err = kfk.KafkaProducer(consts.KafkaCSVLoaderTopic, doctByte)
	if err != nil {
		logs.LogrusObj.Errorf("DocData2Kfk-KafkaCSVLoaderTopic :%+v", err)
		return
	}

	return
}
```

## 消费消息

消费代码在`pkg/kfk/consume/forward_index.go`

1. 设置消费群组

```go
// Consumer Sarama消费者群体的消费者
type ForwardIndexConsumer struct {
	Ready chan bool
}

// Setup is run at the beginning of a new session, before ConsumeClaim
func (consumer *ForwardIndexConsumer) Setup(sarama.ConsumerGroupSession) error {
	close(consumer.Ready)
	return nil
}

// Cleanup is run at the end of a session, once all ConsumeClaim goroutines have exited
func (consumer *ForwardIndexConsumer) Cleanup(sarama.ConsumerGroupSession) error {
	return nil
}
```

2. 具体消费处理

```go
func (consumer *ForwardIndexConsumer) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	ctx := context.Background()
	task := &types.Task{
		Columns:    []string{"doc_id", "title", "body", "url"},
		BiTable:    "data",
		SourceType: consts.DataSourceCSV,
	}
	iDao := dao.NewInputDataDao(ctx)
	for {
		select {
		case message, ok := <-claim.Messages():
			if !ok {
				logs.LogrusObj.Infof("message channel was closed")
				return nil
			}

			if task.SourceType == consts.DataSourceCSV {
				doc := new(types.Document)
				_ = doc.UnmarshalJSON(message.Value)
				inputData := &model.InputData{
					DocId:  doc.DocId,
					Title:  doc.Title,
					Body:   doc.Body,
					Url:    "",
					Score:  0.0,
					Source: task.SourceType,
				}
				err := iDao.CreateInputData(inputData)
				if err != nil {
					logs.LogrusObj.Errorf("iDao.CreateInputData:%+v", err)
				}
			}

			logs.LogrusObj.Infof("Message claimed: value = %s, timestamp = %v, topic = %s", string(message.Value), message.Timestamp, message.Topic)
			session.MarkMessage(message, "")
		case <-session.Context().Done():
			return nil
		}
	}
}
```
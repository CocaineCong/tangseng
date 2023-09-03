package consts

const LayOutTimeFormat = "2006-01-02 15:04:05"

const (
	ForwardBucket   = "forward"
	ForwardCountKey = "forwardCount"
	TermBucket      = "term"
	DictBucket      = "dict"
)

const (
	EngineBufSize         = 10000
	ForwardCountInitValue = "0"
)

const (
	KafkaCrawlTopic     = "kafka-crawl-topic"
	KafkaCSVLoaderTopic = "search-engine-csv-loader-topic"
)

package crawl

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/CocaineCong/tangseng/consts"
	"github.com/CocaineCong/tangseng/pkg/kfk"
)

type SyncCrawl struct {
}

func (s *SyncCrawl) RunSyncCrawl(ctx context.Context) (err error) {
	topic := consts.KafkaCrawlTopic
	msgs, err := kfk.KafkaConsumer(topic)
	if err != nil {
		return
	}

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	for {
		select {
		case msg := <-msgs:
			fmt.Println(msg)

		case <-sigs:
			return
		}
	}
}

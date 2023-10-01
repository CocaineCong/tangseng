package kfk_register

import (
	"context"
	"fmt"

	"github.com/CocaineCong/tangseng/consts"
	"github.com/CocaineCong/tangseng/pkg/kfk/consume"
)

func RunInvertedIndex(ctx context.Context) {
	err := consume.ForwardIndexKafkaConsume(ctx, consts.KafkaCSVLoaderTopic, consts.KafkaCSVLoaderGroupId, consts.KafkaAssignorRoundRobin)
	if err != nil {
		fmt.Println("RunInvertedIndex-ForwardIndexKafkaConsume err :", err)
	}
}

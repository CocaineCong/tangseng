package kfk_register

import (
	"context"
	"fmt"

	"github.com/CocaineCong/tangseng/consts"
	"github.com/CocaineCong/tangseng/pkg/kfk/consume"
)

func RunTireTree(ctx context.Context) {
	err := consume.TrieTreeKafkaConsume(ctx, consts.KafkaTrieTreeTopic, consts.KafkaTrieTreeGroupId, consts.KafkaAssignorRoundRobin)
	if err != nil {
		fmt.Println("RunTireTree-TrieTreeKafkaConsume :", err)
	}
}

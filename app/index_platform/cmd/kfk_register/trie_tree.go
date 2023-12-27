package kfk_register

import (
	"context"
	"github.com/CocaineCong/tangseng/consts"
	"github.com/CocaineCong/tangseng/pkg/kfk/consume"
	log "github.com/CocaineCong/tangseng/pkg/logger"
	"github.com/pkg/errors"
)

func RunTireTree(ctx context.Context) {
	err := consume.TrieTreeKafkaConsume(ctx, consts.KafkaTrieTreeTopic, consts.KafkaTrieTreeGroupId, consts.KafkaAssignorRoundRobin)
	if err != nil {
		log.LogrusObj.Errorf("consume.TrieTreeKafkaConsume failed, original error: %T %v", errors.Cause(err), errors.Cause(err))
		log.LogrusObj.Errorf("stack trace: \n%+v\n", err)
	}
}

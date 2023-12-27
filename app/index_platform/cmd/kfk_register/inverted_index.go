package kfk_register

import (
	"context"
	log "github.com/CocaineCong/tangseng/pkg/logger"
	"github.com/pkg/errors"

	"github.com/CocaineCong/tangseng/consts"
	"github.com/CocaineCong/tangseng/pkg/kfk/consume"
)

func RunInvertedIndex(ctx context.Context) {
	err := consume.ForwardIndexKafkaConsume(ctx, consts.KafkaCSVLoaderTopic, consts.KafkaCSVLoaderGroupId, consts.KafkaAssignorRoundRobin)
	if err != nil {
		log.LogrusObj.Errorf("consume.ForwardIndexKafkaConsume failed, original error: %T %v", errors.Cause(err), errors.Cause(err))
		log.LogrusObj.Errorf("stack trace: \n%+v\n", err)
	}
}

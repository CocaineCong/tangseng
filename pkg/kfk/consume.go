package kfk

import (
	"github.com/IBM/sarama"

	"github.com/CocaineCong/tangseng/consts"
	logs "github.com/CocaineCong/tangseng/pkg/logger"
)

func GetDefaultConsumeConfig(assignor string) *sarama.Config {
	configK := sarama.NewConfig()
	configK.Version = sarama.DefaultVersion

	switch assignor {
	case consts.KafkaAssignorSticky:
		configK.Consumer.Group.Rebalance.GroupStrategies = []sarama.BalanceStrategy{sarama.NewBalanceStrategySticky()}
	case consts.KafkaAssignorRoundRobin:
		configK.Consumer.Group.Rebalance.GroupStrategies = []sarama.BalanceStrategy{sarama.NewBalanceStrategyRoundRobin()}
	case consts.KafkaAssignorRange:
		configK.Consumer.Group.Rebalance.GroupStrategies = []sarama.BalanceStrategy{sarama.NewBalanceStrategyRange()}
	default:
		logs.LogrusObj.Errorf("Unrecognized consumer group partition assignor: %s", assignor)
	}
	configK.Consumer.Offsets.Initial = sarama.OffsetOldest

	return configK
}

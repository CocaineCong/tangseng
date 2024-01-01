// Licensed to the Apache Software Foundation (ASF) under one
// or more contributor license agreements.  See the NOTICE file
// distributed with this work for additional information
// regarding copyright ownership.  The ASF licenses this file
// to you under the Apache License, Version 2.0 (the
// "License"); you may not use this file except in compliance
// with the License.  You may obtain a copy of the License at
//
//   http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing,
// software distributed under the License is distributed on an
// "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
// KIND, either express or implied.  See the License for the
// specific language governing permissions and limitations
// under the License.

//go:build !windows
// +build !windows

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

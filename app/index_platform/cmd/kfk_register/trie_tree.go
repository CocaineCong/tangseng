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

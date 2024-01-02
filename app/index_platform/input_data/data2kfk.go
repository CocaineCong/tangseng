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

package input_data

import (
	"github.com/CocaineCong/tangseng/consts"
	"github.com/CocaineCong/tangseng/pkg/kfk"
	"github.com/CocaineCong/tangseng/types"
	"github.com/pkg/errors"
)

// DocData2Kfk Doc数据处理
func DocData2Kfk(doc *types.Document) (err error) {
	doctByte, _ := doc.MarshalJSON()
	err = kfk.KafkaProducer(consts.KafkaCSVLoaderTopic, doctByte)
	if err != nil {
		return errors.WithMessagef(err, "DocData2Kfk-KafkaCSVLoaderTopic :%v", err)
	}

	return
}

// DocTrie2Kfk Trie树构建
func DocTrie2Kfk(tokens []string) (err error) {
	for _, k := range tokens {
		err = kfk.KafkaProducer(consts.KafkaTrieTreeTopic, []byte(k))
	}

	if err != nil {
		return errors.WithMessagef(err, "DocTrie2Kfk-KafkaTrieTreeTopic :%v", err)
	}

	return
}

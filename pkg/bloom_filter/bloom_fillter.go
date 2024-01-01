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

package bloom_filter

import (
	"hash"
	"hash/fnv"
	"math"
	"sync"
)

// 布隆过滤器，判断是否已经被索引过了

type BloomFilter struct {
	bits        []bool
	numHashFunc int
	hashFunc    hash.Hash64
	mutex       sync.Mutex
}

func NewBloomFilter(numItems int, falsePositiveRate float64) *BloomFilter {
	numBits := int(math.Ceil((float64(numItems) * math.Log(falsePositiveRate)) / math.Log(1.0/math.Pow(2, math.Log(2)))))
	numHashFunc := int(math.Ceil((float64(numBits) / float64(numItems)) * math.Log(2)))
	return &BloomFilter{
		bits:        make([]bool, numBits),
		numHashFunc: numHashFunc,
		hashFunc:    fnv.New64(),
	}
}

func (bf *BloomFilter) Add(item string) {
	bf.mutex.Lock()
	defer bf.mutex.Unlock()

	for i := 0; i < bf.numHashFunc; i++ {
		bf.hashFunc.Reset()
		bf.hashFunc.Write([]byte(item))
		hashValue := bf.hashFunc.Sum64()
		index := hashValue % uint64(len(bf.bits))
		bf.bits[index] = true
	}
}

func (bf *BloomFilter) Contains(item string) bool {
	bf.mutex.Lock()
	defer bf.mutex.Unlock()

	for i := 0; i < bf.numHashFunc; i++ {
		bf.hashFunc.Reset()
		bf.hashFunc.Write([]byte(item))
		hashValue := bf.hashFunc.Sum64()
		index := hashValue % uint64(len(bf.bits))
		if !bf.bits[index] {
			return false
		}
	}
	return true
}

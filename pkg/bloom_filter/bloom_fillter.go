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
	"math"
	"sync"

	"github.com/bits-and-blooms/bitset"
	"github.com/spaolacci/murmur3"
)

const (
	defaultBfNum        = 100000 // 默认的数
	defaultFalsePercent = 0.01   // 默认的错误率
)

// BloomFilter 布隆过滤器，判断是否已经被索引过了
type BloomFilter struct {
	bits              *bitset.BitSet
	numItems          int
	falsePositiveRate float64
	hashFuncs         []hash.Hash64
	mutex             sync.Mutex
}

type Options func(in *BloomFilter)

func WithNumber(number int) Options {
	return func(in *BloomFilter) {
		in.numItems = number
	}
}

func WithFalseRate(rate float64) Options {
	return func(in *BloomFilter) {
		in.falsePositiveRate = rate
	}
}

// NewBloomFilter 新建一个布隆过滤器
func NewBloomFilter(option ...Options) *BloomFilter {
	bf := &BloomFilter{
		numItems:          defaultBfNum,
		falsePositiveRate: defaultFalsePercent,
	}
	for _, opt := range option {
		opt(bf)
	}
	m := getM(bf.numItems, bf.falsePositiveRate)
	k := getK(m, bf.numItems)
	hashFuncs := make([]hash.Hash64, k)
	for i := 0; i < k; i++ {
		hashFuncs[i] = murmur3.New64WithSeed(uint32(i))
	}
	bf.bits = bitset.New(m)
	bf.hashFuncs = hashFuncs
	return bf
}

func getM(numItems int, falsePositiveRate float64) uint {
	return uint(math.Ceil((float64(numItems) * math.Log(falsePositiveRate)) / math.Log(1.0/math.Pow(2, math.Log(2)))))
}

func getK(m uint, numItems int) int {
	return int(math.Ceil((float64(m) / float64(numItems)) * math.Log(2)))
}

func (bf *BloomFilter) Add(item string) {
	bf.mutex.Lock()
	defer bf.mutex.Unlock()

	for _, hashFunc := range bf.hashFuncs {
		hashFunc.Reset()
		_, _ = hashFunc.Write([]byte(item))
		hashValue := hashFunc.Sum64()
		index := hashValue % uint64(bf.bits.Len())
		bf.bits.Set(uint(index))
	}
}

func (bf *BloomFilter) Contains(item string) bool {
	bf.mutex.Lock()
	defer bf.mutex.Unlock()

	for _, hashFunc := range bf.hashFuncs {
		hashFunc.Reset()
		_, _ = hashFunc.Write([]byte(item))
		hashValue := hashFunc.Sum64()
		index := hashValue % uint64(bf.bits.Len())
		if !bf.bits.Test(uint(index)) {
			return false
		}
	}
	return true
}

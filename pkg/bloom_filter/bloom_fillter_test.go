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
	"fmt"
	"testing"
)

func TestBloomFilter(t *testing.T) {
	bf := NewBloomFilter(1000, 0.01)

	bf.Add("呀哈哈哈")
	bf.Add("国王资料")
	bf.Add("国王之泪")

	fmt.Println(bf.Contains("国王"))         // false
	fmt.Println(bf.Contains("哈哈"))         // false
	fmt.Println(bf.Contains("国王之泪"))       // true
	fmt.Println(bf.Contains("durian"))     // false
	fmt.Println(bf.Contains("elderberry")) // false
}

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

package codec

import (
	"strconv"
	"strings"
)

// VBEncodeNumber 将整数编码为VB编码的字符串
func VBEncodeNumber(n uint32) string {
	var bytes []uint32

	for {
		bytes = append(bytes, n%128+128)
		if n < 128 {
			break
		}
		n = n / 128
	}

	var by []string
	for i := len(bytes) - 1; i >= 0; i-- {
		if i < len(bytes)-1 {
			by = append(by, strconv.FormatUint(uint64(bytes[i]), 2)[1:]+" ")
		} else {
			by = append(by, strconv.FormatUint(uint64(bytes[i]), 2))
		}
	}

	return strings.Join(by, "")
}

// VBDecode 将VB编码的字节序列解码为整数数组
func VBDecode(bytestream []byte) []uint64 {
	var numbers []uint64
	n := uint64(0)

	for i := 0; i < len(bytestream); i++ {
		byteStr := strings.Split(string(bytestream[i]), " ")
		l := len(byteStr)

		for j := 0; j < l; j++ {
			var by uint64
			if j < l-1 {
				by, _ = strconv.ParseUint("0b1"+byteStr[j][1:], 2, 8)
			} else {
				by, _ = strconv.ParseUint("0b"+byteStr[j], 2, 8)
			}

			if by < 128 {
				n = 128*n + by
			} else {
				n = 128*(n-1) + by
			}
		}

		numbers = append(numbers, n)
		n = 0
	}

	return numbers
}

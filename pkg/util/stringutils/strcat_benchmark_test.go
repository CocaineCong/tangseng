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

package stringutils

import (
	"bytes"
	"fmt"
	"strconv"
	"strings"
	"testing"
)

// 1. 直接拼接
func BenchmarkString(b *testing.B) {
	elems := make([]string, 100000)
	for i := 0; i < len(elems); i++ {
		elems[i] = strconv.Itoa(i)
	}
	sum := ""
	length := len(elems)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for j := 0; j < length; j++ {
			sum += elems[j]
		}
	}
	b.StopTimer()
}

// 2. fmt.Sprintf("%s",xxxxx)
func BenchmarkFmtSprintfString(b *testing.B) {
	elems := make([]int, 100000)
	for i := 0; i < len(elems); i++ {
		elems[i] = i
	}
	length := len(elems)
	sum := ""
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for j := 0; j < length; j++ {
			sum += fmt.Sprintf("%d", elems[j])
		}
	}
	b.StopTimer()
}

// 3. string.Builder
func BenchmarkBuilderString(b *testing.B) {
	elems := make([]string, 100000)
	for i := 0; i < len(elems); i++ {
		elems[i] = strconv.Itoa(i)
	}
	var builder strings.Builder
	length := len(elems)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for j := 0; j < length; j++ {
			builder.WriteString(elems[j])
		}
	}
	b.StopTimer()
}

// 4. bytes.Builder
func BenchmarkByteBufferString(b *testing.B) {
	elems := make([]string, 100000)
	for i := 0; i < len(elems); i++ {
		elems[i] = strconv.Itoa(i)
	}
	buffer := new(bytes.Buffer)
	length := len(elems)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for j := 0; j < length; j++ {
			buffer.WriteString(elems[j])
		}
	}
	b.StopTimer()
}

// 5. byte 拼接
func BenchmarkByteConcatString(b *testing.B) {
	elems := make([]string, 100000)
	for i := 0; i < len(elems); i++ {
		elems[i] = strconv.Itoa(i)
	}
	length := len(elems)
	buf := make([]byte, 0, len(elems))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for j := 0; j < length; j++ {
			buf = append(buf, elems[j]...)
		}
	}
	fmt.Println(string(buf))
	b.StopTimer()
}

// 6. 官方包 strings.join()
func BenchmarkConcatStringJoins(b *testing.B) {
	elems := make([]string, 100000)
	for i := 0; i < len(elems); i++ {
		elems[i] = strconv.Itoa(i)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		strings.Join(elems, "")
	}
	b.StopTimer()
}

// 2021/14/macbook pro/m1 pro
// BenchmarkString
// BenchmarkString-10               	       1	4554106125 ns/op
// BenchmarkFmtSprintfString
// BenchmarkFmtSprintfString-10     	       1	2223630000 ns/op
// BenchmarkBuilderString
// BenchmarkBuilderString-10        	    2056	    611903 ns/op
// BenchmarkByteBufferString
// BenchmarkByteBufferString-10     	    2511	    522022 ns/op
// BenchmarkByteConcatString
// BenchmarkByteConcatString-10     	    2386	    449875 ns/op
// BenchmarkConcatStringJoins
// BenchmarkConcatStringJoins-10    	    1138	   1046680 ns/op

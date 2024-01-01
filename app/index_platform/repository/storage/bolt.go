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

package storage

import (
	"github.com/pkg/errors"
	bolt "go.etcd.io/bbolt"
)

// Put 通过bolt写入数据
func Put(db *bolt.DB, bucket string, key []byte, value []byte) error {
	return db.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists([]byte(bucket))
		if err != nil {
			return errors.Wrap(err, "failed to create bucket")
		}
		return errors.Wrap(b.Put(key, value), "failed to put data")
	})
}

// Get 通过bolt获取数据
func Get(db *bolt.DB, bucket string, key []byte) (r []byte, err error) {
	err = db.View(func(tx *bolt.Tx) (err error) {
		b := tx.Bucket([]byte(bucket))
		r = b.Get(key)
		if r == nil { // 如果是空的话，直接创建这个key，然后返回这个key的初始值，也就是0
			r = []byte("0")
			return
		}
		return
	})

	return r, errors.Wrap(err, "failed to get data")
}

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

package redis

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"

	"github.com/CocaineCong/tangseng/config"
	logs "github.com/CocaineCong/tangseng/pkg/logger"
)

// RedisClient Redis缓存客户端单例
var RedisClient *redis.Client
var RedisContext = context.Background()

// InitRedis 在中间件中初始化redis链接
func InitRedis() {
	rConfig := config.Conf.Redis
	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", rConfig.RedisHost, rConfig.RedisPort),
		Username: rConfig.RedisUsername,
		Password: rConfig.RedisPassword,
		DB:       rConfig.RedisDbName,
	})
	_, err := client.Ping(RedisContext).Result()
	if err != nil {
		logs.LogrusObj.Errorln(err)
		panic(err)
	}
	RedisClient = client
}

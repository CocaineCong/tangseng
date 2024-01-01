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

package es

import (
	"fmt"
	"log"

	"github.com/CocaineCong/eslogrus"
	elastic "github.com/elastic/go-elasticsearch"
	"github.com/sirupsen/logrus"

	"github.com/CocaineCong/tangseng/config"
)

var EsClient *elastic.Client

// InitEs 初始化es
func InitEs() {
	eConfig := config.Conf.Es
	esConn := fmt.Sprintf("http://%s:%s", eConfig.EsHost, eConfig.EsPort)
	cfg := elastic.Config{
		Addresses: []string{esConn},
	}
	client, err := elastic.NewClient(cfg)
	if err != nil {
		log.Panic(err)
	}
	EsClient = client
}

// EsHookLog 初始化log日志
func EsHookLog() *eslogrus.ElasticHook {
	eConfig := config.Conf.Es
	hook, err := eslogrus.NewElasticHook(EsClient, eConfig.EsHost, logrus.DebugLevel, eConfig.EsIndex)
	if err != nil {
		log.Panic(err)
	}
	return hook
}

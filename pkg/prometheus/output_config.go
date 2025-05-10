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

package prometheus

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/CocaineCong/tangseng/config"
	log "github.com/CocaineCong/tangseng/pkg/logger"
)

// GenerateAllConfigFile generate configuration files
// for all registered services
func GenerateAllConfigFile() {
	service := config.Conf.Services
	if len(service) == 0 {
		return
	}
	for k := range service {
		GenerateConfigFile(k)
	}
}

// GenerateConfigFile generate configuration files
// for the services
func GenerateConfigFile(job string) {
	instance := GetServerAddress(job)

	f, err := os.OpenFile(fmt.Sprintf("./pkg/prometheus/config/files/%s.json", job), os.O_CREATE|os.O_TRUNC|os.O_RDWR, 0777)
	if err != nil {
		log.LogrusObj.Error(fmt.Sprintf("failed open file prometheus/config/files/%s.json", job), err)
		return
	}
	defer f.Close()
	buf, err := json.MarshalIndent(instance.Conf, "", "    ")
	if err != nil {
		log.LogrusObj.Error("failed marshal", err)
		return
	}
	_, err = f.Write(buf)
	if err != nil {
		log.LogrusObj.Error("failed write to file", err)
		return
	}
}

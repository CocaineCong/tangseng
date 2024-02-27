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
	"net/http"
	"strings"

	"github.com/CocaineCong/tangseng/config"
	"github.com/CocaineCong/tangseng/consts"
	log "github.com/CocaineCong/tangseng/pkg/logger"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// GatewayHandler for expose gateway metrics
func GatewayHandler() gin.HandlerFunc {
	EtcdRegister(config.Conf.Server.Metrics, consts.GatewayJobForPrometheus)
	handler := promhttp.Handler()
	return func(c *gin.Context) {
		handler.ServeHTTP(c.Writer, c.Request)
	}
}

// RpcHandler is for launch a http server
// to expose metrics
func RpcHandler(addr string) {
	port := strings.Split(addr, ":")[1]
	http.Handle("/metrics", promhttp.Handler())
	log.LogrusObj.Panic(http.ListenAndServe(":"+port, nil))
}

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

package shutdown

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	log "github.com/CocaineCong/tangseng/pkg/logger"
)

func GracefullyShutdown(server *http.Server) {
	// 创建系统信号接收器接收关闭信号
	done := make(chan os.Signal, 1)
	/**
	os.Interrupt           -> ctrl+c 的信号
	syscall.SIGINT|SIGTERM -> kill 进程时传递给进程的信号
	*/
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	select {
	case sig := <-done:
		{
			log.LogrusObj.Infoln("stopping service, because received signal:", sig)
			if err := server.Shutdown(context.Background()); err != nil {
				log.LogrusObj.Infof("closing http service gracefully failed: :%v", err)
			}
			log.LogrusObj.Infoln("service has stopped")
			os.Exit(0)
		}
	default:
		os.Exit(0)
	}
}

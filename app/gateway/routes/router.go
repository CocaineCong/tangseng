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

package routes

import (
	"github.com/CocaineCong/tangseng/pkg/prometheus"
	"github.com/gin-gonic/gin"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"

	"github.com/CocaineCong/tangseng/app/gateway/http"
	"github.com/CocaineCong/tangseng/app/gateway/middleware"
	"github.com/CocaineCong/tangseng/consts"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
)

func NewRouter() *gin.Engine {
	r := gin.Default()
	r.Use(middleware.Cors(), middleware.ErrorMiddleware(), otelgin.Middleware(consts.ServiceName))
	store := cookie.NewStore([]byte("something-very-secret"))
	r.Use(sessions.Sessions("mysession", store))
	r.GET("/metrics", prometheus.GatewayHandler())
	v1 := r.Group("/api/v1")
	{
		v1.GET("ping", func(context *gin.Context) {
			context.JSON(200, "success")
		})
		// 用户服务
		v1.POST("/user/register", http.UserRegister)
		v1.POST("/user/login", http.UserLogin)
		// 索引平台
		IndexPlatformRegisterHandlers(v1)
		// 搜索平台
		SearchRegisterHandlers(v1)
		// 需要登录保护
		authed := v1.Group("/")
		authed.Use(middleware.AuthMiddleware())
		{
			// 收藏夹模块
			FavoriteRegisterHandlers(authed)
		}
	}
	return r
}

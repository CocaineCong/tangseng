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
	"github.com/gin-gonic/gin"

	"github.com/CocaineCong/tangseng/app/gateway/http"
)

func FavoriteRegisterHandlers(rg *gin.RouterGroup) {
	favoriteGroup := rg.Group("/favorite")
	{
		favoriteGroup.POST("/create", http.CreateFavorite)
		favoriteGroup.GET("/list", http.ListFavorite)
		favoriteGroup.POST("/update", http.UpdateFavorite)
		favoriteGroup.POST("/delete", http.DeleteFavorite)
	}

	favoriteDetailGroup := rg.Group("/favorite_detail")
	{
		favoriteDetailGroup.POST("/create", http.CreateFavoriteDetail)
		favoriteDetailGroup.GET("/list", http.ListFavoriteDetail)
		favoriteDetailGroup.POST("/delete", http.DeleteFavoriteDetail)
	}

}

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

package timeutils

import (
	"fmt"
	"strings"

	"github.com/golang-module/carbon"
)

// GetTodayDate return 2023-10-01
func GetTodayDate() string {
	return carbon.Now().ToDateString()
}

// GetNowTime return 2023-10-01
func GetNowTime() string {
	timeStr := strings.Split(carbon.Now().String(), " ")
	return fmt.Sprintf("%s-%s", timeStr[0], timeStr[1])
}

// GetMonthDate return 2023-10
func GetMonthDate() string {
	year := carbon.Now().Year()
	month := carbon.Now().Month()
	return fmt.Sprintf("%d-%d", year, month)
}

// GetSeasonDate return 2023-Autumn
func GetSeasonDate() string {
	year := carbon.Now().Year()
	season := carbon.Now().Season()
	return fmt.Sprintf("%d-%s", year, season)
}

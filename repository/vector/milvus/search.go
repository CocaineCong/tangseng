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

package milvus

import (
	"github.com/milvus-io/milvus-sdk-go/v2/entity"
)

type MilvusRequest struct {
	TopK           int
	CollectionName string
	VectorField    string
	OutputFields   []string
	Vectors        []entity.Vector
	Partitions     []string
	Expr           string
	MetricType     entity.MetricType
	SearchParams   *MilvusSearchParams
}

func (r *MilvusRequest) SetMetricType(metricType string) {
	switch metricType {
	case "IP":
		r.MetricType = entity.IP
	case "L2":
		r.MetricType = entity.L2
	}
}

func (r *MilvusRequest) AppendVectors(vectors []float32) {
	r.Vectors = append(r.Vectors, entity.FloatVector(vectors))
}

func (r *MilvusRequest) SetSearchParams(params map[string]interface{}) {
	r.SearchParams = &MilvusSearchParams{params: params}
}

type MilvusSearchParams struct {
	params map[string]interface{}
}

func (msp *MilvusSearchParams) AddRadius(radius float64) {

}

func (msp *MilvusSearchParams) AddRangeFilter(rangeFilter float64) {

}

func (msp *MilvusSearchParams) Params() map[string]interface{} {
	return msp.params
}

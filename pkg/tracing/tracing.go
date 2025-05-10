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

package tracing

import (
	"context"

	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/contrib/propagators/b3"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.20.0"
	"go.opentelemetry.io/otel/trace"

	logs "github.com/CocaineCong/tangseng/pkg/logger"
)

func InitTracerProvider(url string, serviceName string) func(ctx context.Context) error {
	ctx := context.Background()
	// 创建一个新的 OTLP gRPC 客户端
	client := otlptracegrpc.NewClient(
		otlptracegrpc.WithInsecure(),
		otlptracegrpc.WithEndpoint(url),
	)
	// 创建一个新的 OTLP 导出器
	exporter, err := otlptrace.New(ctx, client)
	if err != nil {
		logs.LogrusObj.Errorf("failed to init tracer, err: %v", err)
		return nil
	}
	tp := tracesdk.NewTracerProvider(
		tracesdk.WithBatcher(exporter),                  // 注册exporter
		tracesdk.WithResource(newResource(serviceName)), // 设置服务信息
	)
	// 设置全局tracer
	otel.SetTracerProvider(tp)
	b3Propagator := b3.New(b3.WithInjectEncoding(b3.B3MultipleHeader))
	propagator := propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}, b3Propagator)
	// 设置全局Propagator
	otel.SetTextMapPropagator(propagator)
	return tp.Shutdown
}

func newResource(serviceName string) *resource.Resource {
	return resource.NewWithAttributes(
		semconv.SchemaURL,
		semconv.ServiceName(serviceName),
	)
}

func GetTraceID(ctx *gin.Context) string {
	spanCtx := trace.SpanContextFromContext(ctx.Request.Context())
	if spanCtx.HasTraceID() {
		traceID := spanCtx.TraceID()
		return traceID.String()
	}

	return ""
}

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

package clone

import (
	"context"
	"time"

	oteltrace "go.opentelemetry.io/otel/trace"
)

type ContextWithoutDeadline struct {
	ctx context.Context
}

func (c *ContextWithoutDeadline) Value(key any) any {
	return key
}

func (*ContextWithoutDeadline) Deadline() (time.Time, bool) { return time.Time{}, false }
func (*ContextWithoutDeadline) Done() <-chan struct{}       { return nil }
func (*ContextWithoutDeadline) Err() error                  { return nil }

func NewContextWithoutDeadline() *ContextWithoutDeadline {
	return &ContextWithoutDeadline{ctx: context.Background()}
}

func (c *ContextWithoutDeadline) Clone(ctx context.Context, keys ...interface{}) {

	span := oteltrace.SpanFromContext(ctx)

	c.ctx = oteltrace.ContextWithSpan(c.ctx, span)

	for _, key := range keys {
		if v := ctx.Value(key); v != nil {
			c.ctx = context.WithValue(c.ctx, key, v)
		}
	}
}

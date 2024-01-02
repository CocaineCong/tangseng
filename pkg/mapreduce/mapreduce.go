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

package mapreduce

import (
	"context"
	"sync"
	"sync/atomic"

	log "github.com/CocaineCong/tangseng/pkg/logger"
	"github.com/pkg/errors"
)

// MapReduce By chan and goroutine to replace using rpc to make master or worker
// https://github.com/kevwan/mapreduce

const (
	defaultWorkers = 5
	minWorkers     = 1
)

var (
	ErrCancelWithNil           = errors.New("mapreduce cancelled with nil")
	ErrReduceNoOutput          = errors.New("reduce not writing value")
	ErrWriteMoreThanOneProduce = errors.New("more than one element written in reducer")
)

type (
	// MapFunc Map的具体方法函数
	MapFunc[T, U any] func(item T, writer Writer[U])
	// ReduceFunc Reduce的具体函数方法
	ReduceFunc[T, U any] func(item T, writer Writer[U])
	// MapperFunc 执行map的函数
	MapperFunc[T, U any] func(item T, writer Writer[U], cancel func(err error))
	// ReducerFunc 执行reduce的函数
	ReducerFunc[U, V any] func(pipe <-chan U, writer Writer[V], cancel func(err error))
	// GenerateFunc 执行生成的函数
	GenerateFunc[T any] func(source chan<- T)
	// Option 可选择的选择
	Option func(opts *mapReduceOptions)
	// mapperContext map所需要的参数
	mapperContext[T, U any] struct {
		ctx       context.Context
		mapper    MapFunc[T, U]
		source    <-chan T
		panicChan *onceChan
		collector chan<- U
		doneChan  <-chan struct{}
		workers   int
	}
	Writer[T any] interface {
		Write(v T)
	}
	mapReduceOptions struct {
		ctx     context.Context
		workers int
	}
)

func MapReduce[T, U, V any](gen GenerateFunc[T], mapper MapperFunc[T, U], reducer ReducerFunc[U, V], opts ...Option) (v V, err error) {
	panicChan := &onceChan{channel: make(chan any)}
	source := buildSource(gen, panicChan)
	v, err = mapReduceWithPanicChan(source, panicChan, mapper, reducer, opts...)
	return v, errors.WithMessage(err, "mapReduceWithPanicChan error")
}

func newOptions() *mapReduceOptions {
	return &mapReduceOptions{
		ctx:     context.Background(),
		workers: defaultWorkers,
	}
}

func buildOptions(opts ...Option) *mapReduceOptions {
	options := newOptions()
	for _, opt := range opts {
		opt(options)
	}

	return options
}

func mapReduceWithPanicChan[T, U, V any](source <-chan T, panicChan *onceChan, mapper MapperFunc[T, U], reducer ReducerFunc[U, V], opts ...Option) (val V, err error) {
	options := buildOptions(opts...)
	// output is used to write the final result
	output := make(chan V)
	defer func() {
		// reducer can only write once, if more, panic
		for range output {
			log.LogrusObj.Errorln(ErrWriteMoreThanOneProduce)
			panic(ErrWriteMoreThanOneProduce)
		}
	}()

	// collector is used to collect data from mapper, and consume in reducer
	collector := make(chan U, options.workers)
	// if done is closed, all mappers and reducer should stop processing
	done := make(chan struct{})
	writer := newGuardedWriter(options.ctx, output, done)
	var closeOnce sync.Once
	// use atomic.Value to avoid data race
	var retErr atomic.Value
	finish := func() {
		closeOnce.Do(func() {
			close(done)
			close(output)
		})
	}
	cancel := once(func(err error) {
		if err != nil {
			retErr.Store(err)
		} else {
			retErr.Store(ErrCancelWithNil)
		}

		drain(source)
		finish()
	})

	go func() {
		defer func() {
			drain(collector)
			if r := recover(); r != nil {
				panicChan.write(r)
			}
			finish()
		}()

		reducer(collector, writer, cancel)
	}()

	go executeMappers(mapperContext[T, U]{
		ctx: options.ctx,
		mapper: func(item T, w Writer[U]) {
			mapper(item, w, cancel)
		},
		source:    source,
		panicChan: panicChan,
		collector: collector,
		doneChan:  done,
		workers:   options.workers,
	})

	select {
	case <-options.ctx.Done():
		cancel(context.DeadlineExceeded)
		err = context.DeadlineExceeded
	case v := <-panicChan.channel:
		// drain output here, otherwise for loop panic in defer
		drain(output)
		panic(v)
	case v, ok := <-output:
		if e := retErr.Load(); e != nil {
			err = e.(error)
		} else if ok {
			val = v
		} else {
			err = ErrReduceNoOutput
		}
	}

	return val, errors.Wrap(err, "mapReduceWithPanicChan error")
}

func once(fn func(error)) func(error) {
	on := new(sync.Once)
	return func(err error) {
		on.Do(func() {
			fn(err)
		})
	}
}

func buildSource[T any](generate GenerateFunc[T], panicChan *onceChan) chan T {
	source := make(chan T)
	go func() {
		defer func() {
			if r := recover(); r != nil {
				panicChan.write(r)
			}
			close(source)
		}()

		generate(source)
	}()

	return source
}

func newGuardedWriter[T any](ctx context.Context, channel chan<- T, done <-chan struct{}) guardedWriter[T] {
	return guardedWriter[T]{
		ctx:     ctx,
		channel: channel,
		done:    done,
	}
}

type guardedWriter[T any] struct {
	ctx     context.Context
	channel chan<- T
	done    <-chan struct{}
}

func (gw guardedWriter[T]) Write(v T) {
	select {
	case <-gw.ctx.Done():
		return
	case <-gw.done:
		return
	default:
		gw.channel <- v
	}
}

type onceChan struct {
	channel chan any
	wrote   int32
}

func (oc *onceChan) write(val any) {
	if atomic.CompareAndSwapInt32(&oc.wrote, 0, 1) {
		oc.channel <- val
	}
}

// drain drains the channel.
func drain[T any](channel <-chan T) {
	// drain the channel
	for range channel {
	}
}

func executeMappers[T, U any](mCtx mapperContext[T, U]) {
	var wg sync.WaitGroup
	defer func() {
		wg.Wait()
		close(mCtx.collector)
		drain(mCtx.source)
	}()

	var failed int32
	pool := make(chan struct{}, mCtx.workers)
	writer := newGuardedWriter(mCtx.ctx, mCtx.collector, mCtx.doneChan)
	for atomic.LoadInt32(&failed) == 0 {
		select {
		case <-mCtx.ctx.Done():
			return
		case <-mCtx.doneChan:
			return
		case pool <- struct{}{}:
			item, ok := <-mCtx.source
			if !ok {
				<-pool
				return
			}

			wg.Add(1)
			go func() {
				defer func() {
					if r := recover(); r != nil {
						atomic.AddInt32(&failed, 1)
						mCtx.panicChan.write(r)
					}
					wg.Done()
					<-pool
				}()

				mCtx.mapper(item, writer)
			}()
		}
	}
}

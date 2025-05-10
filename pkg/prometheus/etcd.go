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
	"context"
	"fmt"
	"time"

	etcd "go.etcd.io/etcd/client/v3"

	"github.com/CocaineCong/tangseng/config"
	"github.com/CocaineCong/tangseng/consts"
	log "github.com/CocaineCong/tangseng/pkg/logger"
)

// Instance is for marshal conf
type Instance struct {
	Conf []*Conf
}

// Conf is the basic unit of the prometheus detection unit
type Conf struct {
	Targets []string          `json:"targets"`
	Labels  map[string]string `json:"labels"`
}

// EtcdRegister need server address and name
// for register to etcd and keep alive
func EtcdRegister(targets string, job string) {
	client := newClient()
	leaseResp, err := client.Grant(context.Background(), 15)
	if err != nil {
		log.LogrusObj.Error(err)
	}
	key := fmt.Sprintf("%s/%s/%d", consts.PrometheusJobKey, job, leaseResp.ID)
	if _, err = client.Put(context.Background(), key, targets, etcd.WithLease(leaseResp.ID)); err != nil {
		log.LogrusObj.Error(err)
		return
	}

	go keepALive(client, leaseResp.ID)
	go GenerateConfigFile(job)
}

// keepAlive for registered instance
func keepALive(c *etcd.Client, leaseId etcd.LeaseID) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	keepLiveCh, _ := c.KeepAlive(ctx, leaseId)

	for {
		select {
		case <-keepLiveCh:
			break
		case <-time.After(time.Duration(15) * time.Second):
			log.LogrusObj.Error("A server lose heart")
			return
		}
	}
}

// GetServerAddress get all addresses for this job
func GetServerAddress(job string) *Instance {
	client := newClient()
	resp, err := client.Get(context.Background(), fmt.Sprintf("%s/%s", consts.PrometheusJobKey, job), etcd.WithPrefix())
	if err != nil {
		log.LogrusObj.Error("failed get server")
		return nil
	}

	if resp.Count == 0 {
		return nil
	}
	addresses := make([]string, 0)
	for _, v := range resp.Kvs {
		addr := string(v.Value)
		if addr != "" {
			addresses = append(addresses, addr)
		}
	}
	conf := make([]*Conf, 1)
	conf[0] = &Conf{
		Targets: addresses,
		Labels: map[string]string{
			"job": job,
		}}
	return &Instance{
		Conf: conf,
	}
}

// GetAllServerAddress Get addresses for all the job
func GetAllServerAddress() []*Instance {
	service := config.Conf.Services
	if len(service) == 0 {
		return nil
	}
	instances := make([]*Instance, len(service))
	for k := range service {
		instances = append(instances, GetServerAddress(k))
	}
	return instances
}

// newClient return an etcd.Client
func newClient() *etcd.Client {
	client, err := etcd.New(etcd.Config{
		Endpoints:   []string{config.Conf.Etcd.Address},
		DialTimeout: 2 * time.Second,
	})
	if err != nil {
		log.LogrusObj.Error(err)
	}
	return client
}

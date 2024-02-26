package prometheus

import (
	"context"
	"fmt"
	"time"

	"github.com/CocaineCong/tangseng/config"
	"github.com/CocaineCong/tangseng/consts"
	log "github.com/CocaineCong/tangseng/pkg/logger"
	etcd "go.etcd.io/etcd/client/v3"
)

type Instance struct {
	Conf []*Conf
}

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

func keepALive(c *etcd.Client, leaseId etcd.LeaseID) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	keepLiveCh, _ := c.KeepAlive(ctx, leaseId)

	for {
		select {
		case <-keepLiveCh:
			break
		case <-time.After(time.Duration(15) * time.Second):
			log.LogrusObj.Error(fmt.Sprintf("A server lose heart"))
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
	for k, _ := range service {
		instances = append(instances, GetServerAddress(k))
	}
	return instances
}

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

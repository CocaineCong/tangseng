package prometheus

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/CocaineCong/tangseng/config"
	log "github.com/CocaineCong/tangseng/pkg/logger"
)

// GenerateAllConfigFile generate configuration files
// for all registered services
func GenerateAllConfigFile() {
	service := config.Conf.Services
	if len(service) == 0 {
		return
	}
	for k, _ := range service {
		GenerateConfigFile(k)
	}
}

// GenerateConfigFile generate configuration files
// for the services
func GenerateConfigFile(job string) {
	instance := GetServerAddress(job)

	f, err := os.OpenFile(fmt.Sprintf("./pkg/prometheus/config/files/%s.json", job), os.O_CREATE|os.O_TRUNC|os.O_RDWR, 0777)
	if err != nil {
		log.LogrusObj.Error(fmt.Sprintf("failed open file prometheus/config/files/%s.json", job), err)
		return
	}
	defer f.Close()
	buf, err := json.MarshalIndent(instance.Conf, "", "    ")
	if err != nil {
		log.LogrusObj.Error("failed marshal", err)
		return
	}
	_, err = f.Write(buf)
	if err != nil {
		log.LogrusObj.Error("failed write to file", err)
		return
	}
}

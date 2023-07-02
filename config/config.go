package config

import (
	"os"

	"github.com/spf13/viper"
)

var Conf *Config

type Config struct {
	Server   *Server             `yaml:"server"`
	MySQL    *MySQL              `yaml:"mysql"`
	Redis    *Redis              `yaml:"redis"`
	Etcd     *Etcd               `yaml:"etcd"`
	Es       *Es                 `yaml:"es"`
	Services map[string]*Service `yaml:"services"`
	Domain   map[string]*Domain  `yaml:"domain"`
	SeConfig *SeConfig           `yaml:"seConfig"`
}

type SeConfig struct {
	StoragePath      string   `yaml:"storagePath"`
	SourceFiles      []string `yaml:"sourceFiles"`
	MergeChannelSize int64    `yaml:"mergeChannelSize"`
	Version          string   `yaml:"version"`
	SourceWuKoFile   string   `yaml:"sourceWuKoFile"`
}

type Server struct {
	Port      string `yaml:"port"`
	Version   string `yaml:"version"`
	JwtSecret string `yaml:"jwtSecret"`
}

type MySQL struct {
	DriverName string `yaml:"driverName"`
	Host       string `yaml:"host"`
	Port       string `yaml:"port"`
	Database   string `yaml:"database"`
	UserName   string `yaml:"username"`
	Password   string `yaml:"password"`
	Charset    string `yaml:"charset"`
}

type Es struct {
	EsHost  string `yaml:"esHost"`
	EsPort  string `yaml:"esPort"`
	EsIndex string `yaml:"esIndex"`
}

type Redis struct {
	UserName string `yaml:"userName"`
	Address  string `yaml:"address"`
	Password string `yaml:"password"`
}

type Etcd struct {
	Address string `yaml:"address"`
}

type Service struct {
	Name        string   `yaml:"name"`
	LoadBalance bool     `yaml:"loadBalance"`
	Addr        []string `yaml:"addr"`
}

type Domain struct {
	Name string `yaml:"name"`
}

func InitConfig() {
	workDir, _ := os.Getwd()
	viper.SetConfigName("config")
	viper.SetConfigType("yml")
	viper.AddConfigPath(workDir + "/config")
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}
	err = viper.Unmarshal(&Conf)
	if err != nil {
		panic(err)
	}
}

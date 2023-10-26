package config

import (
	"os"

	"github.com/spf13/viper"
)

var Conf *Config

type Config struct {
	Server    *Server             `yaml:"service"`
	MySQL     *MySQL              `yaml:"mysql"`
	Redis     *Redis              `yaml:"redis"`
	Etcd      *Etcd               `yaml:"etcd"`
	Es        *Es                 `yaml:"es"`
	Services  map[string]*Service `yaml:"services"`
	Domain    map[string]*Domain  `yaml:"domain"`
	SeConfig  *SeConfig           `yaml:"SeConfig"`
	Kafka     *Kafka              `yaml:"kafka"`
	StarRocks *StarRocks          `yaml:"starrock"`
	Vector    *VectorConfig       `yaml:"vector"`
	Milvus    *MilvusConfig       `yaml:"milvus"`
}

type VectorConfig struct {
	ServerAddress string `yaml:"server_address"`
	Timeout       int64  `yaml:"timeout"`
}

type MilvusConfig struct {
	Host                   string `yaml:"host"`
	Port                   string `yaml:"port"`
	VectorDimension        int    `yaml:"vector_dimension"`
	DefaultMilvusTableName string `yaml:"default_milvus_table_name"`
	MetricType             string `yaml:"metric_type"`
	Timeout                int    `yaml:"timeout`
}

type StarRocks struct {
	UserName string `yaml:"user_name"`
	Password string `yaml:"password"`
	Database string `yaml:"database"`
	LoadUrl  string `yaml:"load_url"`
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	Charset  string `yaml:"charset"`
}

type SeConfig struct {
	StoragePath      string   `yaml:"StoragePath"`
	SourceFiles      []string `yaml:"SourceFiles"`
	MergeChannelSize int64    `yaml:"MergeChannelSize"`
	Version          string   `yaml:"Version"`
	SourceWuKoFile   string   `yaml:"SourceWuKoFile"`
	MetaPath         string   `yaml:"MetaPath"`
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
	RedisHost     string `yaml:"redisHost"`
	RedisPort     string `yaml:"redisPort"`
	RedisUsername string `yaml:"redisUsername"`
	RedisPassword string `yaml:"redisPassword"`
	RedisDbName   int    `yaml:"redisDbName"`
}

type Etcd struct {
	Address string `yaml:"address"`
}

type Service struct {
	Name        string   `yaml:"name"`
	LoadBalance bool     `yaml:"loadBalance"`
	Addr        []string `yaml:"addr"`
}

type Kafka struct {
	Address []string `yaml:"address"`
}

type Domain struct {
	Name string `yaml:"name"`
}

func InitConfig() {
	workDir, _ := os.Getwd()
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
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

package util

import (
	"log"
	"strconv"
	"strings"

	"github.com/spf13/viper"
)

type DockerComposeYml struct {
	Services ServicesConfig `mapstructure:"services"`
}

type ServicesConfig struct {
	Influxdb InfluxdbConfig `mapstructure:"influxdb"`
	Mongo    MongoConfig    `mapstructure:"mongo"`
}

type MongoConfig struct {
	Env        MongoEnv `mapstructure:"environment"`
	Ports      []string `mapstructure:"ports"`
	PortSrc    int
	PortTarget int
}

type MongoEnv struct {
	Username string `mapstructure:"MONGO_INITDB_ROOT_USERNAME"`
	Password string `mapstructure:"MONGO_INITDB_ROOT_PASSWORD"`
}

type InfluxdbConfig struct {
	Env        InfluxdbEnv `mapstructure:"environment"`
	Ports      []string    `mapstructure:"ports"`
	PortSrc    int
	PortTarget int
}

type InfluxdbEnv struct {
	Mode     string `mapstructure:"DOCKER_INFLUXDB_INIT_MODE"`
	Username string `mapstructure:"DOCKER_INFLUXDB_INIT_USERNAME"`
	Password string `mapstructure:"DOCKER_INFLUXDB_INIT_PASSWORD"`
	Org      string `mapstructure:"DOCKER_INFLUXDB_INIT_ORG"`
	Bucket   string `mapstructure:"DOCKER_INFLUXDB_INIT_BUCKET"`
}

func GetDockerComposeYml(filePath string) *DockerComposeYml {
	v := viper.New()
	v.SetConfigFile(filePath)

	if err := v.ReadInConfig(); err != nil {
		log.Println(err)
	}

	dockComposeConfig := &DockerComposeYml{}
	if err := v.Unmarshal(dockComposeConfig); err != nil {
		log.Println(err)
	}

	influxdbPorts := strings.Split(dockComposeConfig.Services.Influxdb.Ports[0], ":")
	dockComposeConfig.Services.Influxdb.PortSrc, _ = strconv.Atoi(influxdbPorts[0])
	dockComposeConfig.Services.Influxdb.PortTarget, _ = strconv.Atoi(influxdbPorts[1])

	mongoPorts := strings.Split(dockComposeConfig.Services.Mongo.Ports[0], ":")
	dockComposeConfig.Services.Mongo.PortSrc, _ = strconv.Atoi(mongoPorts[0])
	dockComposeConfig.Services.Mongo.PortTarget, _ = strconv.Atoi(mongoPorts[1])

	return dockComposeConfig
}

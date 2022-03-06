package model

import (
	"io/ioutil"
	"path/filepath"

	logger "github.com/webdevolegkuprianov/server_rest_grpc/app/apiserver/logger"
	"gopkg.in/yaml.v2"
)

//config yaml struct
type Service struct {
	KafkaEcoSystem string `yaml:"kafka_eco_system"`
	Spec           struct {
		Ports struct {
			RestServerGrpcClient struct {
				BindAddr string `yaml:"bind_addr"`
			} `yaml:"rest_server_grpc_client"`
			GrpcServer struct {
				BindAddr string `yaml:"bind_addr"`
			} `yaml:"grpc_server"`
		} `yaml:"ports"`
		DB struct {
			Name              string `yaml:"name"`
			Host              string `yaml:"host"`
			Port              uint16 `yaml:"port"`
			User              string `yaml:"user"`
			Password          string `yaml:"password"`
			Database          string `yaml:"database"`
			MaxConnLifetime   int    `yaml:"max_conn_lifetime"`
			MaxConnIdletime   int    `yaml:"max_conn_idletime"`
			MaxConns          int32  `yaml:"max_conns"`
			MinConns          int32  `yaml:"min_conns"`
			HealthCheckPeriod int    `yaml:"health_check_period"`
		} `yaml:"db"`
		Jwt struct {
			TokenDecode string `yaml:"token"`
			LifeTerm    int    `yaml:"term"`
		} `yaml:"jwt"`
		Logs struct {
			Path string `yaml:"path"`
		} `yaml:"logs"`
	} `yaml:"spec"`
}

//New config
func NewConfig() (*Service, error) {

	var service *Service

	f, err := filepath.Abs("/root/config/kafka.yaml")
	if err != nil {
		logger.ErrorLogger.Println(err)
		return nil, err
	}

	y, err := ioutil.ReadFile(f)
	if err != nil {
		logger.ErrorLogger.Println(err)
		return nil, err
	}

	if err := yaml.Unmarshal(y, &service); err != nil {
		logger.ErrorLogger.Println(err)
		return nil, err
	}

	return service, nil

}

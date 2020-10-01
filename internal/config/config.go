package config

import (
	"os"
	"path/filepath"
	"sync"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

//Config struct holds grpc and Server configuration
type Config struct {
	// gRPC server start parameters section
	// gRPC is TCP port to listen by gRPC server
	GRPCPort     string
	DBHost       string
	DBPort       string
	DBUri        string
	DBUser       string
	DBPassword   string
	DBName       string
	GRPCServer   string
	ClientPort   string
	ClientServer string
}

var doOnce sync.Once

// ReadConfig will read the configuration from properties file
func ReadConfig() (*Config, error) {
	pwd, err := os.Getwd()
	if err != nil {
		logrus.Fatal("Unable to get Working Directory.. ", err)
	}
	var v *viper.Viper
	//fmt.Println(pwd)
	doOnce.Do(func() {

		v = viper.New()
		v.SetConfigName("config")
		v.SetConfigType("yaml")
		v.AddConfigPath(filepath.FromSlash(pwd + "/resources"))
		v.AutomaticEnv()
		err = v.ReadInConfig()

		if err != nil {
			logrus.Fatal("Unable to Read Config file ", err)
		}
	})

	config := &Config{
		GRPCPort:     v.GetString("grpcport"),
		DBUri:        v.GetString("mongodb.uri"),
		DBName:       v.GetString("mongodb.database"),
		GRPCServer:   v.GetString("grpcserver"),
		ClientPort:   v.GetString("clientport"),
		ClientServer: v.GetString("clientserver"),
	}

	return config, nil
}

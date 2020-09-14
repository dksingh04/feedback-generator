package config

import (
	"github.com/spf13/viper"
)

//Config struct holds grpc and Server configuration
type Config struct {
	// gRPC server start parameters section
	// gRPC is TCP port to listen by gRPC server
	GRPCPort   string
	DBHost     string
	DBPort     string
	DBUri      string
	DBUser     string
	DBPassword string
	DBName     string
	Server     string
}

// ReadConfig will read the configuration from properties file
func ReadConfig() (*Config, error) {
	v := viper.New()
	v.SetConfigName("config")
	v.SetConfigType("yaml")
	v.AddConfigPath("./internal/config")
	v.AutomaticEnv()
	err := v.ReadInConfig()

	if err != nil {
		return nil, err
	}

	config := &Config{
		GRPCPort: v.GetString("grpcport"),
		DBUri:    v.GetString("mongodb.uri"),
		DBName:   v.GetString("mongodb.database"),
		Server:   v.GetString("server"),
	}

	return config, nil
}

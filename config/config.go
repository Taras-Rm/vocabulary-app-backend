package config

import (
	"fmt"
	"os"

	"github.com/spf13/viper"
)

type AppConfig struct {
	Elastic  ElasticConfig  `yaml:"elastic"`
	Postgres PostgresConfig `yaml:"postgres"`
	Salt     string         `yaml:"salt"`
	AWS      AWSConfig      `yaml:"aws"`
}

type ElasticConfig struct {
	Host string `yaml:"host"`
	Port string `yaml:"port"`
}

type PostgresConfig struct {
	User     string `yaml:"user"`
	Port     string `yaml:"port"`
	Database string `yaml:"database"`
	Password string `yaml:"password"`
	Host     string `yaml:"host"`
}

type AWSConfig struct {
	Region   string `yaml:"region"`
	AccessId string `yaml:"accessId"`
	Secret   string `yaml:"secret"`
}

var Config AppConfig

func init() {
	env, _ := os.LookupEnv("ENV")
	fmt.Println("Envirq " + env)

	if env == "prod" {
		setProdConfig()
	} else {
		viper.AddConfigPath("./config")
		viper.SetConfigName("config")
		viper.SetConfigType("yaml")
		err := viper.ReadInConfig()
		if err != nil {
			fmt.Println(err)
		}

		err = viper.Unmarshal(&Config)
		if err != nil {
			fmt.Println(err)
		}
	}

}

func setProdConfig() {
	data, ok := os.LookupEnv("DB_NAME")
	if !ok {
		fmt.Println("can`t get env")
	}
	Config.Postgres.Database = data

	data, ok = os.LookupEnv("DB_USER")
	if !ok {
		fmt.Println("can`t get env")
	}
	Config.Postgres.User = data

	data, ok = os.LookupEnv("DB_PORT")
	if !ok {
		fmt.Println("can`t get env")
	}
	Config.Postgres.Port = data

	data, ok = os.LookupEnv("DB_PASSWORD")
	if !ok {
		fmt.Println("can`t get env")
	}
	Config.Postgres.Password = data

	data, ok = os.LookupEnv("SALT")
	if !ok {
		fmt.Println("can`t get env")
	}
	Config.Salt = data
}

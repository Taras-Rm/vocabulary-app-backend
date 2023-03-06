package config

import (
	"fmt"

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
}

type AWSConfig struct {
	Region   string `yaml:"region"`
	AccessId string `yaml:"accessId"`
	Secret   string `yaml:"secret"`
}

var Config AppConfig

func init() {
	if viper.Get("ENV") == "prod" {
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
	Config.Postgres.Database = viper.GetString("DB_NAME")
	Config.Postgres.User = viper.GetString("DB_USER")
	Config.Postgres.Port = viper.GetString("DB_PORT")
	Config.Postgres.Password = viper.GetString("DB_PASSWORD")

	Config.Salt = viper.GetString("SALT")

}

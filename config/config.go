package config

import (
	"fmt"
	"os"
	"strconv"

	"github.com/spf13/viper"
)

type AppConfig struct {
	Elastic  ElasticConfig  `yaml:"elastic"`
	Postgres PostgresConfig `yaml:"postgres"`
	Salt     string         `yaml:"salt"`
	AWS      AWSConfig      `yaml:"aws"`
	Hasher   Hasher         `yaml:"hasher"`
}

type ElasticConfig struct {
	Url      string `yaml:"url"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
}

type PostgresConfig struct {
	User     string `yaml:"user"`
	Port     string `yaml:"port"`
	Database string `yaml:"database"`
	Password string `yaml:"password"`
	Host     string `yaml:"host"`
}

type Hasher struct {
	Cost int `yaml:"cost"`
}

type AWSConfig struct {
	Region   string `yaml:"region"`
	AccessId string `yaml:"accessId"`
	Secret   string `yaml:"secret"`
}

var Config AppConfig

type Prod struct{}
type Local struct{}

type AppConfigInterface interface {
	GetEnvVariables() error
}

func (p *Local) GetEnvVariables() error {
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

	return nil
}

func (p *Prod) GetEnvVariables() error {
	// db
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

	data, ok = os.LookupEnv("DB_HOST")
	if !ok {
		fmt.Println("can`t get env")
	}
	Config.Postgres.Host = data

	// salt
	data, ok = os.LookupEnv("SALT")
	if !ok {
		fmt.Println("can`t get env")
	}
	Config.Salt = data

	// elasticsearch
	data, ok = os.LookupEnv("ES_URL")
	if !ok {
		fmt.Println("can`t get env")
	}
	Config.Elastic.Url = data

	data, ok = os.LookupEnv("ES_USERNAME")
	if !ok {
		fmt.Println("can`t get env")
	}
	Config.Elastic.Username = data

	data, ok = os.LookupEnv("ES_PASSWORD")
	if !ok {
		fmt.Println("can`t get env")
	}
	Config.Elastic.Password = data

	data, ok = os.LookupEnv("HS_HASH")
	if !ok {
		fmt.Println("can`t get env")
	}
	dataN, err := strconv.Atoi(data)
	if err != nil {
		fmt.Println("can`t parse env variable")
	}
	Config.Hasher.Cost = dataN

	data, ok = os.LookupEnv("AWS_ID")
	if !ok {
		fmt.Println("can`t get env")
	}
	Config.AWS.AccessId = data

	data, ok = os.LookupEnv("AWS_SECRET")
	if !ok {
		fmt.Println("can`t get env")
	}
	Config.AWS.Secret = data

	data, ok = os.LookupEnv("AWS_REGION")
	if !ok {
		fmt.Println("can`t get env")
	}
	Config.AWS.Region = data

	return nil
}

func init() {
	env, _ := os.LookupEnv("ENV")
	fmt.Println("Envirq " + env)

	var environment AppConfigInterface

	if IsProdEnv() {
		environment = &Prod{}
	} else {
		environment = &Local{}
	}

	err := environment.GetEnvVariables()
	if err != nil {
		fmt.Println(err.Error())
	}
}

func IsProdEnv() bool {
	env, _ := os.LookupEnv("ENV")
	return env == "prod"
}

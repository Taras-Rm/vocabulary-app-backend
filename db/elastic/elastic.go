package elastic

import (
	"crypto/tls"
	"errors"
	"fmt"
	"net/http"
	"time"
	"vacabulary/config"

	elastic "github.com/olivere/elastic/v7"
)

type ElasticClient struct {
	Client *elastic.Client
}

func NewElasticClient(cfg config.ElasticConfig) *ElasticClient {
	var client *elastic.Client

	if config.IsProdEnv() {
		client = getClient(cfg.Host, cfg.Port, cfg.Username, cfg.Password)
	} else {
		client = getClientLocal(cfg.Host, cfg.Port)
	}

	if client == nil {
		panic("missing elasticsearch connection")
	}

	return &ElasticClient{
		Client: client,
	}
}

func getClientLocal(host, port string) *elastic.Client {
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	var client *elastic.Client
	var err error

	if host == "" {
		host = "elasticsearch"
	}
	connectionUrl := fmt.Sprintf("http://%s:%s", host, port)

	for {
		client, err = elastic.NewSimpleClient(
			elastic.SetURL(connectionUrl),
			elastic.SetHealthcheck(true),
		)
		if err != nil {
			time.Sleep(3 * time.Second)
		} else {
			break
		}
	}

	return client
}

func getClient(host, port, username, password string) *elastic.Client {
	connectionUrl := fmt.Sprintf("http://%s:%s", host, port)

	client, err := elastic.NewSimpleClient(
		elastic.SetSniff(false),
		elastic.SetURL(connectionUrl),
		elastic.SetBasicAuth(username, password),
		elastic.SetHealthcheck(true),
	)
	if err != nil {
		panic(err)
	}

	return client
}

func (ec *ElasticClient) GetConnection() (*elastic.Client, error) {
	if ec.Client == nil {
		return nil, errors.New("missing elastic client")
	}

	return ec.Client, nil
}

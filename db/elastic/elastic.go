package elastic

import (
	"crypto/tls"
	"errors"
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
		client = getClient(cfg.Url, cfg.Username, cfg.Password)
	} else {
		client = getClientLocal(cfg.Url)
	}

	if client == nil {
		panic("missing elasticsearch connection")
	}

	return &ElasticClient{
		Client: client,
	}
}

func getClientLocal(url string) *elastic.Client {
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	var client *elastic.Client
	var err error

	for {
		client, err = elastic.NewSimpleClient(
			elastic.SetURL(url),
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

func getClient(url, username, password string) *elastic.Client {
	client, err := elastic.NewSimpleClient(
		elastic.SetSniff(false),
		elastic.SetURL(url),
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

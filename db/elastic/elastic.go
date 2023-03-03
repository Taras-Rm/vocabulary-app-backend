package elastic

import (
	"crypto/tls"
	"errors"
	"fmt"
	"net/http"
	"time"

	elastic "github.com/olivere/elastic/v7"
)

type ElasticClient struct {
	Client *elastic.Client
}

func NewElasticClient(host, port string) *ElasticClient {
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

	if client == nil {
		panic("missing elasticsearch connection")
	}

	return &ElasticClient{
		Client: client,
	}
}

func (ec *ElasticClient) GetClient() (*elastic.Client, error) {
	if ec.Client == nil {
		return nil, errors.New("missing elastic client")
	}

	return ec.Client, nil
}

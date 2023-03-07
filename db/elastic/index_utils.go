package elastic

import (
	"context"
	"errors"
)

func (ec *ElasticClient) CreateVocabularyIndices() error {
	wordsIndex := NewWordsIndex()

	err := ec.createIndicesIfNotExists(wordsIndex)
	if err != nil {
		return err
	}

	return nil
}

func (ec *ElasticClient) createIndicesIfNotExists(indices ...IndexInterface) error {
	client, err := ec.GetConnection()
	if err != nil {
		return err
	}
	ctx := context.Background()

	for _, index := range indices {
		indexName := index.GetName()
		indexMapping := index.GetMapping()

		exists, err := client.IndexExists(indexName).Do(ctx)
		if err != nil {
			return err
		}

		if !exists {
			result, err := client.CreateIndex(indexName).BodyString(indexMapping).Do(ctx)
			if err != nil {
				return err
			}

			if !result.Acknowledged {
				return errors.New("index acknowledged")
			}
		}
	}

	return nil
}

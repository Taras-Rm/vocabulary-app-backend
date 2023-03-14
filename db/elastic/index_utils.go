package elastic

import (
	"context"
	"errors"
	"fmt"
)

func (ec *ElasticClient) CreateUserWordsIndices(userId uint64) error {
	collectionWordsIndex, err := NewCollectionWordsIndex(CollectionWordsIndexContext{UserID: userId})
	if err != nil {
		return err
	}

	err = ec.createIndicesIfNotExists(collectionWordsIndex)
	if err != nil {
		return err
	}

	return nil
}

func (ec *ElasticClient) CreateCollectionAliases(userId uint64, collectionId uint64) error {
	collectionWordsIndex, err := NewCollectionWordsIndex(CollectionWordsIndexContext{UserID: userId})
	if err != nil {
		return err
	}

	collectionWordsAlias, err := NewCollectionWordsIndex(CollectionWordsIndexContext{UserID: userId, CollectionID: collectionId})
	if err != nil {
		return err
	}

	err = ec.createAliacesIfNotExists(collectionWordsIndex.GetName(), collectionWordsAlias)
	if err != nil {
		return err
	}

	return nil
}

func (ec *ElasticClient) createIndicesIfNotExists(indices ...CollectionWordsIndexInterface) error {
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

func (ec *ElasticClient) createAliacesIfNotExists(index string, aliases ...CollectionWordsIndexInterface) error {
	client, err := ec.GetConnection()
	if err != nil {
		return err
	}
	ctx := context.Background()

	aliasesResult, err := client.Aliases().Index(index).Do(ctx)
	if err != nil {
		return err
	}

	for _, alias := range aliases {
		aliasName := alias.GetName()

		if aliasesResult.Indices[index].HasAlias(aliasName) {
			fmt.Printf("Alias %s for index %s already exists - skip", aliasName, index)
		} else {
			aliasFilter, err := alias.GetFilter()
			if err != nil {
				return err
			}

			if aliasFilter == nil {
				return fmt.Errorf("missing alias filter for alias %s", aliasName)
			}

			response, err := client.Alias().AddWithFilter(index, aliasName, aliasFilter).Do(ctx)
			if err != nil {
				return err
			}

			if !response.Acknowledged {
				return errors.New("index acknowledged")
			}
		}
	}

	return nil
}

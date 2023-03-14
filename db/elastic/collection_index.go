package elastic

import (
	"errors"
	"fmt"

	"github.com/olivere/elastic/v7"
)

type CollectionWordsIndex struct {
	name    string
	mapping string
	ctx     *CollectionWordsIndexContext
}

type CollectionWordsIndexContext struct {
	UserID       uint64
	CollectionID uint64
}

type CollectionWordsIndexInterface interface {
	GetName() string
	GetMapping() string
	GetFilter() (elastic.Query, error)
}

func (wi *CollectionWordsIndex) GetMapping() string {
	return wi.mapping
}

func (wi *CollectionWordsIndex) GetName() string {
	return wi.name
}

func (wi *CollectionWordsIndex) GetFilter() (elastic.Query, error) {
	if wi.ctx == nil || wi.ctx.CollectionID == 0 {
		return nil, errors.New("mixxing collectionId")
	}

	return elastic.NewTermQuery("collection_id", wi.ctx.CollectionID), nil
}

func NewCollectionWordsIndex(ctx CollectionWordsIndexContext) (*CollectionWordsIndex, error) {
	collection_words_index := "collection_words"

	if ctx.UserID == 0 {
		return nil, fmt.Errorf("missed userId for CollectionWordsIndex")
	}

	if ctx.CollectionID == 0 {
		collection_words_index = fmt.Sprintf("%s-%v", collection_words_index, ctx.UserID)
	} else {
		collection_words_index = fmt.Sprintf("%s-%v-%v", collection_words_index, ctx.UserID, ctx.CollectionID)
	}

	const collectionWordsMapping = `
	{
		"settings":{
			"number_of_shards": 1
		},
		"mappings":{
			"properties":{
				"collection_id":{
					"type":"integer"
				},
				"word":{
					"type":"text"
				},
				"translation":{
					"type":"text"
				},
				"part_of_speech":{
					"type":"text"
				},
				"scentance":{
					"type":"text"
				},
				"created_at":{
					"type":"date"
				}
			}
		}
	}`

	return &CollectionWordsIndex{
		name:    collection_words_index,
		mapping: collectionWordsMapping,
		ctx:     &ctx,
	}, nil
}

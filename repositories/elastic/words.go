package elastic

import (
	"context"
	"encoding/json"
	"errors"
	"reflect"
	"time"
	myElastic "vacabulary/db/elastic"
	"vacabulary/models"

	"github.com/olivere/elastic/v7"
)

type collectionWordsRepo struct {
	client *elastic.Client
}

type Words interface {
	Create(word models.Word, wordsCtx CollectionWordsOperationCtx) error
	BulkCreate(word []models.Word, wordsCtx CollectionWordsOperationCtx) error
	Update(word models.Word, wordsCtx CollectionWordsOperationCtx) error
	DeleteById(id string, wordsCtx CollectionWordsOperationCtx) error
	Get(origin string, wordsCtx CollectionWordsOperationCtx) (*models.Word, error)
	GetById(origin string, wordsCtx CollectionWordsOperationCtx) (*models.Word, error)
	GetByTranslation(translation string, wordsCtx CollectionWordsOperationCtx) (*models.Word, error)
	GetAll(size, page uint64, wordsCtx CollectionWordsOperationCtx) ([]models.Word, uint64, error)
	GetByWords(words []string, wordsCtx CollectionWordsOperationCtx) ([]models.Word, error)
	Search(settings models.SearchSettings, wordsCtx CollectionWordsOperationCtx) ([]models.Word, error)
}

func NewCollectionWordsRepo(client *elastic.Client) Words {
	return &collectionWordsRepo{
		client: client,
	}
}

type ElasticWord struct {
	CollectionId uint64    `json:"collection_id"`
	Word         string    `json:"word"`
	Translation  string    `json:"translation"`
	PartOfSpeech string    `json:"part_of_speech"`
	Scentance    string    `json:"scentance"`
	CreatedAt    time.Time `json:"created_at"`
}

type CollectionWordsOperationCtx struct {
	UserId       uint64
	CollectionId uint64
}

func (r *collectionWordsRepo) Create(word models.Word, wordsCtx CollectionWordsOperationCtx) error {
	index, err := r.getIndex(wordsCtx)
	if err != nil {
		return err
	}

	elasticWord := ElasticWord{
		CollectionId: word.CollectionId,
		Word:         word.Word,
		Translation:  word.Translation,
		PartOfSpeech: word.PartOfSpeech,
		Scentance:    word.Scentance,
		CreatedAt:    time.Now(),
	}

	ctx := context.Background()
	_, err = r.client.Index().Index(index.GetName()).BodyJson(elasticWord).Refresh("true").Do(ctx)
	if err != nil {
		return err
	}

	return nil
}

func (r *collectionWordsRepo) BulkCreate(words []models.Word, wordsCtx CollectionWordsOperationCtx) error {
	index, err := r.getIndex(wordsCtx)
	if err != nil {
		return err
	}

	elasticWords := []ElasticWord{}
	for _, word := range words {
		elasticWords = append(elasticWords, ElasticWord{
			CollectionId: word.CollectionId,
			Word:         word.Word,
			Translation:  word.Translation,
			PartOfSpeech: word.PartOfSpeech,
			Scentance:    word.Scentance,
			CreatedAt:    time.Now(),
		})
	}

	bulk := r.client.Bulk()
	for _, eWord := range elasticWords {
		req := elastic.NewBulkIndexRequest()

		req.Index(index.GetName())
		req.Doc(eWord)

		bulk.Add(req)
	}

	ctx := context.Background()
	_, err = bulk.Refresh("true").Do(ctx)
	if err != nil {
		return err
	}

	return nil
}

func (r *collectionWordsRepo) DeleteById(id string, wordsCtx CollectionWordsOperationCtx) error {
	index, err := r.getIndex(wordsCtx)
	if err != nil {
		return err
	}

	ctx := context.Background()
	_, err = r.client.Delete().Index(index.GetName()).Refresh("true").Id(id).Do(ctx)
	if err != nil {
		return err
	}

	return nil
}

func (r *collectionWordsRepo) Get(word string, wordsCtx CollectionWordsOperationCtx) (*models.Word, error) {
	index, err := r.getIndex(wordsCtx)
	if err != nil {
		return nil, err
	}

	ctx := context.Background()

	query := elastic.NewBoolQuery() // ("word", word)
	q1 := elastic.NewMatchQuery("word", word)
	// q2 := elastic.NewMatchQuery("collection_id", collectionId)

	query.Must(q1)

	searchResult, err := r.client.Search().Index(index.GetName()).Query(query).Do(ctx)
	if err != nil {
		return nil, err
	}

	var findedWord *models.Word
	for _, word := range searchResult.Each(reflect.TypeOf(ElasticWord{})) {
		word, ok := word.(ElasticWord)
		if !ok {
			continue
		}

		findedWord = &models.Word{
			CollectionId: word.CollectionId,
			Word:         word.Word,
			Translation:  word.Translation,
			PartOfSpeech: word.PartOfSpeech,
			Scentance:    word.Scentance,
			CreatedAt:    word.CreatedAt,
		}
	}

	return findedWord, nil
}

func (r *collectionWordsRepo) GetById(id string, wordsCtx CollectionWordsOperationCtx) (*models.Word, error) {
	index, err := r.getIndex(wordsCtx)
	if err != nil {
		return nil, err
	}

	ctx := context.Background()

	searchResult, err := r.client.Get().Index(index.GetName()).Id(id).Do(ctx)
	if err != nil {
		return nil, err
	}

	var word ElasticWord
	err = json.Unmarshal(searchResult.Source, &word)
	if err != nil {
		return nil, err
	}

	findedWord := models.Word{
		Id:           searchResult.Id,
		CollectionId: word.CollectionId,
		Word:         word.Word,
		Translation:  word.Translation,
		PartOfSpeech: word.PartOfSpeech,
		Scentance:    word.Scentance,
		CreatedAt:    word.CreatedAt,
	}

	return &findedWord, nil
}

func (r *collectionWordsRepo) GetByTranslation(translation string, wordsCtx CollectionWordsOperationCtx) (*models.Word, error) {
	index, err := r.getIndex(wordsCtx)
	if err != nil {
		return nil, err
	}

	ctx := context.Background()

	query := elastic.NewTermQuery("translation", translation)

	searchResult, err := r.client.Search().Index(index.GetName()).Query(query).Do(ctx)
	if err != nil {
		return nil, err
	}
	var words []models.Word
	for _, word := range searchResult.Each(reflect.TypeOf(ElasticWord{})) {
		word, ok := word.(ElasticWord)
		if !ok {
			continue
		}

		words = append(words, models.Word{
			CollectionId: word.CollectionId,
			Word:         word.Word,
			Translation:  word.Translation,
			PartOfSpeech: word.PartOfSpeech,
			Scentance:    word.Scentance,
			CreatedAt:    word.CreatedAt,
		})
	}

	if len(words) == 0 {
		return nil, nil
	}

	return &words[0], nil
}

func (r *collectionWordsRepo) GetAll(size, page uint64, wordsCtx CollectionWordsOperationCtx) ([]models.Word, uint64, error) {
	index, err := r.getIndex(wordsCtx)
	if err != nil {
		return nil, 0, err
	}

	ctx := context.Background()

	// query := elastic.NewMatchQuery("collection_id", collectionId)

	var searchResult *elastic.SearchResult

	if size == 0 && page == 0 {
		// get all words
		searchResult, err = r.client.Search().Index(index.GetName()).Size(1000).Do(ctx)
	} else {
		from := size * (page - 1)
		searchResult, err = r.client.Search().Index(index.GetName()).Size(int(size)).From(int(from)).Do(ctx)
	}
	if err != nil {
		return nil, 0, err
	}

	totalHits := uint64(searchResult.Hits.TotalHits.Value)

	var words []models.Word
	for _, hit := range searchResult.Hits.Hits {
		var word ElasticWord
		err := json.Unmarshal(hit.Source, &word)
		if err != nil {
			continue
		}

		words = append(words, models.Word{
			Id:           hit.Id,
			CollectionId: word.CollectionId,
			Word:         word.Word,
			Translation:  word.Translation,
			PartOfSpeech: word.PartOfSpeech,
			Scentance:    word.Scentance,
			CreatedAt:    word.CreatedAt,
		})
	}

	return words, totalHits, nil
}

func (r *collectionWordsRepo) Update(word models.Word, wordsCtx CollectionWordsOperationCtx) error {
	index, err := r.getIndex(wordsCtx)
	if err != nil {
		return err
	}

	elasticWord := ElasticWord{
		CollectionId: word.CollectionId,
		Word:         word.Word,
		Translation:  word.Translation,
		PartOfSpeech: word.PartOfSpeech,
		Scentance:    word.Scentance,
		CreatedAt:    word.CreatedAt,
	}

	ctx := context.Background()
	_, err = r.client.Update().Index(index.GetName()).Refresh("true").Doc(elasticWord).Id(word.Id).Do(ctx)
	if err != nil {
		return err
	}

	return nil
}

func (r *collectionWordsRepo) GetByWords(words []string, wordsCtx CollectionWordsOperationCtx) ([]models.Word, error) {
	index, err := r.getIndex(wordsCtx)
	if err != nil {
		return nil, err
	}

	ctx := context.Background()

	wordsForSearch := make([]interface{}, len(words))
	for index, value := range words {
		wordsForSearch[index] = value
	}

	query := elastic.NewBoolQuery()

	q1 := elastic.NewTermsQuery("word", wordsForSearch...)
	// q2 := elastic.NewMatchQuery("collection_id", collectionId)

	query.Must(q1)

	searchResult, err := r.client.Search().Index(index.GetName()).Query(query).Do(ctx)
	if err != nil {
		return nil, err
	}
	var findedWords []models.Word
	for _, word := range searchResult.Each(reflect.TypeOf(ElasticWord{})) {
		word, ok := word.(ElasticWord)
		if !ok {
			continue
		}

		findedWords = append(findedWords, models.Word{
			CollectionId: word.CollectionId,
			Word:         word.Word,
			Translation:  word.Translation,
			PartOfSpeech: word.PartOfSpeech,
			Scentance:    word.Scentance,
			CreatedAt:    word.CreatedAt,
		})
	}

	if len(words) == 0 {
		return nil, nil
	}

	return findedWords, nil
}

func (r *collectionWordsRepo) Search(settings models.SearchSettings, wordsCtx CollectionWordsOperationCtx) ([]models.Word, error) {
	index, err := r.getIndex(wordsCtx)
	if err != nil {
		return nil, err
	}
	ctx := context.Background()

	query := elastic.NewBoolQuery()

	q1 := elastic.NewWildcardQuery(settings.SearchBy, "*"+settings.TextForSearch+"*")
	// q2 := elastic.NewMatchQuery("collection_id", collectionId)

	query.Must(q1)

	if len(settings.PartsOfSpeech) != 0 {
		partsOfSpeechArr := make([]interface{}, len(settings.PartsOfSpeech))
		for index, value := range settings.PartsOfSpeech {
			partsOfSpeechArr[index] = value
		}
		q3 := elastic.NewTermsQuery("part_of_speech", partsOfSpeechArr...)

		query.Must(q3)
	}

	searchResult, err := r.client.Search().Index(index.GetName()).Query(query).Do(ctx)
	if err != nil {
		return nil, err
	}
	var words []models.Word
	for _, hit := range searchResult.Hits.Hits {
		var word ElasticWord
		err := json.Unmarshal(hit.Source, &word)
		if err != nil {
			continue
		}

		words = append(words, models.Word{
			Id:           hit.Id,
			CollectionId: word.CollectionId,
			Word:         word.Word,
			Translation:  word.Translation,
			PartOfSpeech: word.PartOfSpeech,
			Scentance:    word.Scentance,
			CreatedAt:    word.CreatedAt,
		})
	}

	if len(words) == 0 {
		return nil, nil
	}

	return words, nil
}

func (r *collectionWordsRepo) getIndex(ctx CollectionWordsOperationCtx) (*myElastic.CollectionWordsIndex, error) {
	index, err := myElastic.NewCollectionWordsIndex(myElastic.CollectionWordsIndexContext{UserID: ctx.UserId, CollectionID: ctx.CollectionId})
	if err != nil {
		return nil, errors.New("failed to get collection words index")
	}

	return index, nil
}

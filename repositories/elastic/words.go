package elastic

import (
	"context"
	"encoding/json"
	"reflect"
	"time"
	myElastic "vacabulary/db/elastic"
	"vacabulary/models"

	"github.com/olivere/elastic/v7"
)

type wordsRepo struct {
	client *elastic.Client
}

type Words interface {
	Create(word models.Word) error
	BulkCreate(word []models.Word) error
	Update(word models.Word) error
	DeleteById(id string) error
	Get(origin string, collectionId uint64) (*models.Word, error)
	GetById(origin string) (*models.Word, error)
	GetByTranslation(translation string) (*models.Word, error)
	GetAll(collectionId uint64, size, page uint64) ([]models.Word, uint64, error)
	GetByWords(words []string, collectionId uint64) ([]models.Word, error)
	Search(collectionId uint64, settings models.SearchSettings) ([]models.Word, error)
}

func NewWordsRepo(client *elastic.Client) Words {
	return &wordsRepo{
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

func (r *wordsRepo) Create(word models.Word) error {
	index := r.getIndex()

	elasticWord := ElasticWord{
		CollectionId: word.CollectionId,
		Word:         word.Word,
		Translation:  word.Translation,
		PartOfSpeech: word.PartOfSpeech,
		Scentance:    word.Scentance,
		CreatedAt:    time.Now(),
	}

	ctx := context.Background()
	_, err := r.client.Index().Index(index.GetName()).BodyJson(elasticWord).Refresh("true").Do(ctx)
	if err != nil {
		return err
	}

	return nil
}

func (r *wordsRepo) BulkCreate(words []models.Word) error {
	index := r.getIndex()

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
	_, err := bulk.Refresh("true").Do(ctx)
	if err != nil {
		return err
	}

	return nil
}

func (r *wordsRepo) DeleteById(id string) error {
	index := r.getIndex()

	ctx := context.Background()
	_, err := r.client.Delete().Index(index.GetName()).Refresh("true").Id(id).Do(ctx)
	if err != nil {
		return err
	}

	return nil
}

func (r *wordsRepo) Get(word string, collectionId uint64) (*models.Word, error) {
	index := r.getIndex()

	ctx := context.Background()

	query := elastic.NewBoolQuery() // ("word", word)
	q1 := elastic.NewMatchQuery("word", word)
	q2 := elastic.NewMatchQuery("collection_id", collectionId)

	query.Must(q1, q2)

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

func (r *wordsRepo) GetById(id string) (*models.Word, error) {
	index := r.getIndex()

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

func (r *wordsRepo) GetByTranslation(translation string) (*models.Word, error) {
	index := r.getIndex()

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

func (r *wordsRepo) GetAll(collectionId uint64, size, page uint64) ([]models.Word, uint64, error) {
	index := r.getIndex()

	ctx := context.Background()

	query := elastic.NewMatchQuery("collection_id", collectionId)

	var searchResult *elastic.SearchResult
	var err error

	if size == 0 && page == 0 {
		// get all words
		searchResult, err = r.client.Search().Index(index.GetName()).Query(query).Size(1000).Do(ctx)
	} else {
		from := size * (page - 1)
		searchResult, err = r.client.Search().Index(index.GetName()).Query(query).Size(int(size)).From(int(from)).Do(ctx)
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

func (r *wordsRepo) Update(word models.Word) error {
	index := r.getIndex()

	elasticWord := ElasticWord{
		CollectionId: word.CollectionId,
		Word:         word.Word,
		Translation:  word.Translation,
		PartOfSpeech: word.PartOfSpeech,
		Scentance:    word.Scentance,
		CreatedAt:    word.CreatedAt,
	}

	ctx := context.Background()
	_, err := r.client.Update().Index(index.GetName()).Refresh("true").Doc(elasticWord).Id(word.Id).Do(ctx)
	if err != nil {
		return err
	}

	return nil
}

func (r *wordsRepo) GetByWords(words []string, collectionId uint64) ([]models.Word, error) {
	index := r.getIndex()

	ctx := context.Background()

	wordsForSearch := make([]interface{}, len(words))
	for index, value := range words {
		wordsForSearch[index] = value
	}

	query := elastic.NewBoolQuery()

	q1 := elastic.NewTermsQuery("word", wordsForSearch...)
	q2 := elastic.NewMatchQuery("collection_id", collectionId)

	query.Must(q1, q2)

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

func (r *wordsRepo) Search(collectionId uint64, settings models.SearchSettings) ([]models.Word, error) {
	index := r.getIndex()

	ctx := context.Background()

	query := elastic.NewBoolQuery()

	q1 := elastic.NewWildcardQuery(settings.SearchBy, "*"+settings.TextForSearch+"*")
	q2 := elastic.NewMatchQuery("collection_id", collectionId)

	query.Must(q1, q2)

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

func (r *wordsRepo) getIndex() *myElastic.WordsIndex {
	index := myElastic.NewWordsIndex()

	return index
}

package elastic

type WordsIndex struct {
	name    string
	mapping string
}

type IndexInterface interface {
	GetName() string
	GetMapping() string
}

func (wi *WordsIndex) GetMapping() string {
	return wi.mapping
}

func (wi *WordsIndex) GetName() string {
	return wi.name
}

func NewWordsIndex() *WordsIndex {
	words_index := "words"

	const wordMapping = `
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

	return &WordsIndex{
		name:    words_index,
		mapping: wordMapping,
	}
}

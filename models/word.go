package models

import "time"

type Word struct {
	Id           string    `json:"id"`
	CollectionId uint64    `json:"collectionId"`
	Word         string    `json:"word"`
	Translation  string    `json:"translation"`
	PartOfSpeech string    `json:"partOfSpeech"`
	Scentance    string    `json:"scentance"`
	CreatedAt    time.Time `json:"createdAt"`
}

type SearchSettings struct {
	TextForSearch string   `json:"textForSearch"`
	SearchBy      string   `json:"searchBy"`
	PartsOfSpeech []string `json:"partsOfSpeech"`
}

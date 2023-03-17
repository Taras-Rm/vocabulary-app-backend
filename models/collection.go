package models

import "time"

type Collection struct {
	Id        uint64    `json:"id"`
	Name      string    `json:"name"`
	OwnerId   uint64    `json:"ownerId"`
	CreatedAt time.Time `json:"createdAt"`
	Words     []Word    `json:"words"`
	LangFrom  string    `json:"langFrom"`
	LangTo    string    `json:"langTo"`
}

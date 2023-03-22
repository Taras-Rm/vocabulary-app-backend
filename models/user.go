package models

import "time"

type User struct {
	Id        uint64    `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Password  string    `json:"-"`
	CreatedAt time.Time `json:"createdAt"`
	IsSuper   bool      `json:"isSuper"`

	Settings *UserSettings `json:"settings"`
}

type UserSettings struct {
	Id       uint64 `json:"id"`
	UserId   uint64 `json:"userId"`
	Language string `json:"language"`
}

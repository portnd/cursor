package models

import "time"

type Discord struct {
	Username string   `json:"username"`
	Content  string   `json:"content"`
	Embeds   []Embeds `json:"embeds"`
}

type Embeds struct {
	Title  string   `json:"title"`
	Color  int      `json:"color"`
	Fields []Fields `json:"fields"`
}

type Fields struct {
	Name   string `json:"name"`
	Value  string `json:"value"`
	Inline bool   `json:"inline"`
}

type MongoDbLog struct {
	IsSuccess bool      `json:"is_success"`
	DateTime  time.Time `json:"date_time"`
	Title     string    `json:"title"`
	Message   string    `json:"message"`
}

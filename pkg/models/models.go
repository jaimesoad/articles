package models

type JSON map[string]interface{}

type Article struct {
	ReadTime int    `json:"readTime"`
	Title    string `json:"title"`
	Author   string `json:"author"`
	HTML     string `json:"-"`
	Date     string `json:"date"`
	Filename string	`json:"filename"`
}

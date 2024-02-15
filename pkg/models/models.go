package models

type JSON map[string]interface{}

type Article struct {
	ReadTime int
	Title    string
	Author   string
	HTML     string
	Date     string
	Filename string
}
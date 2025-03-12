package models

type Album struct {
	ID     int  `json:"id"`
	Title  string  `json:"title"`
	Artist string  `json:"artist"`
	Link   string `json:"link"`
}

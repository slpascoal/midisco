package models

type Album struct {
	ID     int    `json:"id" gorm:"primaryKey;autoIncrement"`
	Title  string `json:"title" gorm:"size:255;not null"`
	Artist string `json:"artist" gorm:"size:255;not null"`
	Link   string `json:"link" gorm:"size:255"`
}

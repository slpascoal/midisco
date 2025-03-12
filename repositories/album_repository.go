package repositories

import (
	"database/sql"
	"errors"
	"log"

	"midisco-api/models"
)

type AlbumRepository interface {
	GetAll() ([]models.Album, error)
	GetByID(id int) (models.Album, error)
	Create(album models.Album) error
	Update(album models.Album) error
	Delete(id int) error
}

type albumRepository struct {
	db *sql.DB
}

func NewAlbumRepository(db *sql.DB) AlbumRepository {
	err := createAlbumsTableIfNotExists(db)
	if err != nil {
		log.Fatal("Erro ao criar tabela de álbuns:", err)
	}
	return &albumRepository{db: db}
}

func createAlbumsTableIfNotExists(db *sql.DB) error {
	query := `
	CREATE TABLE IF NOT EXISTS albums (
		id INT AUTO_INCREMENT PRIMARY KEY,
		title VARCHAR(255) NOT NULL,
		artist VARCHAR(255) NOT NULL,
		link VARCHAR(255)
	);`
	_, err := db.Exec(query)
	return err
}

func (r *albumRepository) GetAll() ([]models.Album, error) {
	rows, err := r.db.Query("SELECT id, title, artist, link FROM albums")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var albums []models.Album
	for rows.Next() {
		var a models.Album
		if err := rows.Scan(&a.ID, &a.Title, &a.Artist, &a.Link); err != nil {
			return nil, err
		}
		albums = append(albums, a)
	}
	return albums, nil
}

func (r *albumRepository) GetByID(id int) (models.Album, error) {
	var album models.Album
	query := "SELECT id, title, artist, link FROM albums WHERE id = ?"
	err := r.db.QueryRow(query, id).Scan(&album.ID, &album.Title, &album.Artist, &album.Link)
	if err == sql.ErrNoRows {
		return album, errors.New("album not found")
	}
	return album, err
}

func (r *albumRepository) Create(album models.Album) error {
	query := "INSERT INTO albums (title, artist, link) VALUES (?, ?, ?)"
	_, err := r.db.Exec(query, album.Title, album.Artist, album.Link)
	return err
}

func (r *albumRepository) Update(album models.Album) error {
	query := "UPDATE albums SET title = ?, artist = ?, link = ? WHERE id = ?"
	res, err := r.db.Exec(query, album.Title, album.Artist, album.Link, album.ID)
	if err != nil {
		return err
	}
	affected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if affected == 0 {
		return errors.New("album not found")
	}
	return nil
}

func (r *albumRepository) Delete(id int) error {
	query := "DELETE FROM albums WHERE id = ?"
	res, err := r.db.Exec(query, id)
	if err != nil {
		return err
	}
	affected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if affected == 0 {
		return errors.New("album not found")
	}
	return nil
}

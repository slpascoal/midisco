package repositories

import (
	"errors"
	"log"

	"midisco-api/models"

	"gorm.io/gorm"
)

type AlbumRepository interface {
	GetAll() ([]models.Album, error)
	GetByID(id int) (models.Album, error)
	Create(album models.Album) error
	Update(album models.Album) error
	Delete(id int) error
}

type albumRepository struct {
	db *gorm.DB
}

func NewAlbumRepository(db *gorm.DB) AlbumRepository {
	// Auto-migra o schema (cria a tabela se não existir)
	if err := db.AutoMigrate(&models.Album{}); err != nil {
		log.Fatal("Erro ao migrar a tabela de álbuns:", err)
	}
	return &albumRepository{db: db}
}

func (r *albumRepository) GetAll() ([]models.Album, error) {
	var albums []models.Album
	result := r.db.Find(&albums)
	return albums, result.Error
}

func (r *albumRepository) GetByID(id int) (models.Album, error) {
	var album models.Album
	result := r.db.First(&album, id)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return album, errors.New("album not found")
	}
	return album, result.Error
}

func (r *albumRepository) Create(album models.Album) error {
	result := r.db.Create(&album)
	return result.Error
}

func (r *albumRepository) Update(album models.Album) error {
	var existing models.Album
	if err := r.db.First(&existing, album.ID).Error; err != nil {
		return errors.New("album not found")
	}

	result := r.db.Model(&existing).Updates(album)
	return result.Error
}

func (r *albumRepository) Delete(id int) error {
	result := r.db.Delete(&models.Album{}, id)
	if result.RowsAffected == 0 {
		return errors.New("album not found")
	}
	return result.Error
}

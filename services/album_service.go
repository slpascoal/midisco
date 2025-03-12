package services

import (
	"midisco-api/models"
	"midisco-api/repositories"
)

type AlbumService interface {
	GetAlbums() ([]models.Album, error)
	GetAlbumByID(id int) (models.Album, error)
	CreateAlbum(album models.Album) error
	UpdateAlbum(album models.Album) error
	DeleteAlbum(id int) error
}

type albumService struct {
	repo repositories.AlbumRepository
}

func NewAlbumService(repo repositories.AlbumRepository) AlbumService {
	return &albumService{repo: repo}
}

func (s *albumService) GetAlbums() ([]models.Album, error) {
	return s.repo.GetAll()
}

func (s *albumService) GetAlbumByID(id int) (models.Album, error) {
	return s.repo.GetByID(id)
}

func (s *albumService) CreateAlbum(album models.Album) error {
	return s.repo.Create(album)
}

func (s *albumService) UpdateAlbum(album models.Album) error {
	return s.repo.Update(album)
}

func (s *albumService) DeleteAlbum(id int) error {
	return s.repo.Delete(id)
}

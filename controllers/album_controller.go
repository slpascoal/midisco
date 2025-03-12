package controllers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"midisco-api/models"
	"midisco-api/services"
)

type AlbumController struct {
	service services.AlbumService
}

func NewAlbumController(service services.AlbumService) *AlbumController {
	return &AlbumController{service: service}
}

func (ac *AlbumController) GetAlbums(c *gin.Context) {
	albums, err := ac.service.GetAlbums()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, albums)
}

func (ac *AlbumController) GetAlbumByID(c *gin.Context) {
	idParam := c.Param("id")

	id, err := strconv.Atoi(idParam)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
        return
    }

	album, err := ac.service.GetAlbumByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Álbum não encontrado"})
		return
	}
	c.JSON(http.StatusOK, album)
}

func (ac *AlbumController) CreateAlbum(c *gin.Context) {
	var album models.Album
	if err := c.ShouldBindJSON(&album); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := ac.service.CreateAlbum(album); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"message": "Álbum criado com sucesso"})
}

func (ac *AlbumController) UpdateAlbum(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	var album models.Album
	if err := c.ShouldBindJSON(&album); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	album.ID = id

	if err := ac.service.UpdateAlbum(album); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Álbum atualizado com sucesso"})
}

func (ac *AlbumController) DeleteAlbum(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	if err := ac.service.DeleteAlbum(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Álbum deletado com sucesso"})
}

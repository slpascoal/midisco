package routes

import (
	"github.com/gin-gonic/gin"
	"midisco-api/controllers"
)

func SetupRouter(albumController *controllers.AlbumController) *gin.Engine {
	router := gin.Default()

	router.GET("/albums", albumController.GetAlbums)
	router.GET("/albums/:id", albumController.GetAlbumByID)
	router.POST("/albums", albumController.CreateAlbum)
	router.PUT("/albums/:id", albumController.UpdateAlbum)
	router.DELETE("/albums/:id", albumController.DeleteAlbum)

	return router
}
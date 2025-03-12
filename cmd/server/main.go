package main

import (
	"fmt"
	"log"
	"time"

	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	"midisco-api/config"
	"midisco-api/controllers"
	"midisco-api/repositories"
	"midisco-api/routes"
	"midisco-api/services"
)

func initDB() *gorm.DB {
	cfg, err := config.Load()
	if err != nil {
		log.Fatal("Erro ao carregar configurações:", err)
	}

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true",
		cfg.DBUser, cfg.DBPassword, cfg.DBHost, cfg.DBPort, cfg.DBName)

	var db *gorm.DB
	maxAttempts := 10
	for attempts := 1; attempts <= maxAttempts; attempts++ {
		db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
		if err == nil {
			sqlDB, _ := db.DB()
			err = sqlDB.Ping()
			if err == nil {
				log.Println("Banco de dados conectado com sucesso!")
				break
			}
		}
		log.Printf("MySQL não está pronto (tentativa %d/%d): %v\n", attempts, maxAttempts, err)
		time.Sleep(2 * time.Second)
	}
	if err != nil {
		log.Fatal("Erro ao conectar ao banco de dados:", err)
	}

	return db
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("Nenhum arquivo .env encontrado. Continuando com as variáveis de ambiente.")
	}

	db := initDB()

	albumRepo := repositories.NewAlbumRepository(db)
	albumService := services.NewAlbumService(albumRepo)
	albumController := controllers.NewAlbumController(albumService)

	router := routes.SetupRouter(albumController)
	router.Run(":8080")
}

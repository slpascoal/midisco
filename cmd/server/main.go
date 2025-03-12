package main

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/joho/godotenv"
	_ "github.com/go-sql-driver/mysql"

	"midisco-api/controllers"
	"midisco-api/repositories"
	"midisco-api/routes"
	"midisco-api/services"
	"midisco-api/config"
)

func initDB() *sql.DB {
	cfg, err := config.Load()
	if err != nil {
		log.Fatal("Erro ao carregar configurações:", err)
	}

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true", 
		cfg.DBUser, cfg.DBPassword, cfg.DBHost, cfg.DBPort, cfg.DBName)

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal("Erro ao conectar ao MySQL:", err)
	}

	maxAttempts := 10
	for attempts := 1; attempts <= maxAttempts; attempts++ {
		err = db.Ping()
		if err == nil {
			log.Println("Banco de dados conectado com sucesso!")
			break
		}
		log.Printf("MySQL não está pronto (tentativa %d/%d): %v\n", attempts, maxAttempts, err)
		time.Sleep(2 * time.Second)
	}

	if err := db.Ping(); err != nil {
		log.Fatal("Erro ao conectar ao banco de dados:", err)
	}

	return db
}

func main() {
	// Carrega o arquivo .env
	err := godotenv.Load()
	if err != nil {
		log.Println("Nenhum arquivo .env encontrado. Continuando com as variáveis de ambiente.")
	}

	db := initDB()

	// Injeção de dependências
	albumRepo := repositories.NewAlbumRepository(db)
	albumService := services.NewAlbumService(albumRepo)
	albumController := controllers.NewAlbumController(albumService)

	// Configuração das rotas e inicialização do servidor
	router := routes.SetupRouter(albumController)
	router.Run(":8080")
}

package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	_ "github.com/go-sql-driver/mysql"
)

// album representa a estrutura dos dados de um álbum.
type album struct {
	ID     string  `json:"id"`
	Title  string  `json:"title"`
	Artist string  `json:"artist"`
	Price  float64 `json:"price"`
}

var db *sql.DB

// initDB realiza a conexão com o MySQL e cria a tabela 'albums' se ela não existir.
func initDB() {
	// Recupera as variáveis de ambiente para a conexão com o MySQL
	mysqlUser := "root"
	mysqlPassword := os.Getenv("MYSQL_PASSWORD")
	mysqlHost := os.Getenv("MYSQL_HOST")
	mysqlPort := os.Getenv("MYSQL_PORT")
	mysqlDB := os.Getenv("MYSQL_DB")

	// Data Source Name (DSN)
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true", mysqlUser, mysqlPassword, mysqlHost, mysqlPort, mysqlDB)
	var err error
	db, err = sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal("Erro ao conectar ao MySQL:", err)
	}

	// Tenta fazer ping no banco até que esteja disponível
	maxAttempts := 10
	for attempts := 1; attempts <= maxAttempts; attempts++ {
		err = db.Ping()
		if err == nil {
			break
		}
		log.Printf("MySQL não está pronto (tentativa %d/%d): %v\n", attempts, maxAttempts, err)
		time.Sleep(2 * time.Second)
	}

	if err != nil {
		log.Fatal("Erro ao fazer ping no MySQL:", err)
	}

	// Cria a tabela 'albums' se ela não existir
	query := `CREATE TABLE IF NOT EXISTS albums (
		id VARCHAR(36) PRIMARY KEY,
		title VARCHAR(100),
		artist VARCHAR(100),
		price DOUBLE
	);`
	_, err = db.Exec(query)
	if err != nil {
		log.Fatal("Erro ao criar tabela:", err)
	}
}

// getAlbums retorna todos os álbuns cadastrados no banco de dados.
func getAlbums(c *gin.Context) {
	rows, err := db.Query("SELECT id, title, artist, price FROM albums")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	var albumsList []album
	for rows.Next() {
		var a album
		if err := rows.Scan(&a.ID, &a.Title, &a.Artist, &a.Price); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		albumsList = append(albumsList, a)
	}

	c.JSON(http.StatusOK, albumsList)
}

// postAlbum insere um novo álbum no banco de dados.
func postAlbum(c *gin.Context) {
	var newAlbum album

	if err := c.ShouldBindJSON(&newAlbum); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	query := "INSERT INTO albums (id, title, artist, price) VALUES (?, ?, ?, ?)"
	_, err := db.Exec(query, newAlbum.ID, newAlbum.Title, newAlbum.Artist, newAlbum.Price)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, newAlbum)
}

// getAlbumByID busca um álbum pelo ID no banco de dados.
func getAlbumByID(c *gin.Context) {
	id := c.Param("id")
	var a album
	query := "SELECT id, title, artist, price FROM albums WHERE id = ?"
	err := db.QueryRow(query, id).Scan(&a.ID, &a.Title, &a.Artist, &a.Price)
	if err == sql.ErrNoRows {
		c.JSON(http.StatusNotFound, gin.H{"message": "album não encontrado"})
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, a)
}

func main() {
	// Carrega o arquivo .env, se existir.
	err := godotenv.Load()
	if err != nil {
		log.Println("Nenhum arquivo .env encontrado. Continuando com as variáveis de ambiente.")
	}

	// Inicializa a conexão com o banco de dados e cria a tabela, se necessário.
	initDB()

	r := gin.Default()

	// Define os endpoints da API.
	r.GET("/albums", getAlbums)
	r.GET("/albums/:id", getAlbumByID)
	r.POST("/albums", postAlbum)

	// Roda a aplicação na porta 8080.
	r.Run(":8080")
}

package connection

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	"log"
	"os"
	"strconv"
)

func CreateConnection() *sqlx.DB {
	// faz a junção de string para conectar ao banco
	var psqlInfo string = ConnectionDB()

	db, err := sqlx.Open("postgres", psqlInfo)
	if err != nil {
		log.Fatal(err)
	}
	return db
}

func ConnectionDB() string {
	// carregar env
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Erro ao carregar o .env")
	}

	// conversão de string para int64
	port := os.Getenv("DB_PORT")

	portN, err := strconv.ParseInt(port, 10, 64)

	if err != nil {
		log.Fatal("Erro ao carregar o .env")
	}

	return fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		os.Getenv("DB_HOST"), portN, os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_NAME"))
}
